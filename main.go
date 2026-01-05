package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	drawSquare(w.area, w.width, w.height, 200, 20, 40)
}

func drawSquare(grid []bool, width, height, startX, startY, size int) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			gx := startX + x
			gy := startY + y

			// Check bounds
			if gx >= 0 && gx < width && gy >= 0 && gy < height {
				grid[gy*width+gx] = true
			}
		}
	}
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
			px[i*4] = 0xFF
			px[i*4+1] = 0xFF
			px[i*4+2] = 0xFF
			px[i*4+3] = 0xFF
		} else {
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
}

func (g *Game) Update() error {
	g.World.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.World.Draw(g.pixels)

	screen.WritePixels(g.pixels)

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
