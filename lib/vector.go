package lib

import (
	"fmt"
	"image/color"
	"math"
)

type Vector struct {
	X, Y, Z float64
}

type unitVector struct {
	Vector
}

func (v Vector) Unit() unitVector {
	return unitVector{v.MulScalar(1 / math.Sqrt(v.Dot(v)))}
}

func I() unitVector {
	return unitVector{Vector{1, 0, 0}}
}

func J() unitVector {
	return unitVector{Vector{0, 1, 0}}
}

func K() unitVector {
	return unitVector{Vector{0, 0, 1}}
}

func Zero() Vector {
	return Vector{0, 0, 0}
}

func White() Vector {
	return Vector{1, 1, 1}
}

func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v Vector) AddScalar(s float64) Vector {
	return Vector{v.X + s, v.Y + s, v.Z + s}
}

func (v Vector) Sub(w Vector) Vector {
	return Vector{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v Vector) SubScalar(s float64) Vector {
	return Vector{v.X - s, v.Y - s, v.Z - s}
}

func (v Vector) Mul(w Vector) Vector {
	return Vector{v.X * w.X, v.Y * w.Y, v.Z * w.Z}
}

func (v Vector) MulScalar(s float64) Vector {
	return Vector{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector) Dot(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vector) Cross(w Vector) Vector {
	return Vector{
		v.Y*w.Z - v.Z*w.Y,
		v.Z*w.X - v.X*w.Z,
		v.X*w.Y - v.Y*w.X,
	}
}

// project v onto w
func (v Vector) Project(w Vector) Vector {
	return w.MulScalar(v.Dot(w) / w.Dot(w))
}

// reflect v over w
func (v Vector) Reflect(w Vector) Vector {
	return v.Sub(w.MulScalar(2 * v.Dot(w) / w.Dot(w)))
}

// rotate v around w by theta radians
func (v Vector) Rotate(w Vector, theta float64) Vector {
	parallel := v.Project(w)
	ortho := v.Sub(parallel)
	orthoLen := math.Sqrt(ortho.Dot(ortho))
	outOfPlane := w.Cross(ortho)
	inplaneDist := math.Cos(theta) / orthoLen
	outofPlaneDist := math.Sin(theta) / math.Sqrt(outOfPlane.Dot(outOfPlane))
	orthoRot := ortho.MulScalar(inplaneDist).Add(outOfPlane.MulScalar(outofPlaneDist)).MulScalar(orthoLen)
	return parallel.Add(orthoRot)
}

func (v Vector) Trim(min, max float64) Vector {
	return Vector{
		math.Min(math.Max(v.X, min), max),
		math.Min(math.Max(v.Y, min), max),
		math.Min(math.Max(v.Z, min), max),
	}
}

func (v Vector) ToColor() (*color.RGBA, error) {
	if v.X < 0 || v.Y < 0 || v.Z < 0 || v.X > 1 || v.Y > 1 || v.Z > 1 {
		return nil, fmt.Errorf("Vector %v cannot be converted to color; all coordinates must be in [0, 1]", v)
	}
	return &color.RGBA{
		uint8(v.X * 255),
		uint8(v.Y * 255),
		uint8(v.Z * 255),
		255,
	}, nil
}

func (v Vector) String() string {
	return fmt.Sprintf("(%v, %v, %v)", v.X, v.Y, v.Z)
}
