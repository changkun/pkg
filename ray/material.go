package main

import (
	"math"
	"math/rand"
)

type material interface {
	scatter(rin *ray, rec *hitRecord, attenuation *vec, scattered *ray) bool
}

type lambertian struct {
	albedo *vec
}

func newLambertian(a *vec) *lambertian {
	return &lambertian{a}
}

func (l *lambertian) scatter(rin *ray, rec *hitRecord, attenuation *vec, scattered *ray) bool {
	target := rec.p.add(&rec.normal).add0(randInUnitSphere(1)).sub0(&rec.p)
	scattered.set(&rec.p, target)
	attenuation.setV(l.albedo)
	return true
}

type metal struct {
	albedo *vec
	fuzzy  float64
}

func newMetal(a *vec, fuzzy float64) *metal {
	if fuzzy < 1 {
		return &metal{a, fuzzy}
	}
	return &metal{a, 1}
}

func (m *metal) scatter(rin *ray, rec *hitRecord, attenuation *vec, scattered *ray) bool {
	reflected := rin.dir.uni().reflect(&rec.normal).add0(randInUnitSphere(m.fuzzy))
	scattered.set(&rec.p, reflected)
	attenuation.setV(m.albedo)
	return scattered.dir.dot(&rec.normal) > 0
}

type dielectric struct {
	refIdx float64
}

func newDielectric(ri float64) *dielectric {
	return &dielectric{ri}
}

func (d *dielectric) scatter(rin *ray, rec *hitRecord, attenuation *vec, scattered *ray) bool {
	var (
		outwardNormal *vec
		niOverNt      float64
		reflectProb   float64
		cosine        float64
		refracted     = newVec(0, 0, 0)
	)
	attenuation.set(1, 1, 1)

	if rin.dir.dot(&rec.normal) > 0 {
		outwardNormal = rec.normal.mult(-1)
		niOverNt = d.refIdx
		cosine = d.refIdx * rin.dir.dot(&rec.normal) / rin.dir.len()
	} else {
		outwardNormal = &rec.normal
		niOverNt = 1.0 / d.refIdx
		cosine = -rin.dir.dot(&rec.normal) / rin.dir.len()

	}
	if rin.dir.refract(outwardNormal, niOverNt, refracted) {
		reflectProb = schlick(cosine, d.refIdx)
	} else {
		reflectProb = 1.0
	}
	if rand.Float64() < reflectProb {
		scattered.set(&rec.p, rin.dir.reflect(&rec.normal))
	} else {
		scattered.set(&rec.p, refracted)
	}
	return true
}

func schlick(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 *= r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

func randInUnitSphere(fuzz float64) (x, y, z float64) {
	for {
		x, y, z = rand.Float64()*2-1, rand.Float64()*2-1, rand.Float64()*2-1
		if x*x+y*y+z*z >= 1.0 {
			return x * fuzz, y * fuzz, z * fuzz
		}
	}
}
