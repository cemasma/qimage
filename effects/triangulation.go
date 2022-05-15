package qimage

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

const MaxEdgeLength = 60
const MinEdgeLength = 5
const TriangleRatio = 0.00001708897089

type Triangulation struct {
	Points    [][]int `json:"points"`
	Simplices [][]int `json:"simplices"`
}

func RandomLines(sourceImagePath, target, colorPalette string) {
	RandomTriangulation(sourceImagePath, target, colorPalette, 60, 1)
}

func RandomTriangulation(sourceImagePath, target, colorPalette string, minEdgeLength, maxEdgeLength int) {
	m, outputImg, palette := initializeQuantization(sourceImagePath, colorPalette)

	if minEdgeLength == 0 || maxEdgeLength == 0 {
		minEdgeLength = MinEdgeLength
		maxEdgeLength = MaxEdgeLength
	}

	maxTriangleCount := (float64(m.Bounds().Max.X) / float64(m.Bounds().Max.Y)) / TriangleRatio

	gatherRandomPoints := func() chan []image.Point {
		results := make(chan []image.Point, int(maxTriangleCount))

		go func() {
			defer close(results)
			for i := 0; i < int(maxTriangleCount); i++ {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))

				v1, v2, v3 := pickRandomPoints(m.Bounds().Max.X, m.Bounds().Max.Y, minEdgeLength, maxEdgeLength, random)

				results <- []image.Point{v1, v2, v3}
			}
		}()

		return results
	}

	for points := range gatherRandomPoints() {
		v1, v2, v3 := points[0], points[1], points[2]

		go dyeTriangle(v1, v2, v3, palette, m, outputImg)
	}

	saveOutputImage(target, outputImg)
}

func UniformTriangulation(sourceImagePath, target, colorPalette string) {
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
	cmd := exec.Command("python3", "effects/main.py", sourceImagePath, "5000")

	_, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	jsonReader, err := os.Open("triangles.json")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := jsonReader.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	bytes, err := ioutil.ReadAll(jsonReader)

	if err != nil {
		log.Fatal(err)
	}

	var triangulation Triangulation
	err = json.Unmarshal(bytes, &triangulation)

	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for _, v := range triangulation.Simplices {
		wg.Add(1)
		i1, i2, i3 := v[0], v[1], v[2]

		p1, p2, p3 := triangulation.Points[i1], triangulation.Points[i2], triangulation.Points[i3]

		v1 := image.Point{X: p1[0], Y: p1[1]}
		v2 := image.Point{X: p2[0], Y: p2[1]}
		v3 := image.Point{X: p3[0], Y: p3[1]}

		go func() {
			dyeTriangle(v1, v2, v3, palette, m, outputImg)

			wg.Done()
		}()
	}
	wg.Wait()

	saveOutputImage(target, outputImg)
}

func dyeTriangle(v1, v2, v3 image.Point, palette color.Palette, m image.Image, outputImg *image.RGBA) {
	points := sortPoints([]image.Point{v1, v2, v3})
	v1, v2, v3 = points[0], points[1], points[2]
	frequentColor, trianglePoints := triangle(v1, v2, v3, m)

	c := frequentColor

	if len(palette) > 0 {
		closestColor := palette.Convert(frequentColor)

		if closestColor != nil {
			c = closestColor
		}
	}

	for _, p := range trianglePoints {
		outputImg.Set(p.X, p.Y, c)
	}
}

func triangle(v1, v2, v3 image.Point, inputImage image.Image) (color.Color, []image.Point) {
	colors := make([]color.Color, 0)
	points := make([]image.Point, 0)

	if v2.Y == v3.Y {
		fillBottomFlatTriangle(v1, v2, v3, inputImage, &colors, &points)
	} else if v1.Y == v2.Y {
		fillTopFlatTriangle(v1, v2, v3, inputImage, &colors, &points)
	} else {
		v4 := image.Point{X: v1.X + ((v2.Y-v1.Y)/(v3.Y-v1.Y))*(v3.X-v1.X), Y: v2.Y}
		fillBottomFlatTriangle(v1, v2, v4, inputImage, &colors, &points)
		fillTopFlatTriangle(v2, v4, v3, inputImage, &colors, &points)
	}

	mostFrequentColor := findMostFrequentColor(colors)

	return mostFrequentColor, points
}

func decodeImage(sourceImagePath string) image.Image {
	reader, err := os.Open(sourceImagePath)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := reader.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	m, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal(err)
	}

	return m
}

func fillBottomFlatTriangle(v1, v2, v3 image.Point, inputImage image.Image, colors *[]color.Color, points *[]image.Point) {
	if v2.Y-v1.Y == 0 || v3.Y-v1.Y == 0 {
		return
	}

	invslope1 := (v2.X - v1.X) / (v2.Y - v1.Y)
	invslope2 := (v3.X - v1.X) / (v3.Y - v1.Y)

	curx1 := v1.X
	curx2 := v1.X

	for scanlineY := v1.Y; scanlineY <= v2.Y; scanlineY++ {
		hLineColor(curx1, scanlineY, curx2, inputImage, colors, points)

		curx1 += invslope1
		curx2 += invslope2
	}
}

func fillTopFlatTriangle(v1, v2, v3 image.Point, inputImage image.Image, colors *[]color.Color, points *[]image.Point) {
	invslope1 := (v3.X - v1.X) / (v3.Y - v1.Y)
	invslope2 := (v3.X - v2.X) / (v3.Y - v2.Y)

	curx1 := v3.X
	curx2 := v3.X

	for scanlineY := v3.Y; scanlineY > v1.Y; scanlineY-- {
		hLineColor(curx1, scanlineY, curx2, inputImage, colors, points)

		curx1 -= invslope1
		curx2 -= invslope2
	}

}

func hLineColor(x1, y, x2 int, inputImage image.Image, colors *[]color.Color, points *[]image.Point) {
	smallX := x1
	bigX := x2

	if bigX < smallX {
		smallX = x2
		bigX = x1
	}

	for ; smallX <= bigX; smallX++ {
		*points = append(*points, image.Point{X: smallX, Y: y})
		*colors = append(*colors, inputImage.At(smallX, y))
	}
}

func pickRandomPoints(maxX, maxY int, minEdgeLength, maxEdgeLength int, random *rand.Rand) (v1, v2, v3 image.Point) {
	randomNumberForX := random.Intn(maxX)
	randomNumberForY := random.Intn(maxY)

	v1.X = randomNumberForX
	v1.Y = randomNumberForY

	verticeArr := []image.Point{v2, v3}

	for i := 0; i < 2; i++ {
		randomNumberForX = random.Intn(maxEdgeLength) + minEdgeLength
		randomNumberForY = random.Intn(maxEdgeLength) + minEdgeLength

		plusOrMinus := random.Intn(1)
		if plusOrMinus == 1 {
			verticeArr[i].X = randomNumberForX + v1.X
		} else {
			verticeArr[i].X = v1.X - randomNumberForX
		}

		plusOrMinus = random.Intn(1)
		if plusOrMinus == 1 {
			verticeArr[i].Y = randomNumberForY + v1.Y
		} else {
			verticeArr[i].Y = v1.Y - randomNumberForY
		}
	}

	v2 = verticeArr[0]
	v3 = verticeArr[1]

	return
}
