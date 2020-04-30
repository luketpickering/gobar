package internal

import (
	"fmt"
)

func ColorText(s string, c string) string {
	return fmt.Sprintf("<span foreground='%v'>%v</span>", c, s)
}

func ColorBackground(s string, c string) string {
	return fmt.Sprintf("<span background='%v'>%v</span>", c, s)
}

func WhiteText(s string) string {
	return ColorText(s, "white")
}

func RedText(s string) string {
	return ColorText(s, "#dc143c")
}

func GreenText(s string) string {
	return ColorText(s, "#32cd32")
}

func BlueText(s string) string {
	return ColorText(s, "#7dacd5")
}
func OrangeBackground(s string) string {
	return ColorBackground(s, "#f29f54")
}
func DarkGreyText(s string) string {
	return ColorText(s, "#323232")
}
func OrangeText(s string) string {
	return ColorText(s, "#f29f54")
}
func BoldText(s string) string {
	return fmt.Sprintf("<span weight='ultrabold'>%v</span>", s)
}
