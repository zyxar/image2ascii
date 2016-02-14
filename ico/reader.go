package ico

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
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
	d := decoder{r: r}
	var cfg image.Config
	var err error
	if err = d.decodeHeader(); err != nil {
		return cfg, err
	}
	if err = d.decodeEntries(); err != nil {
		return cfg, err
	}
	e := d.entries[0]
	buf := make([]byte, e.Size+14)
	n, err := io.ReadFull(d.r, buf[14:])
	if err != nil && err != io.ErrUnexpectedEOF {
		return cfg, err
	}
	buf = buf[:14+n]
	if string(buf[14:14+len(pngHeader)]) == pngHeader {
		return png.DecodeConfig(bytes.NewReader(buf[14:]))
	}

	if err = d.forgeBMPHead(buf); err != nil {
		return cfg, err
	}
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
	r       io.Reader
	head    head
	entries []entry
	images  []image.Image
	cfg     image.Config
}

func (d *decoder) decode(r io.Reader) (err error) {
	d.r = r
	if err = d.decodeHeader(); err != nil {
		return err
	}
	if err = d.decodeEntries(); err != nil {
		return err
	}
	d.images = make([]image.Image, d.head.Number)
	for i, _ := range d.entries {
		d.images[i], err = d.decodeImage(d.entries[i].Size)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeHeader() error {
	binary.Read(d.r, binary.LittleEndian, &(d.head))
	if d.head.Zero != 0 || d.head.Type != 1 {
		return fmt.Errorf("corrupted head: [%x,%x]", d.head.Zero, d.head.Type)
	}
	return nil
}

func (d *decoder) decodeEntries() error {
	n := int(d.head.Number)
	d.entries = make([]entry, n)
	for i := 0; i < n; i++ {
		if err := binary.Read(d.r, binary.LittleEndian, &(d.entries[i])); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeImage(size uint32) (image.Image, error) {
	data := make([]byte, size+14)
	n, err := io.ReadFull(d.r, data[14:])
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	data = data[:14+n]
	// Check if the image is a PNG by the first 8 bytes of the image data
	if string(data[14:14+len(pngHeader)]) == pngHeader {
		return png.Decode(bytes.NewReader(data[14:]))
	}
	// Decode as BMP instead
	if err = d.forgeBMPHead(data); err != nil {
		return nil, err
	}
	return bmp.Decode(bytes.NewReader(data))
}

func (d *decoder) forgeBMPHead(buf []byte) error {
	// calculate image sizes
	// See wikipedia en.wikipedia.org/wiki/BMP_file_format
	data := buf[14:]
	copy(buf[0:2], "\x42\x4D") // Magic number

	dibSize := binary.LittleEndian.Uint32(data[:4])
	w := binary.LittleEndian.Uint32(data[4:8])
	h := binary.LittleEndian.Uint32(data[8:12])
	if h > w {
		binary.LittleEndian.PutUint32(data[8:12], h/2)
	}

	// File size
	binary.LittleEndian.PutUint32(buf[2:6], uint32(len(data)))

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
	return nil
}

const pngHeader = "\x89PNG\r\n\x1a\n"
