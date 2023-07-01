package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type StrictSized struct {
	h, w float32
}

func (s *StrictSized) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// minW, minH := s.w, float32(0)
	// for _, obj := range objects {
	// 	childSize := obj.MinSize()
	// 	minW = fyne.Max(minW, childSize.Width)
	// 	minH += childSize.Height
	// }
	// return fyne.NewSize(s.w, fyne.Max(minH, s.h))
	return fyne.NewSize(s.w, s.h)
}

func (s *StrictSized) Layout(objects []fyne.CanvasObject, size fyne.Size) {
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

func NewStrictSized(w, h float32, contents ...fyne.CanvasObject) *fyne.Container {
	return container.New(&StrictSized{w: w, h: h}, contents...)
}

func WhiteSpace(w, h float32) *fyne.Container {
	return container.New(&StrictSized{w: w, h: h})
}

type StrictWidth struct {
	w float32
}

func (s *StrictWidth) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minH := float32(0)
	for _, obj := range objects {
		childSize := obj.MinSize()
		minH += childSize.Height
	}
	return fyne.NewSize(s.w, minH)
}

func (s *StrictWidth) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, obj := range objects {
		h := obj.MinSize().Height
		obj.Resize(fyne.NewSize(size.Width, h))
		obj.Move(pos)
		pos = pos.Add(fyne.NewPos(0, h))
	}
}

func NewStrictWidth(w float32, contents ...fyne.CanvasObject) *fyne.Container {
	return container.New(&StrictWidth{w: w}, contents...)
}
