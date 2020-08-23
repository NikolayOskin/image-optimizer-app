package main

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
)

const jpegType = "image/jpeg"
const jpgType = "image/jpg"
const pngType = "image/png"

func upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10mb max file size
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		redirectWithErr(w, r, "Could not parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if header == nil {
		redirectWithErr(w, r, "File is corrupted")
		return
	}
	if errWhileUpload(w, r, err) {
		return
	}
	defer file.Close()

	fileType, err := detectFileType(file)
	if err != nil {
		redirectWithErr(w, r, "Could not parse file")
		return
	}
	if fileType != jpegType && fileType != jpgType && fileType != pngType {
		redirectWithErr(w, r, "File is not correct. Allowed types: jpeg, jpg, png")
		return
	}

	result, err := storeCompressedImage(fileType, header, file)
	if err != nil {
		log.Printf("error while optimizing image file: %v", err)
		redirectWithErr(w, r, "something goes wrong")
		return
	}

	redirectToResultPage(w, r, result.fileName, result)
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

func storeCompressedImage(
	imgType string,
	header *multipart.FileHeader,
	file multipart.File,
) (compressResult, error) {
	var result compressResult

	if imgType == jpegType || imgType == jpgType {
		result, err := compressJPEG(header, file)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	if imgType == pngType {
		result, err := compressPNG(header, file)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	return result, errors.New("file type passed is not jpeg or png")
}

func detectFileType(file multipart.File) (string, error) {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	fileType := http.DetectContentType(fileHeader)

	return fileType, nil
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
