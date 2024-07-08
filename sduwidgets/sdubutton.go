package sduwidgets

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
)

const SduButtonFontSize = 20

var sduTTF *truetype.Font = loadTTF()

func NewSduButton(text string, onClick func(args *widget.ButtonClickedEventArgs)) *widget.Button {
	face := GetDefaultFont(SduButtonFontSize)
	button := widget.NewButton(
		widget.ButtonOpts.Image(loadButtonImage()),
		widget.ButtonOpts.Text(text, face, &widget.ButtonTextColor{Idle: color.White, Hover: color.White, Disabled: color.RGBA{R: 100, G: 100, B: 100, A: 255}}),
		widget.ButtonOpts.TextProcessBBCode(false),
		widget.ButtonOpts.TextPosition(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.ButtonOpts.TextPadding(widget.Insets{Top: 5, Left: 20, Right: 20, Bottom: 5}),

		widget.ButtonOpts.ClickedHandler(onClick),
		// widget.ButtonOpts.DisableDefaultKeys(),
	)
	return button
}

func SduButtonReleaseAndClick(button *widget.Button) {
	button.Click()
	button.Focus(false)
}

func SduButtonPushDown(button *widget.Button) {
	button.Focus(true)
}

func SduButtonDisable(button *widget.Button) {
	button.SetState(widget.WidgetGreyed)
	button.GetWidget().Disabled = true
}

func SduButtonEnable(button *widget.Button) {
	button.SetState(widget.WidgetUnchecked)
	button.GetWidget().Disabled = false
}

func UpdateSduButtonText(button *widget.Button, newText string) {
	button.Text().Label = newText
}

func loadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.RGBA{R: 80, G: 80, B: 80, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 85, G: 85, B: 100, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func GetDefaultFont(size float64) font.Face {
	if sduTTF == nil {
		sduTTF = loadTTF()
	}
	return truetype.NewFace(sduTTF, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func loadTTF() *truetype.Font {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("Unable to load fonts: %v\n", err)
	}
	return ttfFont
}
