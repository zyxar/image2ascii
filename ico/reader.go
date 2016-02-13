package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"

	bmp "github.com/jsummers/gobmp"
)

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

const pngHeader = "\x89PNG\r\n\x1a\n"

// A FormatError reports that the input is not a valid ICO.
type FormatError string

func (e FormatError) Error() string { return "invalid ICO format: " + string(e) }

// If the io.Reader does not also have ReadByte, then decode will introduce its own buffering.
type reader interface {
	io.Reader
	io.ByteReader
}

type decoder struct {
	r       reader
	num     uint16
	entries []entry
	images  []image.Image
	cfg     image.Config
}

func (d *decoder) decode(r io.Reader, configOnly bool) error {
	// Add buffering if r does not provide ReadByte.
	if rr, ok := r.(reader); ok {
		d.r = rr
	} else {
		d.r = bufio.NewReader(r)
	}

	if err := d.decodeHeader(); err != nil {
		return err
	}
	if err := d.decodeEntries(configOnly); err != nil {
		return err
	}
	if configOnly {
		cfg, err := d.decodeConfig(d.entries[0])
		if err != nil {
			return err
		}
		d.cfg = cfg
	} else {
		d.images = make([]image.Image, d.num)
		for i, entry := range d.entries {
			img, err := d.decodeImage(entry)
			if err != nil {
				return err
			}
			d.images[i] = img
		}
	}
	return nil
}

func (d *decoder) decodeHeader() error {
	var first, second uint16
	binary.Read(d.r, binary.LittleEndian, &first)
	binary.Read(d.r, binary.LittleEndian, &second)
	binary.Read(d.r, binary.LittleEndian, &d.num)
	if first != 0 {
		return FormatError(fmt.Sprintf("first byte is %d instead of 0", first))
	}
	if second != 1 {
		return FormatError(fmt.Sprintf("second byte is %d instead of 1", second))
	}
	return nil
}

func (d *decoder) decodeEntries(configOnly bool) error {
	n := int(d.num)
	if configOnly {
		n = 1
	}
	d.entries = make([]entry, n)
	for i := 0; i < n; i++ {
		if err := binary.Read(d.r, binary.LittleEndian, &(d.entries[i])); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeImage(e entry) (image.Image, error) {
	data := make([]byte, e.Size)
	io.ReadFull(d.r, data)

	// Check if the image is a PNG by the first 8 bytes of the image data
	if string(data[:len(pngHeader)]) == pngHeader {
		return png.Decode(bytes.NewReader(data))
	}

	// Decode as BMP instead
	bmpBytes, _, _, err := d.setupBMP(e, data)
	if err != nil {
		return nil, err
	}

	return bmp.Decode(bytes.NewReader(bmpBytes))
}

func (d *decoder) decodeConfig(e entry) (cfg image.Config, err error) {
	tmp := make([]byte, e.Size)
	n, err := io.ReadFull(d.r, tmp)
	if n != int(e.Size) {
		return cfg, fmt.Errorf("Only %d of %d bytes read.", n, e.Size)
	}
	if err != nil {
		return cfg, err
	}

	cfg, err = png.DecodeConfig(bytes.NewReader(tmp))
	if err != nil {
		tmp, _, _, _ = d.setupBMP(e, tmp)
		cfg, err = bmp.DecodeConfig(bytes.NewReader(tmp))
	}
	return cfg, err
}

func (d *decoder) setupBMP(e entry, data []byte) ([]byte, []byte, int, error) {
	// Ico files are made up of a XOR mask and an AND mask
	// The XOR mask is the image itself, while the AND mask is a 1 bit-per-pixel alpha channel.
	// setupBMP returns the image as a BMP format byte array, and the mask as a (1bpp) pixel array

	// calculate image sizes
	// See wikipedia en.wikipedia.org/wiki/BMP_file_format
	var imageSize, maskSize int
	if int(e.Size) < len(data) {
		imageSize = int(e.Size)
	} else {
		imageSize = len(data)
	}
	if e.Bits != 32 {
		rowSize := (1 * (int(e.Width) + 31) / 32) * 4
		maskSize = rowSize * int(e.Height)
		imageSize -= maskSize
	}

	img := make([]byte, 14+imageSize)
	mask := make([]byte, maskSize)

	var n int
	// Read in image
	n = copy(img[14:], data[:imageSize])
	if n != imageSize {
		return nil, nil, 0, FormatError(fmt.Sprintf("only %d of %d bytes read.", n, imageSize))
	}
	// Read in mask
	n = copy(mask, data[imageSize:])
	if n != maskSize {
		return nil, nil, 0, FormatError(fmt.Sprintf("only %d of %d bytes read.", n, maskSize))
	}

	var dibSize, w, h uint32
	dibSize = binary.LittleEndian.Uint32(img[14 : 14+4])
	w = binary.LittleEndian.Uint32(img[14+4 : 14+8])
	h = binary.LittleEndian.Uint32(img[14+8 : 14+12])

	if h > w {
		binary.LittleEndian.PutUint32(img[14+8:14+12], h/2)
	}

	// Magic number
	copy(img[0:2], "\x42\x4D")

	// File size
	binary.LittleEndian.PutUint32(img[2:6], uint32(imageSize+14))

	// Calculate offset into image data
	numColors := binary.LittleEndian.Uint32(img[14+32 : 14+36])
	e.Bits = binary.LittleEndian.Uint16(img[14+14 : 14+16])
	e.Size = binary.LittleEndian.Uint32(img[14+20 : 14+24])

	switch int(e.Bits) {
	case 1, 2, 4, 8:
		x := uint32(1 << e.Bits)
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
		offset += binary.LittleEndian.Uint32(img[14+dibSize-8 : 14+dibSize-4])
	}
	binary.LittleEndian.PutUint32(img[10:14], offset)

	return img, mask, int(offset), nil
}

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return nil, err
	}
	return d.images[0], nil
}

func DecodeAll(r io.Reader) ([]image.Image, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return nil, err
	}
	return d.images, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var d decoder
	if err := d.decode(r, true); err != nil {
		return image.Config{}, err
	}
	return d.cfg, nil
}
