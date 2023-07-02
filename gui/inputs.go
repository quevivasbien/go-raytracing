package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/quevivasbien/go-raytracing/lib"
)

const SLIDER_WIDTH float32 = 100

func createInput(defaultText string, parser func(string) (float64, error)) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(defaultText)
	e.Validator = func(s string) error {
		_, err := parser(s)
		return err
	}
	return e
}

// extracts a non-zero positive integer from the given entry
func parseUint(s string) (float64, error) {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil || i <= 0 {
		return 0, fmt.Errorf("Value should be a positive integer, got %s", s)
	}
	return float64(i), nil
}

func parseUnitRange(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f < 0 || f > 1 {
		return 0, fmt.Errorf("Value should be a number between 0 and 1, got %s", s)
	}
	return f, nil
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("Value should be a number, got %s", s)
	}
	return f, nil
}

func NewVectorEntry(v *Vector) *fyne.Container {
	xEntry := createInput("0", parseFloat)
	yEntry := createInput("0", parseFloat)
	zEntry := createInput("0", parseFloat)
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
	return container.NewHBox(xEntry, yEntry, zEntry)
}

func NewUnitSlider() *widget.Slider {
	s := widget.NewSlider(0, 1)
	s.Step = 0.01
	return s
}

func NewColorSwatch(c *Vector) *canvas.Rectangle {
	rgba, _ := c.ToColor()
	rect := canvas.NewRectangle(rgba)
	rect.SetMinSize(fyne.NewSize(20, 20))
	return rect
}

func UpdateRectColor(rect *canvas.Rectangle, c *Vector) {
	rgba, _ := c.ToColor()
	rect.FillColor = rgba
	rect.Refresh()
}

func NewColorEntry(v *Vector) *fyne.Container {
	rEntry := NewUnitSlider()
	gEntry := NewUnitSlider()
	bEntry := NewUnitSlider()
	swatch := NewColorSwatch(v)
	rEntry.OnChanged = func(f float64) {
		v.X = f
		UpdateRectColor(swatch, v)
	}
	gEntry.OnChanged = func(f float64) {
		v.Y = f
		UpdateRectColor(swatch, v)
	}
	bEntry.OnChanged = func(f float64) {
		v.Z = f
		UpdateRectColor(swatch, v)
	}
	return container.NewHBox(
		NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("R"), rEntry),
		NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("G"), gEntry),
		NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("B"), bEntry),
		container.NewVBox(layout.NewSpacer(), swatch),
	)
}

func NewSurfaceEntry(s *Surface) *fyne.Container {
	AmbientEntry := NewUnitSlider()
	DiffuseEntry := NewUnitSlider()
	SpecularEntry := NewUnitSlider()
	AmbientEntry.OnChanged = func(f float64) {
		s.Ambient = f
	}
	DiffuseEntry.OnChanged = func(f float64) {
		s.Diffuse = f
	}
	SpecularEntry.OnChanged = func(f float64) {
		s.Specular = f
	}
	colorEntry := NewColorEntry(&s.Color)
	return container.NewVBox(
		container.NewHBox(
			NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("Ambient"), AmbientEntry),
			NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("Diffuse"), DiffuseEntry),
			NewStrictWidth(SLIDER_WIDTH, widget.NewLabel("Specular"), SpecularEntry),
		),
		colorEntry,
	)
}
