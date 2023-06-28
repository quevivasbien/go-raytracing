package main

import (
	"image/png"
	"os"

	. "github.com/quevivasbien/go-raytracing/lib"
)

func main() {
	camera := DefaultCamera(1920, 1080)
	light := MakeLight(Vector{-1, -4, 2}, 1)
	objects := []Object{
		// large, nonreflective white sphere at back
		Object{
			Shape:   Sphere{Center: Vector{0, 0, 1000}, Radius: 700},
			Surface: Surface{Ambient: 0., Diffuse: 0.9, Specular: 0., Color: Vector{1, 1, 1}},
		},
		// // large, reflective white sphere at back
		Object{
			Shape:   Sphere{Center: Vector{5, 4, 40}, Radius: 20},
			Surface: Surface{Ambient: 0., Diffuse: 0.1, Specular: 0.9, Color: Vector{1, 1, 1}},
		},
		// yellow sphere at bottom
		Object{
			Shape:   Sphere{Center: Vector{0, 1.5, 3}, Radius: 0.5},
			Surface: Surface{Ambient: 0.1, Diffuse: 0.4, Specular: 0., Color: Vector{1, 1, 0}},
		},
		// magenta sphere at top right
		Object{
			Shape:   Sphere{Center: Vector{2, -1, 4}, Radius: 0.5},
			Surface: Surface{Ambient: 0.1, Diffuse: 0.7, Specular: 0., Color: Vector{1, 0, 1}},
		},
		// cyan sphere at top left
		Object{
			Shape:   Sphere{Center: Vector{-2, -1.2, 5}, Radius: 1},
			Surface: Surface{Ambient: 0.1, Diffuse: 0.7, Specular: 0., Color: Vector{0, 1, 1}},
		},
		// green sphere at right
		Object{
			Shape:   Sphere{Center: Vector{4, 0, 0}, Radius: 3.5},
			Surface: Surface{Ambient: 0., Diffuse: 0.2, Specular: 1., Color: Vector{0, 1, 0}},
		},
		// blue sphere at left
		Object{
			Shape:   Sphere{Center: Vector{-6, 0, 12}, Radius: 5},
			Surface: Surface{Ambient: 0., Diffuse: 0.1, Specular: 1., Color: Vector{0, 0, 1}},
		},
	}

	scene := Scene{Camera: camera, Objects: objects, Lights: []Light{light}}
	image := scene.ConcurrentRender()
	f, _ := os.Create("reflective-spheres.png")
	png.Encode(f, image)
}
