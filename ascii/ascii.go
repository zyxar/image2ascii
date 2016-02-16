package ascii

import (
	"fmt"
	"image"
	"io"
	"os"

	"github.com/nfnt/resize"
	"github.com/zyxar/image2ascii/term"

	_ "github.com/zyxar/image2ascii/ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Ascii struct {
	img image.Image
	opt Options
}

type Options struct {
	Width  int
	Height int
	Color  bool
	Invert bool
	Flipx  bool
	Flipy  bool
}

func Encode(w io.Writer, m image.Image, o ...Options) error {
	img, err := decode(m, o...)
	if err != nil {
		return err
	}
	_, err = img.WriteTo(w)
	return err
}

func decode(m image.Image, o ...Options) (a *Ascii, err error) {
	var opt Options
	if o != nil {
		opt = o[0]
	} else {
		opt = Options{
			Color:  false,
			Invert: false,
			Flipx:  false,
			Flipy:  false}
	}

	if opt.Width <= 0 && opt.Height <= 0 {
		opt.Width = term.Width(os.Stdout.Fd())
		opt.Height = round(0.5 * float64(opt.Width) * float64(m.Bounds().Dy()) / float64(m.Bounds().Dx()))
	} else if opt.Height <= 0 {
		opt.Height = round(0.5 * float64(opt.Width) * float64(m.Bounds().Dy()) / float64(m.Bounds().Dx()))
	} else if opt.Width <= 0 {
		opt.Width = round(2 * float64(opt.Height) * float64(m.Bounds().Dx()) / float64(m.Bounds().Dy()))
	}

	img := resize.Resize(uint(opt.Width), uint(opt.Height), m, resize.Lanczos3)
	a = &Ascii{img, opt}
	return
}

func Decode(r io.Reader, o ...Options) (a *Ascii, err error) {
	var img image.Image
	if img, _, err = image.Decode(r); err != nil {
		return
	}
	a, err = decode(img, o...)
	return
}

func (this Ascii) WriteTo(w io.Writer) (n int64, err error) {
	ly := this.opt.Height
	lx := this.opt.Width
	m := 0
	for y := 0; y < ly; y++ {
		for x := 0; x < lx; x++ {
			posX := x
			posY := y
			if this.opt.Flipx {
				posX = lx - 1 - x
			}
			if this.opt.Flipy {
				posY = ly - 1 - y
			}
			r, g, b, _ := this.img.At(posX, posY).RGBA()
			v := round(float64(r+g+b) * float64(ascii_palette_length) / 65535 / 3)
			if this.opt.Invert {
				v = ascii_palette_length - v
			}
			if this.opt.Color {
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
			if this.opt.Color {
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
