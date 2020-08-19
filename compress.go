package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nickalie/go-mozjpegbin"
	"github.com/yusukebe/go-pngquant"
	"image/png"
	"io"
	"io/ioutil"
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

func optimizePNG(path string) (compressResult, error) {
	result := compressResult{}

	finput, err := os.Open(path)
	if err != nil {
		return result, err
	}
	input, err := ioutil.ReadAll(finput)
	if err != nil {
		return result, err
	}
	in := bytes.NewReader(input)
	img, err := png.Decode(in)
	finput.Close()
	if err != nil {
		return result, err
	}

	// Encode image
	out := new(bytes.Buffer)
	cimg, err := pngquant.Compress(img, "1")
	if err != nil {
		return result, err
	}
	err = png.Encode(out, cimg)
	if err != nil {
		return result, err
	}

	outlen := int64(out.Len())

	result.beforeSize = ByteCountSI(in.Size())
	result.afterSize = ByteCountSI(outlen)

	if outlen < in.Size() {
		// Write to file
		f, err := os.Create(path)
		if err != nil {
			return result, err
		}
		_, err = io.Copy(f, out)
		if err != nil {
			return result, err
		}
		f.Close()

		result.saved = (in.Size() - outlen) * 100 / in.Size()
		return result, nil
	} else {
		return result, nil
	}
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
