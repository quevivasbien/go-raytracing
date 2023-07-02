# go-raytracing

This is a simple demo I created while learning to use Go.

The library contains tools for rendering simple scenes with raytracing (on the CPU, but with some support for multiprocessing). The `examples` directory contains code for rendering the scenes below:

![reflective spheres](./examples/reflective-spheres/reflective-spheres.png)
![reflective planes](./examples/reflective-planes/reflective-planes.png)

## GUI

I've also added support for a GUI, using the [Fyne](https://fyne.io/) toolkit, to configure and render scenes. The `main.go` script in the main directory will launch this GUI. 
