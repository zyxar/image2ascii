package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zyxar/image2ascii/ascii"
)

var (
	width  = flag.Int("w", 0, "set image width")
	height = flag.Int("h", 0, "set image height")
	color  = flag.Bool("c", false, "enable colour mode")
	flipx  = flag.Bool("x", false, "enable flip-x mode")
	flipy  = flag.Bool("y", false, "enable flip-y mode")
	invert = flag.Bool("i", false, "enable invert mode")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: image2ascii [options] {IMAGE FILE}\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	r, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer r.Close()
	opt := ascii.Options{
		Width:  *width,
		Height: *height,
		Color:  *color,
		Invert: *invert,
		Flipx:  *flipx,
		Flipy:  *flipy}

	a, err := ascii.Decode(r, opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if _, err = a.WriteTo(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
