package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type Grid struct {
	cells [][]int
}

func InitGrid(x, y int) *Grid {
	g := &Grid{}
	g.cells = make([][]int, y)
	for i, _ := range g.cells {
		g.cells[i] = make([]int, x)
	}
	return g
}

func GridToImage(g *Grid, cellWidth int) *image.RGBA {
	if len(g.cells) == 0 {
		return nil
	}
	width := len(g.cells[0]) * cellWidth
	height := len(g.cells) * cellWidth
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
			if g.cells[y/cellWidth][x/cellWidth] == 1 {
				img.SetRGBA(x, y, color.RGBA{100, 0, 100, 255})
			}
		}
	}
	return img
}

func main() {

	g := InitGrid(8, 4)
	g.cells[1][1] = 1
	g.cells[2][2] = 1
	g.cells[3][3] = 1

	img := GridToImage(g, 10)

	f, err := os.Create("img.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
