package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/bits"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bitset struct {
	BitArr []uint64
	Size   int
}

func NewBitset(size int) *Bitset {
	// always create with enough space
	sliceSize := (size + 63) / 64
	return &Bitset{
		BitArr: make([]uint64, sliceSize),
		Size:   size,
	}
}

func (bs *Bitset) Set(pos int) {
	if pos < 0 || pos >= bs.Size {
		return
	}

	idx := pos / 64
	ofs := pos % 64

	bs.BitArr[idx] |= 1 << ofs
}

func (bs *Bitset) Clear(pos int) {
	if pos < 0 || pos >= bs.Size {
		return
	}
	idx := pos / 64
	ofs := pos % 64
	bs.BitArr[idx] &^= (1 << ofs)
}

func (bs *Bitset) Test(pos int) bool {
	if pos < 0 || pos >= bs.Size {
		return false
	}

	idx := pos / 64
	ofs := pos % 64

	return bs.BitArr[idx]&(1<<ofs) != 0
}

func (bs *Bitset) SetClear(cur int, next int) {
	bs.Clear(cur)
	bs.Set(next)
}

type World struct {
	bs     *Bitset
	width  int
	height int
	ltr    bool
}

func NewWorld(width, height int) *World {
	w := &World{
		bs:     NewBitset(width * height),
		width:  width,
		height: height,
	}
	// w.init()
	return w
}

// func (w *World) init() {
// 	w.FillCircle(10, 10, 30)
// }

func (w *World) Update() {

	for b := w.bs.Size - 1; b >= 0; b-- {
		i := b
		if !w.ltr {
			row := b / w.width
			col := b % w.width
			i = row*w.width + (w.width - 1 - col)
		}

		if b%w.width == 0 {
			w.ltr = !w.ltr
		}

		if w.bs.Test(i) {
			next := i + w.width
			row := next / w.width

			if next > w.bs.Size-1 {
				continue
			}

			if !w.bs.Test(next) {
				w.bs.Clear(i)
				w.bs.Set(next)
				continue
			}

			op := []int{-1, 1}
			fc := rand.Intn(2)
			c := op[fc]
			s := op[1-fc]

			next = i + w.width + c
			nr := next / w.width

			if next < w.bs.Size-1 && nr == row && !w.bs.Test(next) {
				w.bs.SetClear(i, next)
				continue
			}

			next = i + w.width + s
			nr = next / w.width
			if next < w.bs.Size-1 && nr == row && !w.bs.Test(next) {
				w.bs.SetClear(i, next)
				continue
			}

		}

	}

}

func (w *World) FillCircle(cx, cy, radius int) {
	r2 := radius * radius

	for dy := -radius; dy <= radius; dy++ {
		y := cy + dy
		if y < 0 || y >= w.height {
			continue // skip out-of-bounds
		}

		dx := int(math.Sqrt(float64(r2 - dy*dy))) // horizontal distance

		for x := cx - dx; x <= cx+dx; x++ {
			if x < 0 || x >= w.width {
				continue // skip out-of-bounds
			}
			if rand.Intn(2) == 1 {
				pos := y*w.width + x
				w.bs.Set(pos)
			}
		}
	}
}

func bitsOn(x uint64) []int {
	pos := make([]int, 0)

	for x != 0 {
		tz := bits.TrailingZeros64(x)
		pos = append(pos, tz)
		x &= x - 1
	}

	return pos
}

func (w *World) Draw(px []byte) {
	// for i, v := range w.bs.BitArr {
	// 	if v != 0 {
	// 		for _, j := range bitsOn(v) {
	// 			index := i*64 + j
	// 			row := index / w.width
	// 			col := index % w.width

	// px[row*col*4] = 0xC2
	// px[row*col*4+1] = 0xB2
	// px[row*col*4+2] = 0x80
	// px[row*col*4+3] = 0xFF
	// 		}
	// 	}

	// }

	for i := 0; i < w.bs.Size; i++ {
		if w.bs.Test(i) {

			x := i % w.width
			y := i / w.width

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
			px[i*4] = 0
			px[i*4+1] = 0
			px[i*4+2] = 0
			px[i*4+3] = 0
		}
	}
}

func (w *World) Fill(x int, y int) {
	w.FillCircle(x, y, 15)
}

type Game struct {
	World  *World
	pixels []byte

	// mouse position
	mx int
	my int
}

func (g *Game) Update() error {

	g.mx, g.my = ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.World.Fill(g.mx, g.my)
	}

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
