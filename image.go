package main

import (
	"fmt"
	"image"
	"io"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

type Image struct {
	img    image.Image
	width  int
	height int
	color  bool
	invert bool
	flipx  bool
	flipy  bool
}

func Decode(r io.Reader) (i *Image, err error) {
	var im image.Image
	if im, _, err = image.Decode(r); err != nil {
		return
	}
	width := termWidth(os.Stdout.Fd())
	height := round(0.5 * float64(width) * float64(im.Bounds().Dy()) / float64(im.Bounds().Dx()))
	img := resize.Resize(uint(width), uint(height), im, resize.Lanczos3)
	i = &Image{img: img, width: width, height: height}
	return
}

func (i Image) WriteTo(w io.Writer) (n int64, err error) {
	ly := i.height
	lx := i.width
	m := 0
	for y := 0; y < ly; y++ {
		for x := 0; x < lx; x++ {
			posX := x
			posY := y
			if i.flipx {
				posX = lx - 1 - x
			}
			if i.flipy {
				posY = ly - 1 - y
			}
			r, g, b, _ := i.img.At(posX, posY).RGBA()
			v := round(float64(r+g+b) * float64(ascii_palette_length) / 65535 / 3)
			if i.invert {
				v = ascii_palette_length - v
			}
			if i.color {
				vr := float64(r) / 65535
				vg := float64(g) / 65535
				vb := float64(b) / 65535
				if vr-vg > threshold_low && vr-vb > threshold_low {
					m, err = fmt.Fprintf(w, "%s", colorRed)
					if err != nil {
						return
					}
				} else if vg-vr > threshold_low && vg-vb > threshold_low {
					m, err = fmt.Fprintf(w, "%s", colorGreen)
					if err != nil {
						return
					}
				} else if vr-vb > threshold_low && vg-vb > threshold_low && vr+vg > threshold_high {
					m, err = fmt.Fprintf(w, "%s", colorYellow)
					if err != nil {
						return
					}
				} else if vb-vr > threshold_low && vb-vg > threshold_low /*&& Y<0.95*/ {
					m, err = fmt.Fprintf(w, "%s", colorBlue)
					if err != nil {
						return
					}
				} else if vr-vg > threshold_low && vb-vg > threshold_low && vr+vb > threshold_high {
					m, err = fmt.Fprintf(w, "%s", colorMagenta)
					if err != nil {
						return
					}
				} else if vg-vr > threshold_low && vb-vr > threshold_low && vb+vg > threshold_high {
					m, err = fmt.Fprintf(w, "%s", colorCyan)
					if err != nil {
						return
					}
				} else {
					m, err = fmt.Fprintf(w, "%s", colorWhite)
					if err != nil {
						return
					}
				}
				n += int64(m)
			}
			m, err = fmt.Fprintf(w, "%c", ascii_palette[v])
			if err != nil {
				return
			}
			n += int64(m)
			if i.color {
				m, err = fmt.Fprintf(w, "%s", colorReset)
				if err != nil {
					return
				}
				n += int64(m)
			}
		}
		m, err = fmt.Fprintln(w)
		if err != nil {
			return
		}
		n += int64(m)
	}
	return
}

func round(f float64) int {
	return int(f + 0.5)
}
