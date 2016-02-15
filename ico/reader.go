package ico

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	bmp "github.com/jsummers/gobmp"
)

func init() {
	image.RegisterFormat("ico", "\x00\x00\x01\x00?????\x00", Decode, DecodeConfig)
}

// ---- public ----

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	if err := d.decode(r); err != nil {
		return nil, err
	}
	return d.images[d.head.Number-1], nil
}

func DecodeAll(r io.Reader) ([]image.Image, error) {
	var d decoder
	if err := d.decode(r); err != nil {
		return nil, err
	}
	return d.images, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var (
		d   decoder
		cfg image.Config
		err error
	)
	if err = d.decodeHeader(r); err != nil {
		return cfg, err
	}
	if err = d.decodeEntries(r); err != nil {
		return cfg, err
	}
	e := d.entries[0]
	buf := make([]byte, e.Size+14)
	n, err := io.ReadFull(r, buf[14:])
	if err != nil && err != io.ErrUnexpectedEOF {
		return cfg, err
	}
	buf = buf[:14+n]
	if string(buf[14:14+len(pngHeader)]) == pngHeader {
		return png.DecodeConfig(bytes.NewReader(buf[14:]))
	}

	d.forgeBMPHead(buf, &e)
	return bmp.DecodeConfig(bytes.NewReader(buf))
}

// ---- private ----

type entry struct {
	Width   byte
	Height  byte
	Palette byte
	_       byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
}

type head struct {
	Zero   uint16
	Type   uint16
	Number uint16
}

type decoder struct {
	head    head
	entries []entry
	images  []image.Image
	cfg     image.Config
}

func (d *decoder) decode(r io.Reader) (err error) {
	if err = d.decodeHeader(r); err != nil {
		return err
	}
	if err = d.decodeEntries(r); err != nil {
		return err
	}
	d.images = make([]image.Image, d.head.Number)
	for i, _ := range d.entries {
		data := make([]byte, d.entries[i].Size+14)
		n, err := io.ReadFull(r, data[14:])
		if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}
		data = data[:14+n]
		// Check if the image is a PNG by the first 8 bytes of the image data
		if string(data[14:14+len(pngHeader)]) == pngHeader {
			if d.images[i], err = png.Decode(bytes.NewReader(data[14:])); err != nil {
				return err
			}
		} else {
			// Decode as BMP instead
			maskData := d.forgeBMPHead(data, &(d.entries[i]))
			if maskData != nil && len(maskData) < 14+n {
				data = data[:n+14-len(maskData)]
			}
			if d.images[i], err = bmp.Decode(bytes.NewReader(data)); err != nil {
				return err
			}
			bounds := d.images[i].Bounds()
			mask := image.NewAlpha(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			masked := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			for row := 0; row < int(d.entries[i].Height); row++ {
				for col := 0; col < int(d.entries[i].Width); col++ {
					if maskData != nil {
						rowSize := (int(d.entries[i].Width) + 31) / 32 * 4
						if (maskData[row*rowSize+col/8]>>(7-uint(col)%8))&0x01 != 1 {
							mask.SetAlpha(col, int(d.entries[i].Height)-row-1, color.Alpha{255})
						}
					} else {
						rowSize := (int(d.entries[i].Width)*32 + 31) / 32 * 4
						offset := int(binary.LittleEndian.Uint32(data[10:14]))
						mask.SetAlpha(col, int(d.entries[i].Height)-row-1, color.Alpha{data[offset+row*rowSize+col*4+3]})
					}
				}
			}
			draw.DrawMask(masked, masked.Bounds(), d.images[i], bounds.Min, mask, bounds.Min, draw.Src)
			d.images[i] = masked
		}
	}
	return nil
}

func (d *decoder) decodeHeader(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &(d.head))
	if d.head.Zero != 0 || d.head.Type != 1 {
		return fmt.Errorf("corrupted head: [%x,%x]", d.head.Zero, d.head.Type)
	}
	return nil
}

func (d *decoder) decodeEntries(r io.Reader) error {
	n := int(d.head.Number)
	d.entries = make([]entry, n)
	for i := 0; i < n; i++ {
		if err := binary.Read(r, binary.LittleEndian, &(d.entries[i])); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) forgeBMPHead(buf []byte, e *entry) (mask []byte) {
	// See wikipedia en.wikipedia.org/wiki/BMP_file_format
	data := buf[14:]
	copy(buf[0:2], "\x42\x4D") // Magic number

	imageSize := len(data)
	if e.Bits != 32 {
		maskSize := (int(e.Width) + 31) / 32 * 4 * int(e.Height)
		imageSize -= maskSize
		mask = data[imageSize:]
	}

	dibSize := binary.LittleEndian.Uint32(data[:4])
	w := binary.LittleEndian.Uint32(data[4:8])
	h := binary.LittleEndian.Uint32(data[8:12])
	if h > w {
		binary.LittleEndian.PutUint32(data[8:12], h/2)
	}

	// File size
	binary.LittleEndian.PutUint32(buf[2:6], uint32(imageSize))

	// Calculate offset into image data
	numColors := binary.LittleEndian.Uint32(data[32:36])
	bits := binary.LittleEndian.Uint16(data[14:16])

	switch bits {
	case 1, 2, 4, 8:
		x := uint32(1 << bits)
		if numColors == 0 || numColors > x {
			numColors = x
		}
	default:
		numColors = 0
	}

	var numColorsSize uint32
	switch dibSize {
	case 12, 64:
		numColorsSize = numColors * 3
	default:
		numColorsSize = numColors * 4
	}
	offset := 14 + dibSize + numColorsSize
	if dibSize > 40 {
		offset += binary.LittleEndian.Uint32(data[dibSize-8 : dibSize-4])
	}
	binary.LittleEndian.PutUint32(buf[10:14], offset)
	return
}

const pngHeader = "\x89PNG\r\n\x1a\n"
