package sduwidgets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewSduButton_CreatesAButtonWithLabelText(t *testing.T) {
	button := NewSduButton("TestLabel", nil)
	assert.Equal(t, "TestLabel", button.Text().Label)
}

func Test_UpdateSduButtonText_UpdatesLabelText(t *testing.T) {
	button := NewSduButton("TestLabel", nil)
	UpdateSduButtonText(button, "NewLabel")
	assert.Equal(t, "NewLabel", button.Text().Label)
}

func Test_SduButtonDisable_DisablesButton(t *testing.T) {
	button := NewSduButton("TestLabel", nil)
	SduButtonDisable(button)
	assert.True(t, button.GetWidget().Disabled)
}

func Test_SduButtonEnable_EnablesButton(t *testing.T) {
	button := NewSduButton("TestLabel", nil)
	SduButtonDisable(button)
	SduButtonEnable(button)
	assert.False(t, button.GetWidget().Disabled)
}
