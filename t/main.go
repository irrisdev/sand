package main

import "fmt"

func main() {

	w := 10
	h := 10
	size := w * h

	var cr int = 0
	for i := size - 1; i >= 0; i-- {

		if i/w != cr {
			fmt.Println()
			cr = i / w
		}

		// col := i % w
		//(i - w + 1) + col
		fmt.Print((i/w)*w + (w - 1 - (i % w)))
		fmt.Print(" ")

	}

}
