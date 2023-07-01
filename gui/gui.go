package gui

import (
	"fmt"
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/quevivasbien/go-raytracing/lib"
)

const CONCURRENT bool = true
const MIN_WINDOW_WIDTH float32 = 480
const MIN_WINDOW_HEIGHT float32 = 640

func emptyScene(width, height int) Scene {
	return Scene{
		Camera:  DefaultCamera(width, height),
		Objects: []Object{},
		Lights:  []Light{},
	}
}

func createImage(s Scene) *canvas.Image {
	var img *image.RGBA
	if CONCURRENT {
		img = s.ConcurrentRender()
	} else {
		img = s.Render()
	}
	return canvas.NewImageFromImage(img)
}

func resizeWindowToScene(w fyne.Window, s Scene) {
	w.Resize(fyne.NewSize(float32(s.Camera.Width), float32(s.Camera.Height)))
}

func createImageWindow(a fyne.App, width, height float32) fyne.Window {
	w := a.NewWindow("Render")
	w.Resize(fyne.NewSize(width, height))
	return w
}

func addLightMenu(s *Scene, lightList *widget.List) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	intensityEntry := widget.NewSlider(0, 1)
	intensityEntry.Value = 0.5
	intensityEntry.Step = 0.01
	submitButton := widget.NewButton("Add Light", func() {
		s.Lights = append(s.Lights, MakeLight(coords, intensityEntry.Value))
		lightList.Refresh()
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

func showAddLightMenu(a *fyne.App, s *Scene, lightList *widget.List) {
	lightWindow := (*a).NewWindow("Add Light")
	lightWindow.SetContent(addLightMenu(s, lightList))
	lightWindow.Show()
}

func lightsContainer(a *fyne.App, s *Scene) *fyne.Container {
	var lightList *widget.List
	lightList = widget.NewList(
		func() int {
			return len(s.Lights)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Placeholder"))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			pos := s.Lights[i].Position
			obj.(*fyne.Container).Objects = []fyne.CanvasObject{
				widget.NewLabel(fmt.Sprintf("Position: (%v, %v, %v)", pos.X, pos.Y, pos.Z)),
				widget.NewLabel(fmt.Sprintf("Intensity: %v", s.Lights[i].Intensity)),
				layout.NewSpacer(),
				widget.NewButton("Remove", func() {
					s.Lights = append(s.Lights[:i], s.Lights[i+1:]...)
					lightList.Refresh()
				}),
			}
		},
	)
	label := widget.NewLabel("Lights")
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignCenter
	return container.NewVBox(
		label,
		NewStrictSized(MIN_WINDOW_WIDTH, 150, lightList),
		widget.NewButton("Add Light", func() { showAddLightMenu(a, s, lightList) }),
	)
}

func addSphereMenu(s *Scene, objectList *widget.List) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	radiusEntry := createInput("1", parseFloat)
	surface := Surface{}
	surfaceEntry := NewSurfaceEntry(&surface)
	submitButton := widget.NewButton("Add Sphere", func() {
		radius, _ := parseFloat(radiusEntry.Text)
		shape := Sphere{Center: coords, Radius: radius}
		s.Objects = append(s.Objects, Object{Shape: shape, Surface: surface})
		objectList.Refresh()
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

func addPlaneMenu(s *Scene, objectList *widget.List) *fyne.Container {
	coords := Vector{}
	coordEntry := NewVectorEntry(&coords)
	normalVec := Vector{0, 0, -1}
	normalEntry := NewVectorEntry(&normalVec)
	surface := Surface{}
	surfaceEntry := NewSurfaceEntry(&surface)
	submitButton := widget.NewButton("Add Plane", func() {
		shape := Plane{Point: coords, Norm: normalVec.Unit()}
		s.Objects = append(s.Objects, Object{Shape: shape, Surface: surface})
		objectList.Refresh()
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

func addObjectMenu(s *Scene, objectList *widget.List) *fyne.Container {
	addObjectMenu := container.NewVBox(addSphereMenu(s, objectList))
	objectTypeEntry := widget.NewSelect([]string{"Sphere", "Plane"}, func(s string) {})
	objectTypeEntry.Selected = "Sphere"
	menu := container.NewVBox(objectTypeEntry, addObjectMenu)
	objectTypeEntry.OnChanged = func(str string) {
		switch str {
		case "Sphere":
			addObjectMenu.Objects = []fyne.CanvasObject{addSphereMenu(s, objectList)}
		case "Plane":
			addObjectMenu.Objects = []fyne.CanvasObject{addPlaneMenu(s, objectList)}
		}
		addObjectMenu.Refresh()
	}
	return menu
}

// open addObjectMenu in a new window
func showAddObjectMenu(a *fyne.App, s *Scene, objectList *widget.List) {
	w := (*a).NewWindow("Add Object")
	w.SetContent(addObjectMenu(s, objectList))
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
	var objectList *widget.List
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
					objectList.Refresh()
				}),
			}
		},
	)
	label := widget.NewLabel("Objects")
	label.TextStyle.Bold = true
	return container.NewVBox(
		container.NewCenter(label),
		NewStrictSized(MIN_WINDOW_WIDTH, 250, objectList),
		widget.NewButton("Add Object", func() { showAddObjectMenu(a, s, objectList) }),
	)
}

func Launch() {
	a := app.New()
	w := a.NewWindow("Simple Raytracer")
	w.SetMaster()

	widthEntry := createInput("1920", parseUint)
	heightEntry := createInput("1080", parseUint)

	form := widget.NewForm(
		widget.NewFormItem("Width", widthEntry),
		widget.NewFormItem("Height", heightEntry),
	)

	// create button for showing the rendering
	var imageWindow fyne.Window
	scene := emptyScene(1920, 1080)
	scene.Lights = append(scene.Lights, MakeLight(Vector{0, -1, 1}, 1))
	scene.Objects = append(scene.Objects, Object{
		Shape:   Sphere{Center: Vector{0, 0, 5}, Radius: 1},
		Surface: Surface{Ambient: 0, Diffuse: 0.5, Specular: 0, Color: Vector{1, 0, 0}},
	})
	renderButton := widget.NewButton(
		"Render",
		func() {
			width, err := parseUint(widthEntry.Text)
			if err != nil {
				return
			}
			height, err := parseUint(heightEntry.Text)
			if err != nil {
				return
			}
			imageWindow = createImageWindow(a, float32(width), float32(height))
			scene.Camera = scene.Camera.Resized(int(width), int(height), math.Pi/4)
			imageWindow.SetContent(widget.NewLabel("Rendering..."))
			// render in a goroutine so the window can be shown with loading message while rendering
			go func() {
				imageWindow.SetContent(createImage(scene))
			}()
			imageWindow.Show()
		},
	)
	c := container.NewVBox(
		form,
		WhiteSpace(MIN_WINDOW_WIDTH, 50),
		lightsContainer(&a, &scene),
		WhiteSpace(MIN_WINDOW_WIDTH, 50),
		objectsContainer(&a, &scene),
		WhiteSpace(MIN_WINDOW_WIDTH, 50),
		renderButton,
	)
	w.SetContent(c)

	w.ShowAndRun()
}
