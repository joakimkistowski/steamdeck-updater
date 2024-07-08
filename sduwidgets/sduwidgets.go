package sduwidgets

import (
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

var backgroundColor = color.RGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}
var slightlyBrighterBackgroundColor = color.RGBA{R: 0x38, G: 0x38, B: 0x38, A: 0xff}

func BackgroundColor() color.Color {
	return backgroundColor
}

func NewSduText(text string, textColor color.Color, fontSize float64) *widget.Text {
	return widget.NewText(widget.TextOpts.Text(text, GetDefaultFont(fontSize), textColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)
}

func SduTextChangeColor(text *widget.Text, newColor color.Color) {
	text.Color = newColor
}

func SduTextChangeText(text *widget.Text, newText string) {
	text.Label = newText
}
