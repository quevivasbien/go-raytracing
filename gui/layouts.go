package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

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

func StrictSized(w, h float32, contents ...fyne.CanvasObject) *fyne.Container {
	return container.New(&sizedBox{w: w, h: h}, contents...)
}

func WhiteSpace(w, h float32) *fyne.Container {
	return container.New(&sizedBox{w: w, h: h})
}
