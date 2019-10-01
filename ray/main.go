package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"
)

var (
	nx, ny, ns int
)

func init() {
	flag.IntVar(&nx, "w", 1400, "iamge width")
	flag.IntVar(&ny, "h", 700, "image height")
	flag.IntVar(&ns, "s", 50, "ray sample")
	flag.Parse()
}

func setCamera(nx, ny int) *camera {
	lookfrom := newVec(13, 2, 3)
	lookat := newVec(0, 0, 0)
	distToFocus := 10.0
	aperture := 0.1

	cam := newCamera(lookfrom, lookat, newVec(0, 1, 0), 20,
		float64(nx)/float64(ny), aperture, distToFocus)
	return cam
}

func setWorld() *hitableList {
	world := newHitableList()
	world.add(newSphere(newVec(0, -1000, 0), 1000,
		newLambertian(newVec(0.5, 0.5, 0.5))))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := newVec(float64(a)+0.9*rand.Float64(), 0.2,
				float64(b)+0.9*rand.Float64())
			if center.sub(newVec(4, 0.2, 0)).len() > 0.9 { // diffuse
				world.add(newSphere(center, 0.2,
					newLambertian(newVec(rand.Float64()*rand.Float64(),
						rand.Float64()*rand.Float64(),
						rand.Float64()*rand.Float64()))))
			} else if chooseMat < 0.95 { // metal
				world.add(newSphere(center, 0.2, newMetal(
					newVec(0.5*(1+rand.Float64()),
						0.5*(1+rand.Float64()),
						0.5*(1+rand.Float64())),
					0.5*rand.Float64())))
			} else { // glass
				world.add(newSphere(center, 0.2, newDielectric(1.5)))
			}
		}
	}
	world.add(newSphere(newVec(0, 1, 0), 1.0, newDielectric(1.5)))
	world.add(newSphere(newVec(-4, 1, 0), 1.0,
		newLambertian(newVec(0.4, 0.2, 0.1))))
	world.add(newSphere(newVec(4, 1, 0), 1.0,
		newMetal(newVec(0.7, 0.6, 0.5), 0.0)))
	return world
}

func render() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, nx, ny))

	world := setWorld()
	cam := setCamera(nx, ny)

	wg := sync.WaitGroup{}
	wg.Add(ny)
	for j := 0; j < ny; j++ {
		go func(j int) {
			r := newRay()
			for i := 0; i < nx; i++ {
				c := newVec(0, 0, 0)
				for s := 0; s < ns; s++ {
					r.reset()
					u := (float64(i) + rand.Float64()) / float64(nx)
					v := (float64(j) + rand.Float64()) / float64(ny)
					cam.getRay(u, v, r)
					c.addV(r.color(world, 0))
				}
				c.mult0(1.0 / float64(ns))
				c.sqrt0()
				img.Set(i, ny-j+1, color.NRGBA{
					R: uint8(255.99 * c.r()),
					G: uint8(255.99 * c.g()),
					B: uint8(255.99 * c.b()),
					A: 255,
				})
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	return img
}

func main() {
	// defer profile.Start().Stop()

	img := render()
	f, err := os.Create("x.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
