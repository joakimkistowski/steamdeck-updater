package sduwidgets

import (
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image/color"
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
	return NewSduText(fmt.Sprintf("Found %d application to update:", applicationCount), colornames.White, 19)
}

func createRuntimeCountText(runtimeCount int) *widget.Text {
	return NewSduText(fmt.Sprintf("+ %d runtimes", runtimeCount), color.RGBA{R: 220, G: 220, B: 220, A: 255}, 17)
}

func createUpdateNames(applicationNames []string) []widget.PreferredSizeLocateableWidget {
	var updateNames []widget.PreferredSizeLocateableWidget
	for i, name := range applicationNames {
		if i < MaxUpdatesToShow-1 || (i == MaxUpdatesToShow-1 && len(applicationNames) == MaxUpdatesToShow) {
			nameText := NewSduText(name, colornames.White, 17)
			nameText.GetWidget().LayoutData = widget.AnchorLayoutData{HorizontalPosition: widget.AnchorLayoutPositionCenter, VerticalPosition: widget.AnchorLayoutPositionCenter}
			nameTextContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{Top: 5, Left: 5, Right: 5, Bottom: 5}))))
			nameTextContainer.BackgroundImage = image.NewNineSliceColor(slightlyBrighterBackgroundColor)
			nameTextContainer.AddChild(nameText)
			updateNames = append(updateNames, nameTextContainer)
		} else if i == MaxUpdatesToShow-1 {
			moreUpdatesText := NewSduText("...", colornames.White, 17)
			updateNames = append(updateNames, moreUpdatesText)
		} else {
			break
		}
	}

	return updateNames
}
