package main

import (
	"fmt"
	"math/bits"
	"strings"
)

/*
	Set bit
	bitset ^= (1 << 0)

	Extract bit
	bit := (bitset >> 2) & 1

	Test Bit
	bitset&(1<<pos) != 0

*/

// func main() {
// 	bitset := uint64(0)

// 	bitset ^= (1 << 0)
// 	bitset ^= (1 << 2)
// 	bitset ^= (1 << 10)

// 	bit := (bitset >> 2) & 1
// 	fmt.Println(bin64(bit))

// 	fmt.Println(bitset&(1<<3) != 0)

// 	fmt.Println(bin64(bitset))
// }

const w = 20
const h = 10

func main() {

	bitset := make([]uint64, w*h)

	// totalIterations := (len(bitset) * 64) - 1

	bitset[3] = uint64(1) << 2
	fmt.Println(bin64(bitset[3]))

	// for i := totalIterations; i >= 0; i-- {
	// 	scol := i / 64 // points to uint64 in the bitset slice
	// 	// srow := uint(i % 64) // points to the bit inside uint64

	// 	// extract bitset[index] -> srow position
	// 	// t := (bitset[scol] >> srow) & 1

	// 	// t := bitset[scol]&(1<<srow) != 0

	// 	// if t {
	// 	// 	fmt.Println(bin64((bitset[scol] << srow) & 1))
	// 	// }

	// 	if bitset[scol] != 0 {
	// 		fmt.Println("accessing ", scol)
	// 		// fmt.Println(bin64((bitset[scol] << srow) & 1))
	// 		fmt.Println(bitset[scol])
	// 	}

	// }

	// for idx := len(bitset) - 1; idx >= 0; idx-- {
	// 	x := bitset[idx]
	// 	if x != 0 {
	// 		pos := bitsOn(x)
	// 		fmt.Println(idx)
	// 		fmt.Println(pos)
	// 	}

	// }

	// for idx := w*h - 1; idx >= 0; idx-- {
	// 	scol := idx % 64
	// 	srow := uint(idx % 64)

	// 	// create mask of with // 10
	// 	if w <= 64 {
	// 		mask := (uint64(1) << w) - 1
	// 		x := bitset[scol] & mask
	// 		pos := bitsOn(x)
	// 	}

	// }

}

func bitsOn(x uint64) []int {

	pos := []int{}

	for x != 0 {
		tz := bits.TrailingZeros64(x)
		pos = append(pos, tz)
		x &= x - 1
	}

	return pos
}

func bin64(x uint64) string {
	s := fmt.Sprintf("%064b", x)
	var out strings.Builder
	for i, c := range s {
		if i > 0 && i%8 == 0 {
			out.WriteByte(' ')
		}
		out.WriteRune(c)
	}
	return out.String()
}
