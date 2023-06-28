package lib

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"time"
)

const MAX_DEPTH int = 5
const HALO_DROPOFF float64 = 1.
const HALO_THRESHOLD float64 = 0.01

type Camera struct {
	Width, Height                                  int
	Position                                       Vector
	LookAt, Up, Right                              unitVector
	HalfWidth, HalfHeight, PixelWidth, PixelHeight float64
}

func MakeCamera(
	Width, Height int,
	Position Vector,
	LookAt, Up, Right unitVector,
	FieldOfView float64,
) Camera {
	hwRatio := float64(Height) / float64(Width)
	halfWidth := math.Tan(FieldOfView)
	halfHeight := hwRatio * halfWidth
	pixelWidth := 2 * halfWidth / float64(Width-1)
	pixelHeight := 2 * halfHeight / float64(Height-1)
	return Camera{
		Width, Height,
		Position,
		LookAt, Up, Right,
		halfWidth, halfHeight, pixelWidth, pixelHeight,
	}
}

func DefaultCamera(width, height int) Camera {
	return MakeCamera(
		width, height,
		Vector{0, 0, 0},
		unitVector{Vector{0, 0, 1}}, unitVector{Vector{0, 1, 0}}, unitVector{Vector{1, 0, 0}},
		math.Pi/4,
	)
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

func (r Ray) interact(o *Object, loc *Vector, s *Scene, depth int) Vector {
	normal := o.Normal(*loc)
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
		reflection := r.Direction.Reflect(normal.Vector).Unit()
		reflectionRay := Ray{Origin: *loc, Direction: reflection}
		reflectedColor := s.trace(reflectionRay, depth+1)
		specularColor := reflectedColor.MulScalar(o.Surface.Specular)
		color = color.Add(specularColor)
	}

	return color
}

func (s Scene) trace(r Ray, depth int) Vector {
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

func (s Scene) RenderPixel(x int, y int) *color.RGBA {
	// create ray looking at pixel
	xComp := s.Camera.Right.MulScalar(float64(x)*s.Camera.PixelWidth - s.Camera.HalfWidth)
	yComp := s.Camera.Up.MulScalar(float64(y)*s.Camera.PixelHeight - s.Camera.HalfHeight)
	direction := s.Camera.LookAt.Add(xComp).Add(yComp).Unit()
	ray := Ray{Origin: s.Camera.Position, Direction: direction}
	// trace ray
	traceResult := s.trace(ray, 0).Trim(0, 1)
	pixel, error := traceResult.ToColor()
	if pixel == nil {
		fmt.Printf("error when unpacking color: %s\n", error)
		pixel = &color.RGBA{0, 0, 0, 255}
	}
	return pixel
}

func (s Scene) Render() *image.RGBA {
	camera := s.Camera
	img := image.NewRGBA(image.Rect(0, 0, camera.Width, camera.Height))
	timeStart := time.Now()
	for x := 0; x < camera.Width; x++ {
		for y := 0; y < camera.Height; y++ {
			pixel := s.RenderPixel(x, y)
			img.Set(x, y, pixel)
		}
	}
	timeEnd := time.Now()
	fmt.Printf("Rendered in %v\n", timeEnd.Sub(timeStart))
	return img
}

// same functionality as Render, but works in parallel, using all available CPU cores
func (s Scene) ConcurrentRender() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, s.Camera.Width, s.Camera.Height))
	nCores := runtime.NumCPU()
	timeStart := time.Now()
	// render chunks in parallel
	jobSize := s.Camera.Height * s.Camera.Width / nCores
	jobs := make([]Job, 0, nCores)
	for i := 0; i < nCores; i++ {
		start := i * jobSize
		end := start + jobSize
		if i == nCores-1 {
			end = s.Camera.Height * s.Camera.Width
		}
		jobs = append(jobs, Job{Scene: &s, Start: start, End: end})
	}
	c := make(chan JobResult, nCores)
	for _, job := range jobs {
		go job.Run(c)
	}
	for i := 0; i < nCores; i++ {
		result := <-c
		for j, p := range result.Pixels {
			index := result.Start + j
			x := index % s.Camera.Width
			y := index / s.Camera.Width
			img.Set(x, y, p)
		}
	}
	timeEnd := time.Now()
	fmt.Printf("Rendered in %v\n", timeEnd.Sub(timeStart))
	return img
}
