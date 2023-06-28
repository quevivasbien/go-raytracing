package main

import (
	"image/png"
	"os"

	. "github.com/quevivasbien/go-raytracing/lib"
)

func main() {
	camera := DefaultCamera(1920, 1080)
	light := MakeLight(Vector{-1, -4, 6}, 1)
	objects := []Object{
		// spheres in foreground
		Object{
			Shape:   Sphere{Center: Vector{1, -1, 5}, Radius: 0.8},
			Surface: Surface{Ambient: 0, Diffuse: 0.9, Specular: 1., Color: Vector{1, 1, 1}},
		},
		Object{
			Shape:   Sphere{Center: Vector{0, 0, 6}, Radius: 0.8},
			Surface: Surface{Ambient: 0, Diffuse: 0.9, Specular: 1., Color: Vector{1, 1, 1}},
		},
		Object{
			Shape:   Sphere{Center: Vector{-1, 1, 7}, Radius: 0.8},
			Surface: Surface{Ambient: 0, Diffuse: 0.9, Specular: 1., Color: Vector{1, 1, 1}},
		},
		// plane in background
		Object{
			Shape:   Plane{Norm: Vector{0.5, 0.5, -1}.Unit(), Point: Vector{0, 0, 10}},
			Surface: Surface{Ambient: 0, Diffuse: 0.5, Specular: 1, Color: Vector{0, 1, 1}},
		},
		// plane to side
		Object{
			Shape:   Plane{Norm: Vector{-0.5, -0.5, -1}.Unit(), Point: Vector{0, 0, 10}},
			Surface: Surface{Ambient: 0, Diffuse: 0.5, Specular: 1, Color: Vector{1, 0, 1}},
		},
	}

	scene := Scene{Camera: camera, Objects: objects, Lights: []Light{light}}
	image := scene.Render()
	f, _ := os.Create("reflective-planes.png")
	png.Encode(f, image)
}
