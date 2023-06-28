package lib

import (
	"math"
)

const PLANE_TOL float64 = 0.001

type Surface struct {
	// all in range [0, 1]
	Ambient, Diffuse, Specular float64
	Color                      Vector
}

type Shape interface {
	// returns location of intersection between object and ray
	// if no intersection, returns nil
	Intersection(Ray) *Vector
	// returns normal vector at given point
	Normal(Vector) unitVector
}

type Object struct {
	Shape
	Surface
}

type Sphere struct {
	Center Vector
	Radius float64
}

func (s Sphere) Intersection(r Ray) *Vector {
	// check that this sphere is not the ray source
	if s.Center == r.Origin {
		return nil
	}
	rayOriginToCenter := s.Center.Sub(r.Origin)
	scalarProd := rayOriginToCenter.Dot(r.Direction.Vector)
	if scalarProd < 0 {
		// ray is pointing away from sphere
		return nil
	}
	// rsq is squared distance between sphere center and projection of rayOriginToCenter onto ray
	rsq := rayOriginToCenter.Dot(rayOriginToCenter) - scalarProd*scalarProd
	sRadSq := s.Radius * s.Radius
	if rsq > sRadSq {
		// ray misses sphere
		return nil
	}
	// ray hits sphere
	// find distance from ray origin to intersection
	lengthInSphere := math.Sqrt(sRadSq - rsq)
	originToIntersection := r.Direction.MulScalar((scalarProd - lengthInSphere))
	intersection := originToIntersection.Add(r.Origin)
	return &intersection
}

func (s Sphere) Normal(p Vector) unitVector {
	return p.Sub(s.Center).Unit()
}

// a single-sided plane
type Plane struct {
	Point Vector     // a point on the plane
	Norm  unitVector // normal vector facing away from viewable side of plane
}

func (p Plane) Intersection(r Ray) *Vector {
	dist := p.Point.Sub(r.Origin).Dot(p.Norm.Vector) / r.Direction.Dot(p.Norm.Vector)
	if dist < PLANE_TOL {
		// ray is not pointing toward the plane
		return nil
	}
	intersection := r.Direction.MulScalar(dist).Add(r.Origin)
	// fmt.Println(dist, intersection)
	return &intersection
}

func (p Plane) Normal(v Vector) unitVector {
	return p.Norm
}
