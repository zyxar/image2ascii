# Image to ASCII
[![GoDoc][1]][2]
[![Build Status][3]][4]

[1]: https://godoc.org/github.com/zyxar/image2ascii/ascii?status.svg
[2]: https://godoc.org/github.com/zyxar/image2ascii/ascii
[3]: https://travis-ci.org/zyxar/image2ascii.png?branch=master
[4]: https://travis-ci.org/zyxar/image2ascii

inspired by [jp2a](http://csl.sublevel3.org/jp2a/)

## Install

`go get github.com/zyxar/image2ascii`

## Supported [Formats](https://en.wikipedia.org/wiki/Image_file_formats)

- [x] `jpeg`
- [x] `png`
- [x] `gif`
- [x] `bmp`
- [x] `tiff`
- [x] `ico`
- [x] `vp8l`
- [x] `webp`
- [ ] [`exif`](https://en.wikipedia.org/wiki/Exchangeable_image_file_format)
- [ ] [`bgp`](https://en.wikipedia.org/wiki/Better_Portable_Graphics)
- [ ] [`svg`](https://en.wikipedia.org/wiki/Scalable_Vector_Graphics)


## Cmd

```
Usage: image2ascii [options] {IMAGE FILE}
  -c  enable colour mode
  -h int
      set image height
  -i  enable invert mode
  -w int
      set image width
  -x  enable flip-x mode
  -y  enable flip-y mode
```

## Package

```godoc
package ascii
    import "github.com/zyxar/image2ascii/ascii"


FUNCTIONS

func Encode(w io.Writer, m image.Image, c ...Config) error

TYPES

type Ascii struct {
    // contains filtered or unexported fields
}

func Decode(r io.Reader, c ...Config) (a *Ascii, err error)

func (this Ascii) WriteTo(w io.Writer) (n int64, err error)

type Config struct {
    Width  int
    Height int
    Color  bool
    Invert bool
    Flipx  bool
    Flipy  bool
}
```

```godoc
package ico
    import "github.com/zyxar/image2ascii/ico"


FUNCTIONS

func Decode(r io.Reader) (image.Image, error)

func DecodeAll(r io.Reader) ([]image.Image, error)

func DecodeConfig(r io.Reader) (image.Config, error)
```

```godoc
package term
    import "github.com/zyxar/image2ascii/term"


FUNCTIONS

func Width(fd uintptr) int
```

## Alternative

[node-jp2a](https://github.com/zyxar/node-jp2a): `npm install -g node-jp2a`


## License
[Apache 2.0](http://opensource.org/licenses/Apache-2.0)
