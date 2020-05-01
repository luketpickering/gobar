package pangoutils

import (
	"fmt"
	"strings"
	"image/color"
)

const (
	Light = "light"
	Normal = "normal"
	Bold = "bold"
	Ultrabold = "ultrabold"
)

//Color 'constants'
var White color.RGBA = color.RGBA{0xFF,0xFF,0xFF,0xFF}
var Black color.RGBA = color.RGBA{0x00,0x00,0x00,0x00}
var Red color.RGBA = color.RGBA{0xDC,0x14,0x3C,0x00}
var Green color.RGBA = color.RGBA{0x32,0xCD,0x32,0x00}
var Blue color.RGBA = color.RGBA{0x7D,0xAC,0xD5,0x00}
var Orange color.RGBA = color.RGBA{0xF2,0x9F,0x54,0x00}
var DarkGrey color.RGBA = color.RGBA{0x32,0x14,0x3C,0x00}

type PangoStrUnit struct {
	s string
	fg color.RGBA
	fg_s bool
	bg color.RGBA
	bg_s bool
	weight string
	weight_s bool
}

func MakePangoStrU(s string) PangoStrUnit {
	return PangoStrUnit{s,color.RGBA{0xFF,0xFF,0xFF,0xFF}, false, color.RGBA{0xFF,0xFF,0xFF,0xFF}, false, Normal, false}
}

func (p PangoStrUnit) SetFGColor(c color.RGBA) PangoStrUnit {
	p.fg = c
	p.fg_s = true
	return p
}

func (p PangoStrUnit) SetBGColor(c color.RGBA) PangoStrUnit {
	p.bg = c
	p.bg_s = true
	return p
}

func (p PangoStrUnit) SetFontWeight(w string) PangoStrUnit {
	p.weight = w
	p.weight_s = true
	return p
}

func (p PangoStrUnit) String() string {
	strbuilder := []string{"<span"}

	if p.fg_s{
		strbuilder = append(strbuilder,
		  fmt.Sprintf(" foreground='#%.2x%.2x%.2x'", p.fg.R, p.fg.G, p.fg.B))
	}
	if p.bg_s{
		strbuilder = append(strbuilder,
		  fmt.Sprintf(" background='#%.2x%.2x%.2x'", p.bg.R, p.bg.G, p.bg.B))
	}
	if p.weight_s {
		strbuilder = append(strbuilder, " weight='", p.weight, "'")
	}
	strbuilder = append(strbuilder, ">",p.s,"</span>")

	return strings.Join(strbuilder,"")
}


type PangoStr []PangoStrUnit

func (ps PangoStr) String() string {
	var strbuilder []string

	for _, v := range ps {
		strbuilder = append(strbuilder, v.String())
	}
	
	return strings.Join(strbuilder,"")
}