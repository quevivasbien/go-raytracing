package gui

import (
	"fmt"
	"image"
	"math"
	"strconv"

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

func createImageWindow(a fyne.App, width, height int) fyne.Window {
	w := a.NewWindow("Render")
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	return w
}

func createUintInput(defaultText string) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(defaultText)
	e.Validator = validateUint
	return e
}

func createFloatInput(defaultText string) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(defaultText)
	e.Validator = validateFloat
	return e
}

// extracts a non-zero positive integer from the given entry
func parseUint(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil || i <= 0 {
		return 0, fmt.Errorf("Value should be a positive integer, got %s", s)
	}
	return int(i), nil
}

func validateUint(s string) error {
	_, err := parseUint(s)
	return err
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("Value should be a number, got %s", s)
	}
	return f, nil
}

func validateFloat(s string) error {
	_, err := parseFloat(s)
	return err
}

func coordinateEntry() (*fyne.Container, *Vector) {
	xEntry := createFloatInput("0")
	yEntry := createFloatInput("0")
	zEntry := createFloatInput("0")
	v := Vector{}
	xEntry.OnChanged = func(s string) {
		x, _ := parseFloat(s)
		v.X = x
	}
	yEntry.OnChanged = func(s string) {
		y, _ := parseFloat(s)
		v.Y = y
	}
	zEntry.OnChanged = func(s string) {
		z, _ := parseFloat(s)
		v.Z = z
	}
	return container.New(layout.NewGridLayout(3), xEntry, yEntry, zEntry), &v
}

func addLightMenu(s *Scene, lightList *widget.List) *fyne.Container {
	coordEntry, coords := coordinateEntry()
	intensityEntry := widget.NewSlider(0, 1)
	intensityEntry.Value = 0.5
	intensityEntry.Step = 0.01
	submitButton := widget.NewButton("Add Light", func() {
		s.Lights = append(s.Lights, MakeLight(*coords, intensityEntry.Value))
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

func lightsContainer(s *Scene) *fyne.Container {
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
	return container.NewVBox(
		container.NewCenter(widget.NewLabel("Lights")),
		StrictSized(MIN_WINDOW_WIDTH, 150, lightList),
		addLightMenu(s, lightList),
	)
}

func surfaceEntry() (*fyne.Container, *Surface) {
	s := Surface{}
	AmbientEntry, DiffuseEntry, SpecularEntry := widget.NewSlider(0, 1), widget.NewSlider(0, 1), widget.NewSlider(0, 1)
	AmbientEntry.Value = 0
	DiffuseEntry.Value = 0
	SpecularEntry.Value = 0
	AmbientEntry.Step = 0.01
	DiffuseEntry.Step = 0.01
	SpecularEntry.Step = 0.01
	AmbientEntry.OnChanged = func(f float64) {
		s.Ambient = f
	}
	DiffuseEntry.OnChanged = func(f float64) {
		s.Diffuse = f
	}
	SpecularEntry.OnChanged = func(f float64) {
		s.Specular = f
	}
	return container.NewHBox(
		container.NewVBox(widget.NewLabel("Ambient"), AmbientEntry),
		container.NewVBox(widget.NewLabel("Diffuse"), DiffuseEntry),
		container.NewVBox(widget.NewLabel("Specular"), SpecularEntry),
	), &s
}

func addSphereMenu(s *Scene, objectList *widget.List) *fyne.Container {
	coordEntry, coords := coordinateEntry()
	radiusEntry := createFloatInput("1")
	surfaceEntry, surface := surfaceEntry()
	submitButton := widget.NewButton("Add Sphere", func() {
		radius, _ := parseFloat(radiusEntry.Text)
		shape := Sphere{Center: *coords, Radius: radius}
		s.Objects = append(s.Objects, Object{Shape: shape, Surface: *surface})
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

func addObjectMenu(s *Scene, objectList *widget.List) *fyne.Container {
	addObjectMenu := addSphereMenu(s, objectList)
	objectTypeEntry := widget.NewSelect([]string{"Sphere"}, func(s string) {})
	objectTypeEntry.Selected = "Sphere"
	menu := container.NewVBox(objectTypeEntry, addObjectMenu)
	objectTypeEntry.OnChanged = func(str string) {
		switch str {
		case "Sphere":
			addObjectMenu = addSphereMenu(s, objectList)
		default:
			addObjectMenu = nil
		}
		menu.Refresh()
	}
	return menu
}

func objectsContainer(s *Scene) *fyne.Container {
	var objectList *widget.List
	objectList = widget.NewList(
		func() int {
			return len(s.Objects)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Placeholder"))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*fyne.Container).Objects = []fyne.CanvasObject{
				widget.NewLabel(fmt.Sprint(s.Objects[i])),
				layout.NewSpacer(),
				widget.NewButton("Remove", func() {
					s.Objects = append(s.Objects[:i], s.Objects[i+1:]...)
					objectList.Refresh()
				}),
			}
		},
	)
	return container.NewVBox(
		container.NewCenter(widget.NewLabel("Objects")),
		StrictSized(MIN_WINDOW_WIDTH, 150, objectList),
		addObjectMenu(s, objectList),
	)
}

func Launch() {
	a := app.New()
	w := a.NewWindow("Simple Raytracer")
	w.SetMaster()

	widthEntry := createUintInput("1920")
	heightEntry := createUintInput("1080")

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
			imageWindow = createImageWindow(a, width, height)
			scene.Camera = scene.Camera.Resized(width, height, math.Pi/4)
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
		layout.NewSpacer(),
		lightsContainer(&scene),
		objectsContainer(&scene),
		layout.NewSpacer(),
		renderButton,
	)
	w.SetContent(c)
	w.Resize(fyne.NewSize(MIN_WINDOW_WIDTH, MIN_WINDOW_HEIGHT))

	w.ShowAndRun()
}
