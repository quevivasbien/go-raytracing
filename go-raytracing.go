package main

import (
	"image/png"
	"os"

	. "github.com/quevivasbien/go-raytracing/lib"
)

func main() {
	camera := DefaultCamera(640, 480)
	light := Light{Position: Vector{X: 0, Y: 0.5, Z: 2}, Intensity: 1, Threshold: 0.01}
	sphere := Object{
		Shape:   Sphere{Center: Vector{X: 0, Y: 0, Z: 3}, Radius: 1},
		Surface: Surface{Ambient: 0.4, Diffuse: 0.7, Specular: 0.2, Color: Vector{X: 1, Y: 0, Z: 0}},
	}
	scene := Scene{Camera: camera, Objects: []Object{sphere}, Lights: []Light{light}}
	image := scene.Render()
	f, _ := os.Create("image.png")
	png.Encode(f, image)
}
