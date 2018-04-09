package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
)

type Grid struct {
	Cells     [][]int
	Tiles     map[int]Cell
	CellWidth int
	GridLines GridLine
}

type Cell struct {
	Color color.RGBA
	Image string
	img   image.Image
}

type GridLine struct {
	Color color.RGBA
	Width int
}

func InitGrid(x, y, cellWidth int) *Grid {
	g := &Grid{}
	g.CellWidth = cellWidth
	g.Cells = make([][]int, y)
	for i, _ := range g.Cells {
		g.Cells[i] = make([]int, x)
	}
	g.Tiles = make(map[int]Cell)
	return g
}

func LoadTiles(g *Grid) {
	fmt.Println("Loading tiles...")
	for i, v := range g.Tiles {
		if v.Image != "" {
			file, err := os.Open(v.Image)
			if err != nil {
				fmt.Println("Could not open ", v.Image, err)
				continue
			}
			defer file.Close()

			im, _, err := image.Decode(file)
			if err != nil {
				fmt.Println("Could not load ", v.Image, err)
				v.img = nil
				g.Tiles[i] = v
			} else {
				fmt.Println(v.Image)
				v.img = im
				g.Tiles[i] = v
			}
		}
	}
}

func GridToImage(g *Grid) *image.RGBA {
	if len(g.Cells) == 0 {
		return nil
	}
	width := len(g.Cells[0])
	height := len(g.Cells)
	img := image.NewRGBA(image.Rect(0, 0, width*g.CellWidth, height*g.CellWidth))

	// paint cells
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v, ok := g.Tiles[g.Cells[y][x]]
			if ok {
				v.Paint(x, y, g.CellWidth, img)
			}
		}
	}

	if g.GridLines.Width == 0 {
		return img
	}

	// add grid lines
	for y := 0; y < height*g.CellWidth; y++ {
		for x := 0; x < width*g.CellWidth; x++ {
			if y%g.CellWidth == 0 || y == height*g.CellWidth-1 {
				for t := -g.GridLines.Width / 2; t <= g.GridLines.Width/2; t++ {
					img.SetRGBA(x, y+t, g.GridLines.Color)
				}
			}
			if x%g.CellWidth == 0 || x == width*g.CellWidth-1 {
				for t := -g.GridLines.Width / 2; t <= g.GridLines.Width/2; t++ {
					img.SetRGBA(x+t, y, g.GridLines.Color)
				}
			}
		}
	}

	return img
}

// Paint a sub region of img with the contents of this cell.
// Color fills the cell with that color
// Image fills the cell with the contents of Image (possibly scaled)
func (a Cell) Paint(x, y, w int, img *image.RGBA) {
	for i := y * w; i < (y+1)*w; i++ {
		for j := x * w; j < (x+1)*w; j++ {
			img.SetRGBA(j, i, a.Color)
		}
	}

	if a.img == nil {
		return
	}

	bounds := a.img.Bounds()
	avgw := int(bounds.Size().X / w)
	avgh := int(bounds.Size().Y / w)

	for i := 0; i < w; i++ {
		for j := 0; j < w; j++ {

			avgr, avgg, avgb := 0, 0, 0
			for k := 0; k < avgh; k++ {
				for l := 0; l < avgw; l++ {
					rgba := color.RGBAModel.Convert(a.img.At(j*avgw+l, i*avgh+k))
					//avg_col.A += rgba.(color.RGBA).A
					avgr += int(rgba.(color.RGBA).R)
					avgg += int(rgba.(color.RGBA).G)
					avgb += int(rgba.(color.RGBA).B)
				}
			}
			avg_col := color.RGBA{
				uint8(avgr / (avgw * avgh)),
				uint8(avgg / (avgw * avgh)),
				uint8(avgb / (avgw * avgh)),
				255,
			}
			img.SetRGBA(x*w+j, y*w+i, avg_col)
		}
	}
}

func main() {

	var gridfile = flag.String("grid", "grid.json", "grid file to load")
	var imagefile = flag.String("image", "img.png", "image file to write")
	flag.Parse()

	// read gridfile and unmarhal into grid object
	dat, err := ioutil.ReadFile(*gridfile)
	if err != nil {
		panic("could not read grid file ")
	}
	var g Grid
	if err := json.Unmarshal(dat, &g); err != nil {
		panic(err)
	}

	// load tile images if any
	LoadTiles(&g)

	// make the image
	img := GridToImage(&g)

	// save it
	f, err := os.Create(*imagefile)
	if err != nil {
		panic("could not create image file")
	}
	defer f.Close()
	png.Encode(f, img)
}
