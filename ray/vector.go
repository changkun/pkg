package main

import "math"

type vec struct {
	x, y, z float64 // r, g, b
}

func newVec(x, y, z float64) *vec {
	return &vec{x, y, z}
}

func (v *vec) r() float64 {
	return v.x
}
func (v *vec) g() float64 {
	return v.y
}
func (v *vec) b() float64 {
	return v.z
}
func (v *vec) setV(vv *vec) *vec {
	v.x = vv.x
	v.y = vv.y
	v.z = vv.z
	return v
}
func (v *vec) set(x, y, z float64) *vec {
	v.x = x
	v.y = y
	v.z = z
	return v
}
func (v *vec) reset() {
	v.x = 0
	v.y = 0
	v.z = 0
}
func (v *vec) add0(x, y, z float64) *vec {
	v.x += x
	v.y += y
	v.z += z
	return v
}
func (v *vec) addV(vv *vec) *vec {
	v.x += vv.x
	v.y += vv.y
	v.z += vv.z
	return v
}
func (v *vec) sqrt0() {
	v.x = math.Sqrt(v.x)
	v.y = math.Sqrt(v.y)
	v.z = math.Sqrt(v.z)
}
func (v *vec) add(vv *vec) *vec {
	return &vec{v.x + vv.x, v.y + vv.y, v.z + vv.z}
}
func (v *vec) sub0(vv *vec) *vec {
	v.x -= vv.x
	v.y -= vv.y
	v.z -= vv.z
	return v
}
func (v *vec) sub(vv *vec) *vec {
	return &vec{v.x - vv.x, v.y - vv.y, v.z - vv.z}
}
func (v *vec) neg() *vec {
	return &vec{-v.x, -v.y, -v.z}
}
func (v *vec) mult0(vv float64) *vec {
	v.x *= vv
	v.y *= vv
	v.z *= vv
	return v
}
func (v *vec) mult(vv float64) *vec {
	return &vec{v.x * vv, v.y * vv, v.z * vv}
}
func (v *vec) multV0(vv *vec) {
	v.x *= vv.x
	v.y *= vv.y
	v.z *= vv.z
}
func (v vec) dot(vv *vec) float64 {
	return v.x*vv.x + v.y*vv.y + v.z*vv.z
}
func (v *vec) len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}
func (v *vec) cross(vv *vec) *vec {
	return &vec{
		v.y*vv.z - v.z*vv.y,
		-(v.x*vv.z - v.z*vv.x),
		v.x*vv.y - v.y*vv.x,
	}
}
func (v *vec) uni() *vec {
	return v.mult(1.0 / v.len())
}

func (v *vec) reflect(n *vec) *vec {
	return v.sub(n.mult(v.dot(n) * 2))
}

func (v *vec) refract(n *vec, niOverNt float64, refracted *vec) bool {
	uv := v.uni()
	dt := uv.dot(n)
	discriminant := 1.0 - niOverNt*(1-dt*dt)
	if discriminant <= 0 {
		return false
	}
	*refracted = *uv.sub(n.mult(dt)).mult(niOverNt).sub(n.mult(math.Sqrt(discriminant)))
	return true
}
