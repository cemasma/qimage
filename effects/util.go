package qimage

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
)

func initializeQuantization(sourceImagePath, colorPalette string) (image.Image, *image.RGBA, color.Palette) {
	m := decodeImage(sourceImagePath)
	outputImg := image.NewRGBA(m.Bounds())
	palette := parsePalette(colorPalette)

	return m, outputImg, palette
}

func findMostFrequentColor(colors []color.Color) color.Color {
	type colorIndex struct {
		index int
		count int
	}

	colorIndexMap := make(map[string]colorIndex)

	for index, v := range colors {
		r, g, b, _ := v.RGBA()
		key := strconv.FormatUint(uint64(r), 10) + ":" + strconv.FormatUint(uint64(g), 10) + ":" + strconv.FormatUint(uint64(b), 10)
		val, ok := colorIndexMap[key]

		if !ok {
			colorIndexMap[key] = colorIndex{index, 1}
		} else {
			val.index++
			colorIndexMap[key] = val
		}
	}

	var mostFrequentColor colorIndex

	for _, v := range colorIndexMap {
		if v.count > mostFrequentColor.count {
			mostFrequentColor = v
		}
	}

	if len(colors) == 0 {
		return nil
	}

	return colors[mostFrequentColor.index]
}

func sortPoints(vertices []image.Point) []image.Point {
	for i := 0; i < len(vertices); i++ {
		for j := i + 1; j < len(vertices); j++ {
			if vertices[i].Y > vertices[j].Y {
				vertices[i], vertices[j] = vertices[j], vertices[i]
			}
		}
	}

	return vertices
}

func saveOutputImage(target string, rgba *image.RGBA) {
	f, err := os.Create(target)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := f.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	err = png.Encode(f, rgba)
	if err != nil {
		log.Fatal(err)
	}
}
