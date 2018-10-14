package main

import (
	"fmt"
	"github.com/tokenme/probab/dst"
)

func main() {
	PDF()
}

func PDF() {
	fmt.Println("test of Logistic distribution: PDF")
	fn := dst.LogisticPDF(0, 2000.33)
	t := 0.0
	for t < 10000 {
		x := fn(t) * 10000
		fmt.Println(t, "\t", x)
		t += 100
	}
}
