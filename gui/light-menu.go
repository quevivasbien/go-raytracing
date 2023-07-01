package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/quevivasbien/go-raytracing/lib"
)

func addLightMenu(s *Scene, refreshCallback func()) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	intensityEntry := widget.NewSlider(0, 1)
	intensityEntry.Value = 0.5
	intensityEntry.Step = 0.01
	submitButton := widget.NewButton("Add Light", func() {
		s.Lights = append(s.Lights, MakeLight(coords, intensityEntry.Value))
		refreshCallback()
	})
	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Position", coordEntry),
			widget.NewFormItem("Intensity", intensityEntry),
		),
		WhiteSpace(0, 10),
		submitButton,
	)
}

func showAddLightMenu(a *fyne.App, s *Scene, refreshCallback func()) {
	lightWindow := (*a).NewWindow("Add Light")
	lightWindow.SetContent(addLightMenu(s, refreshCallback))
	lightWindow.Show()
}

func lightsContainer(a *fyne.App, s *Scene) *fyne.Container {
	expander := canvas.NewRectangle(color.Transparent)
	expander.SetMinSize(fyne.NewSize(0, 0))
	var lightList *widget.List
	var outerContainer *fyne.Container
	refreshObjectList := func() {
		listHeight := (lightList.MinSize().Height + 5) * float32(lightList.Length())
		expander.SetMinSize(fyne.NewSize(0, fyne.Min(listHeight, 150)))
		lightList.Refresh()
	}
	lightList = widget.NewList(
		func() int {
			return len(s.Lights)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			pos := s.Lights[i].Position
			obj.(*fyne.Container).Objects = []fyne.CanvasObject{
				widget.NewLabel(fmt.Sprintf("Position: (%v, %v, %v)", pos.X, pos.Y, pos.Z)),
				widget.NewLabel(fmt.Sprintf("Intensity: %v", s.Lights[i].Intensity)),
				layout.NewSpacer(),
				widget.NewButton("Remove", func() {
					s.Lights = append(s.Lights[:i], s.Lights[i+1:]...)
					refreshObjectList()
				}),
			}
		},
	)
	outerContainer = container.NewBorder(nil, nil, expander, nil, lightList)
	label := widget.NewLabel("Lights")
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter
	return container.NewVBox(
		label,
		outerContainer,
		widget.NewButton("Add Light", func() { showAddLightMenu(a, s, refreshObjectList) }),
	)
}
