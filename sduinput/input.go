package sduinput

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputHandler struct {
	onAPressed  func()
	onAReleased func()
	onBPressed  func()
	onBReleased func()
}

func (in *InputHandler) RegisterAButtonHandlers(onAPressed func(), onAReleased func()) {
	in.onAPressed = onAPressed
	in.onAReleased = onAReleased
}

func (in *InputHandler) RegisterBButtonHandlers(onBPressed func(), onBReleased func()) {
	in.onBPressed = onBPressed
	in.onBReleased = onBReleased
}

func (in *InputHandler) CheckForInput() {
	in.checkKeyboardInput()
	var i ebiten.GamepadID = 0
	for ; i < 4; i++ {
		in.checkGamepadInput(i)
	}
}

func (in *InputHandler) checkKeyboardInput() {
	var justPressedKeys = make([]ebiten.Key, 0)
	justPressedKeys = inpututil.AppendJustPressedKeys(justPressedKeys)
	for _, key := range justPressedKeys {
		keyname := ebiten.KeyName(key)
		if keyname == "a" {
			in.onAPressed()
		} else if keyname == "b" {
			in.onBPressed()
		}
	}
	var justReleasedKeys = make([]ebiten.Key, 0)
	justReleasedKeys = inpututil.AppendJustReleasedKeys(justReleasedKeys)
	for _, key := range justReleasedKeys {
		keyname := ebiten.KeyName(key)
		if keyname == "a" {
			in.onAReleased()
		} else if keyname == "b" {
			in.onBReleased()
		}
	}
}

func (in *InputHandler) checkGamepadInput(gamepadId ebiten.GamepadID) {
	var justPressedButtons = make([]ebiten.GamepadButton, 0)
	justPressedButtons = inpututil.AppendJustPressedGamepadButtons(gamepadId, justPressedButtons)
	for _, gamepadButton := range justPressedButtons {
		if gamepadButton == ebiten.GamepadButton0 {
			in.onAPressed()
		} else if gamepadButton == ebiten.GamepadButton1 {
			in.onBPressed()
		}
	}
	var justReleasedButtons = make([]ebiten.GamepadButton, 0)
	justReleasedButtons = inpututil.AppendJustReleasedGamepadButtons(gamepadId, justPressedButtons)
	for _, gamepadButton := range justReleasedButtons {
		if gamepadButton == ebiten.GamepadButton0 {
			in.onAReleased()
		} else if gamepadButton == ebiten.GamepadButton1 {
			in.onBReleased()
		}
	}
}
