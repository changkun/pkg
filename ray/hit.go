package main

type hitRecord struct {
	t      float64
	p      vec
	normal vec
	mat    material
}

type hitable interface {
	hit(r *ray, tmin, tmax float64, rec *hitRecord) bool
}

type hitableList struct {
	objs []hitable
}

func newHitableList() *hitableList {
	return &hitableList{}
}

func (l *hitableList) add(obj hitable) {
	l.objs = append(l.objs, obj)
}

func (l hitableList) hit(r *ray, tmin, tmax float64, rec *hitRecord) bool {
	hitAnything := false
	closestSoFar := tmax
	for i := range l.objs {
		if l.objs[i].hit(r, tmin, closestSoFar, rec) {
			hitAnything = true
			closestSoFar = rec.t
		}
	}
	return hitAnything
}
