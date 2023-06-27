package main

import (
	"image/png"
	"os"

	. "github.com/quevivasbien/go-raytracing/lib"
)

func main() {
	camera := DefaultCamera(1920, 1080)
	light := Light{Position: Vector{-1, -4, 2}, Intensity: 1, Threshold: 0.01}
	sphere1 := Object{
		Shape:   Sphere{Center: Vector{0, 1, 3}, Radius: 0.5},
		Surface: Surface{Ambient: 0.1, Diffuse: 0.4, Specular: 0.05, Color: Vector{1, 1, 0}},
	}
	sphere2 := Object{
		Shape:   Sphere{Center: Vector{2, -1, 4}, Radius: 1},
		Surface: Surface{Ambient: 0.1, Diffuse: 0.7, Specular: 0.05, Color: Vector{1, 0, 1}},
	}
	sphere3 := Object{
		Shape:   Sphere{Center: Vector{-2, -1, 5}, Radius: 1},
		Surface: Surface{Ambient: 0.1, Diffuse: 0.7, Specular: 0.05, Color: Vector{0, 1, 1}},
	}
	sphere4 := Object{
		Shape:   Sphere{Center: Vector{0, 0, 40}, Radius: 30},
		Surface: Surface{Ambient: 0., Diffuse: 0.2, Specular: 0.9, Color: Vector{1, 1, 1}},
	}
	scene := Scene{Camera: camera, Objects: []Object{sphere1, sphere2, sphere3, sphere4}, Lights: []Light{light}}
	image := scene.ConcurrentRender()
	f, _ := os.Create("testimage.png")
	png.Encode(f, image)
}
