package main

import (
	"bufio"
	"fmt"
	"github.com/NikolayOskin/image-optimizer-app/pngquant"
	"github.com/nickalie/go-mozjpegbin"
	"mime/multipart"
	"os"
)

type compressResult struct {
	beforeSize string
	afterSize  string
	saved      int64
}

func optimizeJPEG(fileHeader *multipart.FileHeader, file multipart.File) (compressResult, error) {
	result := compressResult{}
	inputFileSize := fileHeader.Size

	in := bufio.NewReader(file)

	out, err := os.Create(imagesPath + fileHeader.Filename)
	if err != nil {
		return result, err
	}
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

	out, err := os.Create(imagesPath + fileHeader.Filename)
	if err != nil {
		return result, err
	}
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
