package main

import (
	"math"
)

type sphere struct {
	center *vec
	radius float64
	mat    material
}

func newSphere(c *vec, r float64, mat material) *sphere {
	return &sphere{c, r, mat}
}

func (s *sphere) hit(r *ray, tmin, tmax float64, rec *hitRecord) bool {
	ocx := r.orig.x - s.center.x
	ocy := r.orig.y - s.center.y
	ocz := r.orig.z - s.center.z
	a := r.dir.dot(&r.dir)
	b := 2.0 * (ocx*r.dir.x + ocy*r.dir.y + ocz*r.dir.z)
	c := (ocx*ocx + ocy*ocy + ocz*ocz) - s.radius*s.radius
	d := b*b - 4*a*c
	if d > 0 {
		tmp := (-b - math.Sqrt(d)) / (2.0 * a)
		if tmp < tmax && tmp > tmin {
			rec.t = tmp
			rec.p.set(r.pointAt(rec.t))
			rec.normal.set(
				(rec.p.x-s.center.x)*(1.0/s.radius),
				(rec.p.y-s.center.y)*(1.0/s.radius),
				(rec.p.z-s.center.z)*(1.0/s.radius),
			)
			rec.mat = s.mat
			return true
		}
		tmp = (-b + math.Sqrt(d)) / (2.0 * a)
		if tmp < tmax && tmp > tmin {
			rec.t = tmp
			rec.p.set(r.pointAt(rec.t))
			rec.normal.set(
				(rec.p.x-s.center.x)*(1.0/s.radius),
				(rec.p.y-s.center.y)*(1.0/s.radius),
				(rec.p.z-s.center.z)*(1.0/s.radius),
			)
			rec.mat = s.mat
			return true
		}
	}
	return false
}
