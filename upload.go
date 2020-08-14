package main

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const jpegType = "image/jpeg"
const jpgType = "image/jpg"
const pngType = "image/png"

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10mb max file size
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		_, _ = w.Write([]byte("could not parse form"))
	}

	file, header, err := r.FormFile("file")

	switch err {
	case nil:
	case http.ErrMissingFile:
		redirectWithErr(w, r, "You didn't choose file to upload")
		return
	default:
		_, _ = w.Write([]byte("something goes wrong"))
		return
	}

	defer file.Close()

	imageType := parseImageType(file)

	if imageType != jpegType && imageType != jpgType && imageType != pngType {
		redirectWithErr(w, r, "File is not correct. Allowed types: jpeg, png")
		return
	}

	path, err := storeFile(&file, header.Filename)
	if err != nil {
		redirectWithErr(w, r, err.Error())
		return
	}

	if imageType == jpegType || imageType == jpgType {
		optimizeJPEG(path)
	}

	http.Redirect(w, r, "/", 301)
	return
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

func redirectWithErr(w http.ResponseWriter, r *http.Request, err string) {
	data.Errors = append(data.Errors, validationError{Error: err})
	http.Redirect(w, r, "/", 301)
}
