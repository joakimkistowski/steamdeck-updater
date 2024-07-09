package main

import (
	"fmt"
	"log"
	"os"
	"steamdeckupdater/flatpakintegration"
	"steamdeckupdater/sduinput"
	"steamdeckupdater/sduwidgets"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const DefaultWindowWidth = 1280
const DefaultWindowHeight = 800

const UpdateStateContainerHeight = DefaultWindowHeight * 0.9
const ButtonContainerHeight = DefaultWindowHeight * 0.1

const UpgradingMessage = "Upgrading"

type SteamdeckUpdateStatus int

const (
	Empty          SteamdeckUpdateStatus = iota
	Refreshing     SteamdeckUpdateStatus = iota
	Refreshed      SteamdeckUpdateStatus = iota
	GettingUpdates SteamdeckUpdateStatus = iota
	GottenUpdates  SteamdeckUpdateStatus = iota
	Updating       SteamdeckUpdateStatus = iota
	Done           SteamdeckUpdateStatus = iota
)

type GraphicalSteamdeckUpdaterApp struct {
	UpdateStatus               SteamdeckUpdateStatus
	Updates                    *flatpakintegration.SDUFlatpakContainerList
	UpgradeResult              *flatpakintegration.SDUFlatpakUpgradeResult
	inputHandler               *sduinput.InputHandler
	updateStatusChangeNotifier chan SteamdeckUpdateStatus
	okButton                   *widget.Button
	cancelButton               *widget.Button
	updateStateContainer       *widget.Container
	displayMessageText         *widget.Text
	ui                         *ebitenui.UI
}

func (app *GraphicalSteamdeckUpdaterApp) Update() error {
	app.checkForStateUpdatesNonBlocking()
	// called on the very first Update cycle only
	if app.UpdateStatus == Empty {
		app.triggerRemotesUpdateAsync()
	} else if app.UpdateStatus == Updating {
		app.updateRegularMessage(getCurrentUpgradingMessage())
	}
	app.inputHandler.CheckForInput()
	app.ui.Update()
	return nil
}

func (app *GraphicalSteamdeckUpdaterApp) Draw(screen *ebiten.Image) {
	screen.Fill(sduwidgets.BackgroundColor())
	app.ui.Draw(screen)
}

func (app *GraphicalSteamdeckUpdaterApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return DefaultWindowWidth, DefaultWindowHeight
}

func (app *GraphicalSteamdeckUpdaterApp) initializeUi() {
	uiContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
	)
	app.updateStateContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{Top: 10, Left: 10, Right: 10, Bottom: 10}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(DefaultWindowWidth, UpdateStateContainerHeight)),
	)
	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Padding(widget.Insets{Top: 10, Left: 10, Right: 10, Bottom: 10}),
				widget.GridLayoutOpts.Spacing(DefaultWindowWidth*0.7, 10),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
			),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(DefaultWindowWidth, ButtonContainerHeight)),
	)
	app.displayMessageText = sduwidgets.NewSduText("Loading ...", colornames.Black, 40)
	app.displayMessageText.GetWidget().LayoutData = widget.AnchorLayoutData{HorizontalPosition: widget.AnchorLayoutPositionCenter, VerticalPosition: widget.AnchorLayoutPositionCenter}
	app.updateStateContainer.AddChild(app.displayMessageText)

	app.cancelButton = sduwidgets.NewSduButton("(B) Cancel / Close", func(args *widget.ButtonClickedEventArgs) {
		os.Exit(0)
	})
	app.okButton = sduwidgets.NewSduButton("(A) Update", func(args *widget.ButtonClickedEventArgs) {
		if app.UpdateStatus == GottenUpdates {
			// && (len(app.Updates.ApplicationPackages) > 0 || len(app.Updates.RuntimePackages) > 0
			app.triggerFullUpgrade()
		}
	})
	sduwidgets.SduButtonDisable(app.okButton)
	uiContainer.AddChild(app.updateStateContainer)
	uiContainer.AddChild(buttonContainer)
	buttonContainer.AddChild(app.cancelButton)
	buttonContainer.AddChild(app.okButton)

	app.ui = &ebitenui.UI{
		Container: uiContainer,
	}
}

func (app *GraphicalSteamdeckUpdaterApp) registerInputHandlers() {
	app.inputHandler.RegisterAButtonHandlers(app.onAPressed, app.onAReleased)
	app.inputHandler.RegisterBButtonHandlers(app.onBPressed, app.onBReleased)
}

func main() {
	ebiten.SetWindowSize(DefaultWindowWidth, DefaultWindowHeight)
	ebiten.SetWindowTitle("Steam Deck Updater")

	app := &GraphicalSteamdeckUpdaterApp{updateStatusChangeNotifier: make(chan SteamdeckUpdateStatus), inputHandler: &sduinput.InputHandler{}}
	app.initializeUi()
	app.registerInputHandlers()
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}

func (app *GraphicalSteamdeckUpdaterApp) triggerRemotesUpdateAsync() {
	app.performStateChange(Refreshing)
	go func() {
		flatpakintegration.UpdateFlatpakRemotes()
		app.updateStatusChangeNotifier <- Refreshed
	}()
}

func (app *GraphicalSteamdeckUpdaterApp) triggerGettingUpdates() {
	app.performStateChange(GettingUpdates)
	go func() {
		app.Updates = flatpakintegration.GetUpdateableFlatpaks()
		app.updateStatusChangeNotifier <- GottenUpdates
	}()
}

func (app *GraphicalSteamdeckUpdaterApp) triggerFullUpgrade() {
	app.performStateChange(Updating)
	go func() {
		app.UpgradeResult = flatpakintegration.UpgradeAllFlatpaks()
		app.updateStatusChangeNotifier <- Done
	}()

}

func (app *GraphicalSteamdeckUpdaterApp) checkForStateUpdatesNonBlocking() {
	select {
	case newState := <-app.updateStatusChangeNotifier:
		app.performStateChange(newState)
	default:
	}
}

func (app *GraphicalSteamdeckUpdaterApp) performStateChange(newState SteamdeckUpdateStatus) {
	app.UpdateStatus = newState
	switch newState {
	case Refreshing:
		app.configureDefaultButtonStates()
		app.updateRegularMessage("Refreshing update cache")
	case Refreshed:
		app.configureDefaultButtonStates()
		app.updateRegularMessage("Cache refreshed")
		app.triggerGettingUpdates()
	case GettingUpdates:
		app.configureDefaultButtonStates()
		app.updateRegularMessage("Getting updates")
	case GottenUpdates:
		sduwidgets.SduButtonEnable(app.cancelButton)
		if app.Updates.IsErrorGettingPackages {
			sduwidgets.SduButtonDisable(app.okButton)
			app.updateErrorMessage("Error getting updates. Is a Flatpak user installation present?")
		} else {
			if len(app.Updates.ApplicationPackages) == 0 && len(app.Updates.RuntimePackages) == 0 {
				app.configureDefaultButtonStates()
				app.updateRegularMessage("No updates found")
			} else if len(app.Updates.ApplicationPackages) == 0 {
				app.updateRegularMessage(fmt.Sprintf("%d runtime updates found", len(app.Updates.RuntimePackages)))
				sduwidgets.SduButtonEnable(app.okButton)
			} else {
				sduwidgets.SduButtonEnable(app.okButton)
				app.updateStateContainer.RemoveChildren()
				app.updateStateContainer.AddChild(sduwidgets.CreateUpdateDisplay(app.getUpdateApplicationNames(), len(app.Updates.RuntimePackages)))
				//app.updateStateContainer.AddChild(sduwidgets.CreateUpdateDisplay([]string{"one", "two", "three", "four", "one", "two", "three", "four", "one", "two", "three", "four", "one", "two", "three", "four", "asd"}, 1))
			}
		}
	case Updating:
		app.updateRegularMessage("Upgrading ...")
		sduwidgets.SduButtonDisable(app.cancelButton)
		sduwidgets.SduButtonDisable(app.okButton)
	case Done:
		app.configureDefaultButtonStates()
		if app.UpgradeResult.IsErrorUpgrading {
			app.updateErrorMessage("Error upgrading: " + app.UpgradeResult.Message)
		} else {
			app.updateRegularMessage("Upgrade successful!")
		}
	default:
		app.configureDefaultButtonStates()
	}
}

func (app *GraphicalSteamdeckUpdaterApp) updateRegularMessage(message string) {
	app.updateStateContainer.RemoveChildren()
	app.updateStateContainer.AddChild(app.displayMessageText)
	sduwidgets.SduTextChangeText(app.displayMessageText, message)
	sduwidgets.SduTextChangeColor(app.displayMessageText, colornames.White)
}

func (app *GraphicalSteamdeckUpdaterApp) updateErrorMessage(message string) {
	app.updateRegularMessage(message)
	sduwidgets.SduTextChangeColor(app.displayMessageText, colornames.Red)
}

func (app *GraphicalSteamdeckUpdaterApp) configureDefaultButtonStates() {
	sduwidgets.SduButtonDisable(app.okButton)
	sduwidgets.SduButtonEnable(app.cancelButton)
}

func getCurrentUpgradingMessage() string {
	timeInTwoSeconds := time.Now().UnixMilli() % 2000
	if timeInTwoSeconds < 666 {
		return UpgradingMessage + " ."
	} else if timeInTwoSeconds < 1333 {
		return UpgradingMessage + " .."
	} else {
		return UpgradingMessage + " ..."
	}
}

func (app *GraphicalSteamdeckUpdaterApp) getUpdateApplicationNames() []string {
	applicationNames := make([]string, 0)
	for _, name := range app.Updates.ApplicationPackages {
		applicationNames = append(applicationNames, name.Name)
	}
	return applicationNames
}

func (app *GraphicalSteamdeckUpdaterApp) onAPressed() {
	if app.okButton != nil && !app.okButton.GetWidget().Disabled {
		sduwidgets.SduButtonPushDown(app.okButton)
	}
}

func (app *GraphicalSteamdeckUpdaterApp) onAReleased() {
	if app.okButton != nil && !app.okButton.GetWidget().Disabled {
		sduwidgets.SduButtonReleaseAndClick(app.okButton)
	}
}

func (app *GraphicalSteamdeckUpdaterApp) onBPressed() {
	if app.cancelButton != nil && !app.cancelButton.GetWidget().Disabled {
		sduwidgets.SduButtonPushDown(app.cancelButton)
	}
}

func (app *GraphicalSteamdeckUpdaterApp) onBReleased() {
	if app.cancelButton != nil && !app.cancelButton.GetWidget().Disabled {
		sduwidgets.SduButtonReleaseAndClick(app.cancelButton)
	}
}
