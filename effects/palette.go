package qimage

import (
	"image/color"
	"strconv"
	"strings"
)

func parsePalette(colorPalette string) color.Palette {
	palette := make(color.Palette, 0)
	splittedPalette := strings.Split(colorPalette, ":")

	for _, v := range splittedPalette {
		splittedRGBA := strings.Split(v, ",")

		if len(splittedRGBA) < 3 {
			continue
		}

		r, rErr := strconv.Atoi(splittedRGBA[0])
		g, gErr := strconv.Atoi(splittedRGBA[1])
		b, bErr := strconv.Atoi(splittedRGBA[2])

		if rErr != nil || gErr != nil || bErr != nil {
			continue
		}

		alpha := 255

		if len(splittedRGBA) == 4 {
			a, err := strconv.Atoi(splittedRGBA[4])

			if err == nil {
				alpha = a
			}
		}

		palette = append(palette, color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(alpha),
		})
	}

	return palette
}

func ProcessPalette(sourceImagePath, target, colorPalette string) {
	m, outputImg, palette := initializeQuantization(sourceImagePath, colorPalette)

	for x := 0; x < m.Bounds().Max.X; x++ {
		for y := 0; y < m.Bounds().Max.Y; y++ {
			c := m.At(x, y)
			closestColor := palette.Convert(c)

			if closestColor != nil {
				c = closestColor
			}
			outputImg.Set(x, y, c)
		}
	}

	saveOutputImage(target, outputImg)
}
