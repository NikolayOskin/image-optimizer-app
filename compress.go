package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/NikolayOskin/image-optimizer-app/pngquant"
	"github.com/nickalie/go-mozjpegbin"
)

type compressResult struct {
	beforeSize string
	afterSize  string
	saved      int64
	fileName   string
}

func optimizeJPEG(fileHeader *multipart.FileHeader, file multipart.File) (compressResult, error) {
	result := compressResult{}
	inputFileSize := fileHeader.Size

	in := bufio.NewReader(file)

	out, name, err := createUniqueFile(fileHeader.Filename)
	if err != nil {
		return result, err
	}
	result.fileName = name
	defer out.Close()

	cjpeg := mozjpegbin.NewCJpeg()
	cjpeg.Optimize(true).Quality(70).Output(out).Input(in)
	err = cjpeg.Run()
	if err != nil {
		return result, err
	}

	outStat, err := out.Stat()
	if err != nil {
		return result, err
	}

	outlen := outStat.Size()
	result.beforeSize = ByteCountSI(inputFileSize)
	result.afterSize = ByteCountSI(outlen)

	if outlen < inputFileSize {
		result.saved = (inputFileSize - outlen) * 100 / inputFileSize
		return result, nil
	}
	result.saved = 0
	return result, nil
}

func optimizePNG(fileHeader *multipart.FileHeader, file multipart.File) (compressResult, error) {
	result := compressResult{}
	inputFileSize := fileHeader.Size

	out, name, err := createUniqueFile(fileHeader.Filename)
	if err != nil {
		return result, err
	}
	result.fileName = name
	defer out.Close()

	err = pngquant.Compress(file, out, "1")
	if err != nil {
		return result, err
	}

	outStat, err := out.Stat()
	if err != nil {
		return result, err
	}

	outlen := outStat.Size()
	result.beforeSize = ByteCountSI(inputFileSize)
	result.afterSize = ByteCountSI(outlen)

	if outlen < inputFileSize {
		result.saved = (inputFileSize - outlen) * 100 / inputFileSize
		return result, nil
	}
	result.saved = 0
	return result, nil
}

func createUniqueFile(filename string) (*os.File, string, error) {
	if fileExists(imagesPath + filename) {
		r := time.Now().Unix() + rand.Int63n(100)
		file, err := os.Create(imagesPath + strconv.Itoa(int(r)) + filename)
		if err != nil {
			return nil, "", err
		}
		return file, strconv.Itoa(int(r)) + filename, nil
	} else {
		file, err := os.Create(imagesPath + filename)
		if err != nil {
			return nil, "", err
		}
		return file, filename, nil
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
