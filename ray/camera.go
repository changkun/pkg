package main

import (
	"math"
	"math/rand"
)

type camera struct {
	origin, lowerLeftCorner *vec
	horizontal, vertical    *vec
	u, v, w                 *vec
	lensRadius              float64
}

func newCamera(lookfrom, lookat, vup *vec, vfov, aspect, aperture, focusDist float64) *camera {

	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	cam := &camera{
		origin:     lookfrom,
		lensRadius: aperture / 2,
	}
	cam.w = lookfrom.sub(lookat).uni()
	cam.u = vup.cross(cam.w).uni()
	cam.v = cam.w.cross(cam.u)
	cam.lowerLeftCorner = cam.origin.sub(
		cam.u.mult(focusDist * halfWidth)).sub(
		cam.v.mult(focusDist * halfHeight)).sub(
		cam.w.mult(focusDist))
	cam.horizontal = cam.u.mult(2 * focusDist * halfWidth)
	cam.vertical = cam.v.mult(2 * focusDist * halfHeight)
	return cam
}

func (c *camera) getRay(s, t float64, r *ray) {
	rdx, rdy := randInUnitDisk(c.lensRadius)
	offset := c.u.mult(rdx).add(c.v.mult(rdy))

	r.orig.addV(c.origin).addV(offset)
	r.dir.addV(c.horizontal).mult0(s).addV(c.lowerLeftCorner).addV(
		c.vertical.mult(t)).sub0(c.origin).sub0(offset)
}

func randInUnitDisk(r float64) (x, y float64) {
	for {
		x, y = 2*rand.Float64()-1, 2*rand.Float64()-1
		if x*x+y*y < 1.0 {
			return x * r, y * r
		}
	}
}
