package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// World represents the game state.
type World struct {
	area []bool
	// bitset []uint64
	width  int
	height int
	ltr    bool
}

// NewWorld creates a new world.
func NewWorld(width, height int) *World {
	w := &World{
		area: make([]bool, width*height),
		// bitset: make([]uint64, width*height/64),
		width:  width,
		height: height,
	}
	w.init()
	return w
}

func (w *World) init() {
	// drawSquare(w.area, w.width, w.height, 200, 20, 100)
}

func drawSquare(grid []bool, width, height, startX, startY, size int) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			gx := startX + x
			gy := startY + y

			// Check bounds
			if rand.Intn(5) == 1 {
				if gx >= 0 && gx < width && gy >= 0 && gy < height {
					grid[gy*width+gx] = true
				}
			}
		}
	}
}

func FillCircle(grid []bool, w, h, cx, cy, radius int) {
	r2 := radius * radius

	for dy := -radius; dy <= radius; dy++ {
		y := cy + dy
		if y < 0 || y >= h {
			continue // skip out-of-bounds
		}

		dx := int(math.Sqrt(float64(r2 - dy*dy))) // horizontal distance

		for x := cx - dx; x <= cx+dx; x++ {
			if x < 0 || x >= w {
				continue // skip out-of-bounds
			}
			if rand.Intn(5) == 1 {
				grid[y*w+x] = true
			}
		}
	}
}

func (w *World) Fill(x int, y int) {
	FillCircle(w.area, w.width, w.height, x, y, 15)
}

func (w *World) Update() {

	// moves sand down once per tick
	for b := len(w.area) - 1; b > 0; b-- {

		i := b
		if !w.ltr {
			i = (b/w.width)*w.width + (w.width - 1 - (b % w.width))
		}

		w.ltr = !w.ltr

		if w.area[i] {
			next := i + w.width
			row := next / w.width

			if next > len(w.area)-1 {
				continue
			}

			// 1. check if cell below is free
			if !w.area[next] {
				w.area[i] = false
				w.area[next] = true
				continue
			}
			options := []int{-1, 1}
			firstIndex := rand.Intn(2) // 0 or 1
			first := options[firstIndex]

			// right
			next = i + w.width + first
			nr := next / w.width
			if next < len(w.area)-1 && nr == row && !w.area[next] {
				w.area[i] = false
				w.area[next] = true
				continue

			}
			// 2. check if right cell is free

			// 3. check if left cell is free
			next = i + w.width - options[1-firstIndex]
			nr = next / w.width

			if next < len(w.area)-1 && nr == row && !w.area[next] {
				w.area[i] = false
				w.area[next] = true
				continue

			}
		}
	}
}

func (w *World) Draw(px []byte) {

	for i, v := range w.area {
		if v {
			// Determine coordinates
			x := i % w.width
			y := i / w.width

			// Deterministic light/dark sand pattern
			if (x+y)%2 == 0 {
				// Light sand
				px[i*4] = 0xC2
				px[i*4+1] = 0xB2
				px[i*4+2] = 0x80
				px[i*4+3] = 0xFF
			} else {
				// Dark sand
				px[i*4] = 0xA9
				px[i*4+1] = 0x91
				px[i*4+2] = 0x5F
				px[i*4+3] = 0xFF
			}

		} else {
			// Transparent / empty
			px[i*4] = 0
			px[i*4+1] = 0
			px[i*4+2] = 0
			px[i*4+3] = 0
		}
	}
}

type Game struct {
	// logic
	World *World

	// rgba representation of world.area
	pixels []byte

	// mouse position
	mx int
	my int
}

func (g *Game) Update() error {

	g.mx, g.my = ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.World.Fill(ebiten.CursorPosition())
	}

	g.World.Update()
	g.World.Update()
	g.World.Update()
	g.World.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.World.Draw(g.pixels)

	screen.WritePixels(g.pixels)

	vector.StrokeCircle(
		screen,
		float32(g.mx), float32(g.my),
		15,
		2,
		color.RGBA{140, 140, 140, 1},
		true,
	)

	msg := fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

const (
	// Settings
	screenWidth  = 640
	screenHeight = 360
)

func main() {
	g := &Game{
		World: NewWorld(screenWidth, screenHeight),
	}

	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Sand Simulation")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
