package lib

import "image/color"

type Job struct {
	Scene      *Scene
	Start, End int
}

type JobResult struct {
	Pixels []*color.RGBA
	Start  int
}

func (j Job) Run(c chan JobResult) {
	pixels := make([]*color.RGBA, 0, j.End-j.Start)
	for i := j.Start; i < j.End; i++ {
		x := i % j.Scene.Camera.Width
		y := i / j.Scene.Camera.Width
		pixel := j.Scene.RenderPixel(x, y)
		pixels = append(pixels, pixel)
	}
	c <- JobResult{Pixels: pixels, Start: j.Start}
}
