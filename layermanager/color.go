package layermanager

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
)

func RandomChoiceColor() color.Color {
	colors := []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
		colornames.Orange,
		colornames.Purple,
		colornames.Cyan,
		colornames.Magenta,
		colornames.Yellow,
	}
	return colors[rand.Intn(len(colors))]
}

func RandomColor() color.Color {
	r := rand.Uint32() >> 24
	g := rand.Uint32() >> 24
	b := rand.Uint32() >> 24

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
