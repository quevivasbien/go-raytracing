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

func addSphereMenu(s *Scene, refreshCallback func()) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	radiusEntry := createInput("1", parseFloat)
	surface := Surface{}
	surfaceEntry := NewSurfaceEntry(&surface)
	submitButton := widget.NewButton("Add Sphere", func() {
		radius, _ := parseFloat(radiusEntry.Text)
		shape := Sphere{Center: coords, Radius: radius}
		s.Objects = append(s.Objects, Object{Shape: shape, Surface: surface})
		refreshCallback()
	})
	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Position", coordEntry),
			widget.NewFormItem("Radius", radiusEntry),
			widget.NewFormItem("Surface", surfaceEntry),
		),
		WhiteSpace(0, 10),
		submitButton,
	)
}

func addPlaneMenu(s *Scene, refreshCallback func()) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	normalVec := Vector{0, 0, -1}
	normalEntry := NewVectorEntry(&normalVec)
	surface := Surface{}
	surfaceEntry := NewSurfaceEntry(&surface)
	submitButton := widget.NewButton("Add Plane", func() {
		shape := Plane{Point: coords, Norm: normalVec.Unit()}
		s.Objects = append(s.Objects, Object{Shape: shape, Surface: surface})
		refreshCallback()
	})
	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Normal vector", normalEntry),
			widget.NewFormItem("Point", coordEntry),
			widget.NewFormItem("Surface", surfaceEntry),
		),
		WhiteSpace(0, 10),
		submitButton,
	)
}

func addObjectMenu(s *Scene, refreshCallback func()) *fyne.Container {
	addObjectMenu := container.NewVBox(addSphereMenu(s, refreshCallback))
	objectTypeEntry := widget.NewSelect([]string{"Sphere", "Plane"}, func(s string) {})
	objectTypeEntry.Selected = "Sphere"
	menu := container.NewVBox(objectTypeEntry, addObjectMenu)
	objectTypeEntry.OnChanged = func(str string) {
		switch str {
		case "Sphere":
			addObjectMenu.Objects = []fyne.CanvasObject{addSphereMenu(s, refreshCallback)}
		case "Plane":
			addObjectMenu.Objects = []fyne.CanvasObject{addPlaneMenu(s, refreshCallback)}
		}
		addObjectMenu.Refresh()
	}
	return menu
}

// open addObjectMenu in a new window
func showAddObjectMenu(a *fyne.App, s *Scene, refreshCallback func()) {
	w := (*a).NewWindow("Add Object")
	w.SetContent(addObjectMenu(s, refreshCallback))
	w.Show()
}

func sphereInfo(s Sphere) *fyne.Container {
	label := widget.NewLabel("Sphere")
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter
	return container.NewHBox(
		label,
		widget.NewLabel(fmt.Sprintf("Center: %v", s.Center)),
		widget.NewLabel(fmt.Sprintf("Radius: %v", s.Radius)),
	)
}

func planeInfo(p Plane) *fyne.Container {
	label := widget.NewLabel("Plane")
	label.TextStyle.Bold = true
	return container.NewHBox(
		label,
		widget.NewLabel(fmt.Sprintf("Point: %v", p.Point)),
		widget.NewLabel(fmt.Sprintf("Normal: %v", p.Norm)),
	)
}

func shapeInfo(s Shape) *fyne.Container {
	sphere, ok := s.(Sphere)
	if ok {
		return sphereInfo(sphere)
	}
	plane, ok := s.(Plane)
	if ok {
		return planeInfo(plane)
	}
	return container.NewHBox(
		widget.NewLabel(fmt.Sprintf("Shape: %v", s)),
	)
}

func surfaceInfo(s Surface) *fyne.Container {
	return container.NewHBox(
		container.NewHBox(widget.NewLabel("Color"), container.NewVBox(layout.NewSpacer(), NewColorSwatch(&s.Color), layout.NewSpacer())),
		widget.NewLabel(fmt.Sprintf("Ambient: %.2f", s.Ambient)),
		widget.NewLabel(fmt.Sprintf("Diffuse: %.2f", s.Diffuse)),
		widget.NewLabel(fmt.Sprintf("Specular: %.2f", s.Specular)),
	)
}

func objectInfo(o Object) *fyne.Container {
	return container.NewVBox(
		shapeInfo(o.Shape),
		surfaceInfo(o.Surface),
	)
}

func objectsContainer(a *fyne.App, s *Scene) *fyne.Container {
	expander := canvas.NewRectangle(color.Transparent)
	expander.SetMinSize(fyne.NewSize(0, 0))
	var objectList *widget.List
	var outerContainer *fyne.Container
	refreshObjectList := func() {
		listHeight := (objectList.MinSize().Height + 5) * float32(objectList.Length())
		expander.SetMinSize(fyne.NewSize(0, fyne.Min(listHeight, 250)))
		objectList.Refresh()
	}
	objectList = widget.NewList(
		func() int {
			return len(s.Objects)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(container.NewVBox(widget.NewLabel(""), widget.NewLabel("")))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*fyne.Container).Objects = []fyne.CanvasObject{
				objectInfo(s.Objects[i]),
				layout.NewSpacer(),
				widget.NewButton("Remove", func() {
					s.Objects = append(s.Objects[:i], s.Objects[i+1:]...)
					refreshObjectList()
				}),
			}
		},
	)
	outerContainer = container.NewBorder(nil, nil, expander, nil, objectList)
	label := widget.NewLabel("Objects")
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter
	return container.NewVBox(
		label,
		outerContainer,
		widget.NewButton("Add Object", func() { showAddObjectMenu(a, s, refreshObjectList) }),
	)
}
