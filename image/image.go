package image

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
	image  image.Image
	config Config
}

type Config struct {
	Width  int
	Height int
	Color  bool
	Invert bool
	Flipx  bool
	Flipy  bool
}

func Decode(r io.Reader, c ...Config) (i *Image, err error) {
	var im image.Image
	if im, _, err = image.Decode(r); err != nil {
		return
	}
	var conf Config
	if c != nil {
		conf = c[0]
	} else {
		conf = Config{
			Color:  false,
			Invert: false,
			Flipx:  false,
			Flipy:  false}
	}

	if conf.Width <= 0 && conf.Height <= 0 {
		conf.Width = termWidth(os.Stdout.Fd())
		conf.Height = round(0.5 * float64(conf.Width) * float64(im.Bounds().Dy()) / float64(im.Bounds().Dx()))
	} else if conf.Height <= 0 {
		conf.Height = round(0.5 * float64(conf.Width) * float64(im.Bounds().Dy()) / float64(im.Bounds().Dx()))
	} else if conf.Width <= 0 {
		conf.Width = round(2 * float64(conf.Height) * float64(im.Bounds().Dx()) / float64(im.Bounds().Dy()))
	}

	img := resize.Resize(uint(conf.Width), uint(conf.Height), im, resize.Lanczos3)
	i = &Image{img, conf}
	return
}

func (i Image) WriteTo(w io.Writer) (n int64, err error) {
	ly := i.config.Height
	lx := i.config.Width
	m := 0
	for y := 0; y < ly; y++ {
		for x := 0; x < lx; x++ {
			posX := x
			posY := y
			if i.config.Flipx {
				posX = lx - 1 - x
			}
			if i.config.Flipy {
				posY = ly - 1 - y
			}
			r, g, b, _ := i.image.At(posX, posY).RGBA()
			v := round(float64(r+g+b) * float64(ascii_palette_length) / 65535 / 3)
			if i.config.Invert {
				v = ascii_palette_length - v
			}
			if i.config.Color {
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
			if i.config.Color {
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