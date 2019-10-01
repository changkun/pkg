package main

import "math"

const maxDepth = 50

type ray struct {
	orig, dir vec
}

func newRay() *ray {
	return &ray{}
}

func (r *ray) set(orig, dir *vec) {
	r.orig.setV(orig)
	r.dir.setV(dir)
}

func (r *ray) reset() {
	r.orig.reset()
	r.dir.reset()
}

func (r *ray) pointAt(t float64) (x, y, z float64) {
	x = r.orig.x + r.dir.x*t
	y = r.orig.y + r.dir.y*t
	z = r.orig.z + r.dir.z*t
	return
}

func (r *ray) hintSphere(center *vec, radius float64) float64 {
	oc := r.orig.sub(center)
	a := r.dir.dot(&r.dir)
	b := 2.0 * oc.dot(&r.dir)
	c := oc.dot(oc) - radius*radius
	d := b*b - 4*a*c
	if d < 0 {
		return -1
	}
	return (-b - math.Sqrt(d)) / (2.0 * a)
}

func (r *ray) color(world *hitableList, depth int) *vec {
	t := 0.5 * (r.dir.y*(1.0/r.dir.len()) + 1.0)
	attenuation := newVec(1-0.5*t, 1-0.3*t, 1)

	rec := &hitRecord{}
	if world.hit(r, 0.001, math.MaxFloat64, rec) {
		scattered := newRay()
		if depth < maxDepth && rec.mat.scatter(r, rec, attenuation, scattered) {
			tmp := scattered.color(world, depth+1)
			attenuation.multV0(tmp)
		}
	}

	return attenuation
}
