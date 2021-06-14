package main

import (
	"fmt"
	"git.maze.io/go/math32"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	. "ray/core"
	"time"
)

func toRGBA(x Vector3) color.RGBA {
	r := uint8(255.99 * x.X)
	g := uint8(255.99 * x.Y)
	b := uint8(255.99 * x.Z)
	return color.RGBA{r, g, b, 0xFF}
}

func color32ToRGBA(c Color32) color.RGBA {
	c.R = Saturate32(c.R)
	c.G = Saturate32(c.G)
	c.B = Saturate32(c.B)
	c.A = Saturate32(c.A)
	r := uint8(255.99 * c.R)
	g := uint8(255.99 * c.G)
	b := uint8(255.99 * c.B)
	a := uint8(255.99 * c.A)
	return color.RGBA{r, g, b, a}
}

func radiance(ray Ray, world HittableList, maxDepth int32) Color32 {
	li := Vector3{}
	throughput := Vector3{1.0, 1.0, 1.0}
	hitRecord := HitRecord{}
	for depth := int32(0); depth < maxDepth; depth++ {
		if !world.Hit(ray, 0.001, Infinity32, &hitRecord) {
			unitDirection := NormalizeVector3(ray.Direction)
			t := 0.5 * (unitDirection.Y + 1.0)
			v := AddVector3(MulVector3(1.0-t, Vector3{1.0, 1.0, 1.0}), MulVector3(t, Vector3{0.5, 0.7, 1.0}))
			li = AddVector3(HadamardDotVector3(throughput, v), li)
			break
		}
		coordinate := NewCoordinate(hitRecord.Normal)
		wow := ray.Direction.Minus()
		wo := coordinate.WorldToLocal(wow)
		materialSample := hitRecord.Material.Sample(wo, rand.Float32(), rand.Float32())
		if materialSample.Weight.IsZero() {
			ray.Origin = AddVector3(MulVector3(Epsilon32, ray.Direction), hitRecord.Position)
		} else {
			wiw := coordinate.LocalToWorld(materialSample.Scattered)
			throughput = HadamardDotVector3(throughput, materialSample.Weight)

			ray.Origin = hitRecord.Position
			ray.Direction = wiw
		}
		//Russian roulette
		if 6 <= depth {
			continueProbability := math32.Min(throughput.Length(), 0.9)
			if continueProbability <= rand.Float32() {
				break
			}
			throughput = DivVector3(throughput, continueProbability)
		}
	}
	return Color32{li.X, li.Y, li.Z, 1.0}
}

func generateScene() HittableList {
	world := NewHittableList()
	world.AddHittable(&Sphere{Vector3{0.0, -1000.0, 0.0}, 1000.0, &Lambertian{Vector3{0.5, 0.5, 0.5}}})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Vector3{float32(a) + 0.9*rand.Float32(), 0.2, float32(b) + 0.9*rand.Float32()}
			v := SubVector3(center, Vector3{4.0, 0.2, 0.0})
			l := v.Length()
			if l <= 0.9 {
				continue
			}
			selection := rand.Float32()
			color := Vector3{rand.Float32(), rand.Float32(), rand.Float32()}
			if selection < 0.4 {
				world.AddHittable(&Sphere{center, 0.2, &Lambertian{color}})
			} else if selection < 0.8 {
				roughness := rand.Float32()*0.5 + 0.1
				world.AddHittable(&Sphere{center, 0.2, &Metal{color, roughness, 0.9}})
			} else {
				world.AddHittable(&Sphere{center, 0.2, &Dielectric{color, rand.Float32()}})
			}
		}
	}
	world.AddHittable(&Sphere{Vector3{0.0, 1.0, 0.0}, 1.0, &Dielectric{Vector3{1.0, 1.0, 1.0}, 1.5}})
	world.AddHittable(&Sphere{Vector3{-4.0, 1.0, 0.0}, 1.0, &Lambertian{Vector3{0.4, 0.2, 0.1}}})
	world.AddHittable(&Sphere{Vector3{4.0, 1.0, 0.0}, 1.0, &Metal{Vector3{0.7, 0.6, 0.5}, 0.05, 0.9}})
	return world
}

func render(name string, width, height, spp, maxDepth int32) {
	fmt.Printf("start render %v ...\n", name)
	start := time.Now()

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	//random := rand.New(rand.NewSource(1))
	//rand.Seed(1)
	screenSamples := GoldenSet(int(spp), random)
	//lensSamples := GoldenSet(int(spp), random)
	lensSamples := SamplerSet(int(spp), NewSamplerJitteredR2(0.05, time.Now().UnixNano()))
	//screenSamples := RandomSet(int(spp), random)
	//lensSamples := RandomSet(int(spp), random)
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(width), int(height)}})
	world := generateScene()
	sigma := float32(0.5)
	gauss0 := float32(1.0/math32.Sqrt(2.0*math32.Pi*sigma*sigma))
	gauss1 := float32(-1.0/(2.0*sigma*sigma))

	camera := NewCameraPerspectiveLens(uint32(width), uint32(height), DegToRad32*45.0, 0.01)
	camera.LookAt(Vector3{9.0, 1.2, 2.5}, Vector3{0.0, 0.0, 0.0}, Vector3{0.0, 1.0, 0.0})
	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			acc := Color32{}
			weight := float32(0.0)
			for s := int32(0); s < spp; s++ {
				ray := camera.GenerateRay(uint32(x), uint32(y), screenSamples[s], lensSamples[s])
				c := radiance(ray, world, maxDepth)
				dx := 2.0 * screenSamples[s].X - 1.0
				dy := 2.0 * screenSamples[s].Y - 1.0
				w := gauss0 * math32.Exp(gauss1*(dx*dx + dy*dy))
				weight += w
				c = MulColor32(w, c)
				acc = AddColor32(acc, c)
			}

			if Epsilon32<weight {
				acc = MulColor32(1.0/weight, acc)
			}
			acc = LinearToSRGB(acc)
			img.Set(int(x), int(height-y-1), color32ToRGBA(acc))
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("done (%v ms)\n", int64(elapsed/time.Millisecond))

	file, err := os.Create(name)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		return
	}
}

func main() {
	var width int32 = 400
	var height int32 = 300
	var numSamples int32 = 200
	var maxDepth int32 = 20
	render("outimage.png", width, height, numSamples, maxDepth)
}
