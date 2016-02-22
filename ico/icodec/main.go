package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/zyxar/image2ascii/ico"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

var (
	format = flag.String("t", "jpg", "set target file format: [jpg|png|gif|tiff|bmp]")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: icodec [options] {ICON DIR/FILE}\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	dst := flag.Arg(0)
	info, err := os.Stat(dst)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	*format = strings.ToLower(*format)
	conv := func(dst string) {
		if err = convert(dst); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println("[DONE]", dst)
		}
	}
	if info.IsDir() {
		files, err := filepath.Glob(filepath.Join(dst, "*.ico"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if len(files) == 0 {
			fmt.Fprintln(os.Stderr, "no icon files found")
			return
		}
		for _, f := range files {
			conv(f)
		}
	} else {
		conv(dst)
	}
}

func convert(f string) (err error) {
	var rd *os.File
	if rd, err = os.Open(f); err != nil {
		return
	}
	defer rd.Close()
	var images []image.Image
	if images, err = ico.DecodeAll(rd); err != nil {
		err = fmt.Errorf("%s not a ico file: %q", filepath.Base(f), err)
		return
	}

	for i := range images {
		var dst string
		if len(images) == 1 {
			dst = fmt.Sprintf("%s.%s", filepath.Base(f), *format)
		} else {
			dst = fmt.Sprintf("%s-%d.%s", filepath.Base(f), i, *format)
		}
		var w *os.File
		if w, err = os.Create(dst); err != nil {
			continue
		}
		switch *format {
		case "jpg", "jpeg":
			err = jpeg.Encode(w, images[i], &jpeg.Options{Quality: 100})
		case "png":
			err = png.Encode(w, images[i])
		case "gif":
			err = gif.Encode(w, images[i], &gif.Options{})
		case "bmp":
			err = bmp.Encode(w, images[i])
		case "tiff":
			err = tiff.Encode(w, images[i], &tiff.Options{})
		default:
			err = fmt.Errorf("unsupported format: %s", *format)
		}
		if err != nil {
			return
		}
	}
	return
}
