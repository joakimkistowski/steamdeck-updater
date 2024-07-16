package sduwidgets

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
	"testing"
)

func Test_NewSduText_CreatesText(t *testing.T) {
	text := NewSduText("TestLabel", colornames.Black, 12)
	assert.Equal(t, "TestLabel", text.Label)
	assert.Equal(t, colornames.Black, text.Color)
	assert.Equal(t, 12, text.Face.Metrics().Height.Ceil())
}

func Test_SduTextChangeColor_ChangesColor(t *testing.T) {
	text := NewSduText("TestLabel", colornames.Black, 12)
	SduTextChangeColor(text, colornames.White)
	assert.Equal(t, colornames.White, text.Color)
}

func Test_SduTextChangeText_ChangesText(t *testing.T) {
	text := NewSduText("TestLabel", colornames.Black, 12)
	SduTextChangeText(text, "NewLabel")
	assert.Equal(t, "NewLabel", text.Label)
}
