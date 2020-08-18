package main

import (
	"bytes"
	"github.com/nickalie/go-mozjpegbin"
	pngquant "github.com/yusukebe/go-pngquant"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
)

func optimizeJPEG(path string) (int64, error) {
	finput, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	input, err := ioutil.ReadAll(finput)
	if err != nil {
		return 0, err
	}
	in := bytes.NewReader(input)
	img, err := jpeg.Decode(in)
	finput.Close()
	if err != nil {
		return 0, err
	}

	// Encode image
	out := new(bytes.Buffer)
	err = mozjpegbin.Encode(out, img, &mozjpegbin.Options{
		Quality:  70,
		Optimize: true,
	})
	if err != nil {
		return 0, err
	}

	outlen := int64(out.Len())
	if outlen < in.Size() {
		// Write to file
		f, err := os.Create(path)
		if err != nil {
			return 0, err
		}
		_, err = io.Copy(f, out)
		if err != nil {
			return 0, err
		}
		f.Close()

		saved := (in.Size() - outlen) * 100 / in.Size()
		return saved, nil
	} else {
		return 0, nil
	}
}

func optimizePNG(path string) (int64, error) {
	finput, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	input, err := ioutil.ReadAll(finput)
	if err != nil {
		return 0, err
	}
	in := bytes.NewReader(input)
	img, err := png.Decode(in)
	finput.Close()
	if err != nil {
		return 0, err
	}

	// Encode image
	out := new(bytes.Buffer)
	cimg, err := pngquant.Compress(img, "1")
	if err != nil {
		return 0, err
	}
	err = png.Encode(out, cimg)
	if err != nil {
		return 0, err
	}

	outlen := int64(out.Len())
	if outlen < in.Size() {
		// Write to file
		f, err := os.Create(path)
		if err != nil {
			return 0, err
		}
		_, err = io.Copy(f, out)
		if err != nil {
			return 0, err
		}
		f.Close()

		saved := (in.Size() - outlen) * 100 / in.Size()
		return saved, nil
	} else {
		return 0, nil
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
