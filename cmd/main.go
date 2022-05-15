package main

import (
	"flag"
	"fmt"
	qimage "qimage/effects"
)

func main() {
	sourceImagePath := flag.String("f", "", "--f image.png")
	target := flag.String("o", "output.png", "-o output.png")
	randomTriangulation := flag.Bool("tr", false, "--tr true")
	uniformTriangulation := flag.Bool("tu", false, "--tu true")
	lineEffect := flag.Bool("l", false, "--l true")
	colorPalette := flag.String("cp", "", "--cp 255,255,255:0,0,0")

	flag.Parse()

	if *sourceImagePath == "" {
		fmt.Println("file parameter cannot be empty.")
		return
	}

	if !*randomTriangulation && !*uniformTriangulation && !*lineEffect && *colorPalette == "" {
		fmt.Println("at least one effect should be chosen.")
		return
	}

	if *randomTriangulation {
		qimage.RandomTriangulation(*sourceImagePath, *target, *colorPalette, 0, 0)
	}

	if *uniformTriangulation {
		qimage.UniformTriangulation(*sourceImagePath, *target, *colorPalette)
	}

	if *lineEffect {
		qimage.RandomLines(*sourceImagePath, *target, *colorPalette)
	}

	if *colorPalette != "" && !*randomTriangulation && !*uniformTriangulation && !*lineEffect {
		qimage.ProcessPalette(*sourceImagePath, *target, *colorPalette)
	}
}
