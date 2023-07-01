package gui

import (
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
