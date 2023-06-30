package main

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
	e.Validator = func(s string) error {
		_, err := parseUint(s)
		return err
	}
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

func addLightMenu(s *Scene, lightList *widget.List) *fyne.Container {
	validator := func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("Value should be a number, got %s", s)
		} else {
			return nil
		}
	}
	xEntry := widget.NewEntry()
	xEntry.SetText("0")
	xEntry.Validator = validator
	yEntry := widget.NewEntry()
	yEntry.SetText("0")
	yEntry.Validator = validator
	zEntry := widget.NewEntry()
	zEntry.SetText("0")
	zEntry.Validator = validator
	intensityEntry := widget.NewSlider(0, 1)
	intensityEntry.Value = 0.5
	intensityEntry.Step = 0.01
	submitButton := widget.NewButton("Add Light", func() {
		x, _ := strconv.ParseFloat(xEntry.Text, 64)
		y, _ := strconv.ParseFloat(yEntry.Text, 64)
		z, _ := strconv.ParseFloat(zEntry.Text, 64)
		s.Lights = append(s.Lights, MakeLight(Vector{x, y, z}, intensityEntry.Value))
		lightList.Refresh()
	})
	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("X", xEntry),
			widget.NewFormItem("Y", yEntry),
			widget.NewFormItem("Z", zEntry),
			widget.NewFormItem("Intensity", intensityEntry),
		),
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
		strictSized(MIN_WINDOW_WIDTH, 150, lightList),
		addLightMenu(s, lightList),
	)
}

func addObjectMenu(s *Scene, objectList *widget.List) *fyne.Container {
	return container.NewVBox(widget.NewLabel("Placeholder"))
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
		strictSized(MIN_WINDOW_WIDTH, 150, objectList),
		addObjectMenu(s, objectList),
	)
}

type sizedBox struct {
	h, w float32
}

func (s *sizedBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// minW, minH := s.w, float32(0)
	// for _, obj := range objects {
	// 	childSize := obj.MinSize()
	// 	minW = fyne.Max(minW, childSize.Width)
	// 	minH += childSize.Height
	// }
	// return fyne.NewSize(s.w, fyne.Max(minH, s.h))
	return fyne.NewSize(s.w, s.h)
}

func (s *sizedBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	// for _, obj := range objects {
	// 	h := obj.MinSize().Height
	// 	obj.Resize(fyne.NewSize(size.Width, h))
	// 	obj.Move(pos)
	// 	pos = pos.Add(fyne.NewPos(0, size.Height))
	// }
	objHeight := size.Height / float32(len(objects))
	for _, obj := range objects {
		obj.Resize(fyne.NewSize(size.Width, objHeight))
		obj.Move(pos)
		pos = pos.Add(fyne.NewPos(0, objHeight))
	}
}

func strictSized(w, h float32, contents ...fyne.CanvasObject) *fyne.Container {
	return container.New(&sizedBox{w: w, h: h}, contents...)
}

func main() {
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
			fmt.Println(scene.Lights, len(scene.Lights))
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
