package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	. "github.com/quevivasbien/go-raytracing/lib"
)

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

func NewColorEntry(v *Vector) *fyne.Container {
	rEntry := createInput("0", parseUnitRange)
	gEntry := createInput("0", parseUnitRange)
	bEntry := createInput("0", parseUnitRange)
	rEntry.OnChanged = func(s string) {
		r, _ := parseUnitRange(s)
		v.X = r
	}
	gEntry.OnChanged = func(s string) {
		g, _ := parseUnitRange(s)
		v.Y = g
	}
	bEntry.OnChanged = func(s string) {
		b, _ := parseUnitRange(s)
		v.Z = b
	}
	return container.NewHBox(
		container.NewVBox(widget.NewLabel("R"), rEntry),
		container.NewVBox(widget.NewLabel("G"), gEntry),
		container.NewVBox(widget.NewLabel("B"), bEntry),
	)
}

func NewSurfaceEntry(s *Surface) *fyne.Container {
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
	colorEntry := NewColorEntry(&s.Color)
	return container.NewHBox(
		container.NewVBox(widget.NewLabel("Ambient"), AmbientEntry),
		container.NewVBox(widget.NewLabel("Diffuse"), DiffuseEntry),
		container.NewVBox(widget.NewLabel("Specular"), SpecularEntry),
		colorEntry,
	)
}
