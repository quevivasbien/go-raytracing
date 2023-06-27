package lib

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

const MAX_DEPTH uint32 = 5
const HALO_DROPOFF float64 = 1.
const HALO_THRESHOLD float64 = 0.01

type Camera struct {
	Width, Height uint32
	Position      Vector
	LookAt        unitVector
	Up            unitVector
	Right         unitVector
	FieldOfView   float64
}

func DefaultCamera(width, height uint32) Camera {
	return Camera{
		width, height,
		Vector{0, 0, 0}, unitVector{Vector{0, 0, 1}}, unitVector{Vector{0, 1, 0}}, unitVector{Vector{1, 0, 0}},
		math.Pi / 4,
	}
}

type Light struct {
	Position  Vector
	Intensity float64
	Threshold float64
}

func MakeLight(position Vector, intensity float64) Light {
	threshold := -math.Log(HALO_THRESHOLD/intensity) / HALO_DROPOFF
	return Light{position, intensity, threshold}
}

type Scene struct {
	Camera  Camera
	Objects []Object
	Lights  []Light
}

type Ray struct {
	Origin    Vector
	Direction unitVector
}

// find the first object that the ray intersects
// returns the object and the location of intersection
func (r Ray) firstIntersection(objects *[]Object) (*Object, *Vector) {
	var fi *Object
	var fiLoc *Vector
	var fiSqDist float64
	for i := range *objects {
		object := (*objects)[i] // need to do it this way so addresses remain the same
		// fmt.Printf("checking object at %v\n", object.RaySource())
		loc := object.Intersection(r)
		if loc == nil {
			continue
		}
		delta := (*loc).Sub(r.Origin)
		sqDist := delta.Dot(delta)
		// fmt.Printf("found intersection with object at %v with address %p\n", object.RaySource(), &object)
		if fi == nil || sqDist < fiSqDist {
			fi = &object
			fiLoc = loc
			fiSqDist = sqDist
		}
	}
	return fi, fiLoc
}

func (s Scene) visibleLights(p Vector) []*Light {
	var visible []*Light
	for _, light := range s.Lights {
		r := Ray{Origin: light.Position, Direction: p.Sub(light.Position).Unit()}
		fi, _ := r.firstIntersection(&s.Objects)
		if fi != nil {
			visible = append(visible, &light)
		}
	}
	return visible
}

func (s Scene) checkForLight(r Ray) Vector {
	out := Zero()
	for _, light := range s.Lights {
		// get scalar projection of ray connecting light and camera onto ray
		// if this is negative, the light is behind the camera
		cameraToLight := light.Position.Sub(r.Origin)
		scalarProjection := cameraToLight.Dot(r.Direction.Vector)
		if scalarProjection < 0 {
			continue
		}
		// get distance between projected ray and light; use this to calculate light intensity
		distance := math.Sqrt(cameraToLight.Dot(cameraToLight) - scalarProjection*scalarProjection)
		if distance > light.Threshold {
			continue
		}
		intensity := light.Intensity * math.Exp(-HALO_DROPOFF*distance)
		out = out.Add(White().MulScalar(intensity))
	}
	return out
}

func (r Ray) interact(o *Object, loc *Vector, s *Scene, depth uint32) Vector {
	normal := (*loc).Sub(o.RaySource()).Unit().Vector
	color := o.Surface.Color.MulScalar(o.Surface.Ambient)
	if o.Surface.Diffuse > 0 {
		diffusion := 0.
		for _, light := range s.visibleLights(*loc) {
			contribution := light.Intensity * normal.Dot(light.Position.Sub(*loc).Unit().Vector)
			if contribution > 0 {
				diffusion += contribution
			}
		}
		if diffusion > 1 {
			diffusion = 1
		}
		diffuseColor := o.Surface.Color.MulScalar(diffusion * o.Surface.Diffuse)
		color = color.Add(diffuseColor)
	}
	if o.Surface.Specular > 0 {
		reflection := r.Direction.Reflect(normal).Unit()
		reflectionRay := Ray{Origin: *loc, Direction: reflection}
		reflectedColor := s.trace(reflectionRay, depth+1)
		specularColor := reflectedColor.MulScalar(o.Surface.Specular)
		color = color.Add(specularColor)
	}

	return color
}

func (s Scene) trace(r Ray, depth uint32) Vector {
	if depth > MAX_DEPTH {
		return Zero()
	}
	fi, fiLoc := r.firstIntersection(&s.Objects)
	if fi == nil {
		return s.checkForLight(r)
	}
	// fmt.Printf("Intersected with %v at %v\n", fi, fiLoc)
	return r.interact(fi, fiLoc, &s, depth)
}

func (s Scene) Render() *image.RGBA {
	camera := s.Camera
	hwRatio := float64(camera.Height) / float64(camera.Width)
	halfWidth := math.Tan(camera.FieldOfView)
	halfHeight := hwRatio * halfWidth
	pixelWidth := 2 * halfWidth / float64(camera.Width-1)
	pixelHeight := 2 * halfHeight / float64(camera.Height-1)
	img := image.NewRGBA(image.Rect(0, 0, int(camera.Width), int(camera.Height)))
	for x := uint32(0); x < camera.Width; x++ {
		for y := uint32(0); y < camera.Height; y++ {
			// create ray looking at pixel
			xComp := camera.Right.MulScalar(float64(x)*pixelWidth - halfWidth)
			yComp := camera.Up.MulScalar(float64(y)*pixelHeight - halfHeight)
			direction := camera.LookAt.Add(xComp).Add(yComp).Unit()
			ray := Ray{Origin: camera.Position, Direction: direction}
			// trace ray
			traceResult := s.trace(ray, 0).Trim(0, 1)
			pixel, error := traceResult.ToColor()
			if pixel == nil {
				fmt.Printf("error when unpacking color: %s\n", error)
				pixel = &color.RGBA{0, 0, 0, 255}
			}
			img.Set(int(x), int(y), pixel)
		}
	}
	return img
}
