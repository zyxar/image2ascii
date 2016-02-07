package main

import (
	"flag"
	"fmt"
	"os"
)

var (
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer r.Close()
	i, err := Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	i.color = *color
	i.invert = *invert
	i.flipy = *flipy
	i.flipx = *flipx
	if _, err = i.WriteTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
