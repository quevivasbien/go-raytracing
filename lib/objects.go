package lib

import "math"

type Surface struct {
	// all in range [0, 1]
	Ambient, Diffuse, Specular float64
	Color                      Vector
}

type Shape interface {
	// returns location of light source
	RaySource() Vector
	// returns location of intersection between object and ray
	// if no intersection, returns nil
	Intersection(Ray) *Vector
}

type Object struct {
	Shape
	Surface
}

type Sphere struct {
	Center Vector
	Radius float64
}

func (s Sphere) RaySource() Vector {
	return s.Center
}

func (s Sphere) Intersection(r Ray) *Vector {
	// check that this sphere is not the ray source
	if s.Center == r.Origin {
		return nil
	}
	rayOriginToCenter := s.Center.Sub(r.Origin)
	scalarProd := rayOriginToCenter.Dot(r.Direction.Vector)
	// check that ray is not pointing away from sphere
	if scalarProd < 0 {
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
