package sduwidgets

import (
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image/color"
	"strings"
)

const MaxUpdatesToShow = 16

func CreateUpdateDisplay(applicationNames []string, runtimeCount int) *widget.Container {
	updateDisplay := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(40),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{HorizontalPosition: widget.AnchorLayoutPositionCenter, VerticalPosition: widget.AnchorLayoutPositionCenter}),
		),
	)

	updateList := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(4),
				widget.GridLayoutOpts.Padding(widget.Insets{Top: 10, Left: 10, Right: 10, Bottom: 10}),
				widget.GridLayoutOpts.Spacing(20, 20),
				widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, []bool{false, false, false, false}),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(1000, 0),
		),
	)
	updateDisplay.AddChild(createApplicationCountText(len(applicationNames)))
	updateDisplay.AddChild(updateList)
	updateDisplay.AddChild(createRuntimeCountText(runtimeCount))

	applicationNameWidgets := createUpdateNames(applicationNames)
	for _, nameWidget := range applicationNameWidgets {
		updateList.AddChild(nameWidget)
	}
	return updateDisplay
}

func createApplicationCountText(applicationCount int) *widget.Text {
	applicationsNameText := "applications"
	if applicationCount == 1 {
		applicationsNameText = "application"

	}
	return NewSduText(fmt.Sprintf("Found %d %s to update:", applicationCount, applicationsNameText), colornames.White, 20)
}

func createRuntimeCountText(runtimeCount int) *widget.Text {
	runtimesNameText := "runtimes"
	if runtimeCount == 1 {
		runtimesNameText = "runtime"
	}
	return NewSduText(fmt.Sprintf("+ %d %s", runtimeCount, runtimesNameText), color.RGBA{R: 220, G: 220, B: 220, A: 255}, 18)
}

func createUpdateNames(allApplicationNames []string) []widget.PreferredSizeLocateableWidget {
	var updateNames []widget.PreferredSizeLocateableWidget
	applicationNamesToDisplay := getApplicationNamesToDisplay(allApplicationNames)
	moreUpdatesText := NewSduText("...", colornames.White, 18)
	for i, name := range applicationNamesToDisplay {
		if i < MaxUpdatesToShow-1 || (i == MaxUpdatesToShow-1 && len(applicationNamesToDisplay) == MaxUpdatesToShow) {
			nameText := NewSduText(name, colornames.White, 17)
			nameText.GetWidget().LayoutData = widget.AnchorLayoutData{HorizontalPosition: widget.AnchorLayoutPositionCenter, VerticalPosition: widget.AnchorLayoutPositionCenter}
			nameTextContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{Top: 5, Left: 5, Right: 5, Bottom: 5}))))
			nameTextContainer.BackgroundImage = image.NewNineSliceColor(slightlyBrighterBackgroundColor)
			nameTextContainer.AddChild(nameText)
			updateNames = append(updateNames, nameTextContainer)
		} else if i == MaxUpdatesToShow-1 {
			updateNames = append(updateNames, moreUpdatesText)
		} else {
			break
		}
	}
	if len(applicationNamesToDisplay) < len(allApplicationNames) && len(applicationNamesToDisplay) < MaxUpdatesToShow {
		updateNames = append(updateNames, moreUpdatesText)
	}

	return updateNames
}

func getApplicationNamesToDisplay(applicationNames []string) []string {
	namesToDisplay := make([]string, 0)
	removedNull := false
	for _, name := range applicationNames {
		if strings.Contains(name, "(null)") {
			removedNull = true
		} else {
			namesToDisplay = append(namesToDisplay, name)
		}

	}
	if len(namesToDisplay) > MaxUpdatesToShow || (len(namesToDisplay) == MaxUpdatesToShow && removedNull) {
		return namesToDisplay[:MaxUpdatesToShow-1]
	}
	return namesToDisplay
}
