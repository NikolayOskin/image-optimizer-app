package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (r *compressResult) countCompressedSizes(inputFileSize int64, outFile *os.File) error {
	if inputFileSize == 0 {
		return errors.New("inputFileSize can not be 0")
	}

	outStat, err := outFile.Stat()
	if err != nil {
		return err
	}

	outlen := outStat.Size()
	r.beforeSize = ByteCountSI(inputFileSize)
	r.afterSize = ByteCountSI(outlen)

	if outlen < inputFileSize {
		r.saved = (inputFileSize - outlen) * 100 / inputFileSize
		return nil
	}
	return nil
}

func compressJPEG(fileHeader *multipart.FileHeader, file multipart.File) (compressResult, error) {
	result := compressResult{}
	inputFileSize := fileHeader.Size

	in := bufio.NewReader(file)

	outFile, name, err := createUniqueImageFile(fileHeader.Filename)
	if err != nil {
		return result, err
	}
	result.fileName = name
	defer outFile.Close()

	cjpeg := mozjpegbin.NewCJpeg()
	cjpeg.Optimize(true).Quality(70).Output(outFile).Input(in)
	err = cjpeg.Run()
	if err != nil {
		return result, err
	}

	err = result.countCompressedSizes(inputFileSize, outFile)
	if err != nil {
		return result, err
	}
	return result, nil
}

func compressPNG(fileHeader *multipart.FileHeader, file multipart.File) (compressResult, error) {
	result := compressResult{}
	inputFileSize := fileHeader.Size

	outFile, name, err := createUniqueImageFile(fileHeader.Filename)
	if err != nil {
		return result, err
	}
	result.fileName = name
	defer outFile.Close()

	err = pngquant.Compress(file, outFile, "1")
	if err != nil {
		return result, err
	}

	err = result.countCompressedSizes(inputFileSize, outFile)
	if err != nil {
		return result, err
	}
	return result, nil
}

func createUniqueImageFile(filename string) (*os.File, string, error) {
	if fileExists(filepath.Join(imagesPath, filename)) {
		r := time.Now().Unix() + rand.Int63n(100)
		file, err := os.Create(filepath.Join(imagesPath, strconv.Itoa(int(r))+filename))
		if err != nil {
			return nil, "", err
		}
		return file, strconv.Itoa(int(r)) + filename, nil
	}
	file, err := os.Create(filepath.Join(imagesPath, filename))
	if err != nil {
		return nil, "", err
	}
	return file, filename, nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ByteCountSI - convert byte size to kb/mb
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
