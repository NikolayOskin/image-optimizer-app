package main

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const jpegType = "image/jpeg"
const jpgType = "image/jpg"
const pngType = "image/png"

func upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10mb max file size
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		_, _ = w.Write([]byte("could not parse form"))
	}

	file, header, err := r.FormFile("file")
	if errWhileUpload(w, r, err) {
		return
	}
	defer file.Close()

	imageType := parseImageType(file)
	if imageType != jpegType && imageType != jpgType && imageType != pngType {
		redirectWithErr(w, r, "File is not correct. Allowed types: jpeg, jpg, png")
		return
	}

	path, err := storeFile(&file, header.Filename)
	if err != nil {
		redirectWithErr(w, r, err.Error())
		return
	}

	result, err := optimizeImage(imageType, path)
	if err != nil {
		log.Printf("error while optimizing image file: %v", err)
		redirectWithErr(w, r, "something goes wrong")
		return
	}

	redirectToResultPage(w, r, header.Filename, result)
}

func errWhileUpload(w http.ResponseWriter, r *http.Request, err error) bool {
	switch err {
	case nil:
	case http.ErrMissingFile:
		redirectWithErr(w, r, "You didn't choose file to upload")
		return true
	default:
		_, _ = w.Write([]byte("something goes wrong"))
		return true
	}
	return false
}

func optimizeImage(imageType string, path string) (compressResult, error) {
	var result compressResult

	if imageType == jpegType || imageType == jpgType {
		result, err := optimizeJPEG(path)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	if imageType == pngType {
		result, err := optimizePNG(path)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	return result, errors.New("file type passed is not jpeg or png")
}

func storeFile(file *multipart.File, filename string) (string, error) {
	uploadedFile, err := os.OpenFile(imagesPath+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", errors.New("could not upload file")
	}
	defer uploadedFile.Close()

	if _, err := io.Copy(uploadedFile, *file); err != nil {
		return "", errors.New("could not upload file")
	}
	return imagesPath + filename, nil
}

func parseImageType(file multipart.File) string {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		panic(err.Error())
	}
	if _, err := file.Seek(0, 0); err != nil {
		panic(err.Error())
	}
	fileType := http.DetectContentType(fileHeader)

	return fileType
}

func redirectToResultPage(
	w http.ResponseWriter,
	r *http.Request,
	filename string,
	result compressResult,
) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["filename"] = filename
	session.Values["beforeSize"] = result.beforeSize
	session.Values["afterSize"] = result.afterSize
	session.Values["saved"] = result.saved

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/result", 301)
}

func redirectWithErr(w http.ResponseWriter, r *http.Request, errText string) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(errText)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 301)
}
