package main

import (
	"fmt"
	pgu "github.com/luketpickering/gobar/pangoutils"
)

func main() {
	my := pgu.NewPangoStrU("hello")

	fmt.Println(my)

	my.SetFGColor(pgu.Red)

	fmt.Println(my)

	my.SetBGColor(pgu.Red)

	fmt.Println(my)

	my.SetFontWeight(pgu.Bold)

	fmt.Println(my)
}