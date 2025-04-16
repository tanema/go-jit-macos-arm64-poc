package main

import (
	_ "embed"
	"fmt"
)

//go:embed build/write_flat.bin
var write []byte

func main() {
	write = write[:24]
	fmt.Println("var code = []byte{")
	for i := 0; i <= len(write)-1; i += 4 {
		fmt.Printf("\t0x%02X, 0x%02X, 0x%02X, 0x%02X,\n", write[i], write[i+1], write[i+2], write[i+3])
	}
	fmt.Println("}")
}
