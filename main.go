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

func radiance(ray Ray, world *HittableList, maxDepth int32, envMap *SphereMap) Color32 {
	li := Vector3{}
	throughput := Vector3{1.0, 1.0, 1.0}
	hitRecord := HitRecord{}
	for depth := int32(0); depth < maxDepth; depth++ {
		if !world.Hit(ray, 0.001, Infinity32, &hitRecord) {
			unitDirection := NormalizeVector3(ray.Direction)
			//t := 0.5 * (unitDirection.Y + 1.0)
			//v := AddVector3(MulVector3(1.0-t, Vector3{1.0, 1.0, 1.0}), MulVector3(t, Vector3{0.5, 0.7, 1.0}))
			v := envMap.Sample(unitDirection)
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
	//world.AddHittable(&Sphere{Vector3{0.0, -1000.0, 0.0}, 1000.0, &Lambertian{Vector3{0.5, 0.5, 0.5}}})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Vector3{float32(a) + 0.9*rand.Float32(), 0.2, float32(b) + 0.9*rand.Float32()}
			v := SubVector3(center, Vector3{4.0, 0.2, 0.0})
			l := v.Length()
			if l <= 0.9 {
				continue
			}
			color := Vector3{rand.Float32(), rand.Float32(), rand.Float32()}
			roughness := rand.Float32()*0.9 + 0.01
			metallic := rand.Float32()*0.9 + 0.01
			world.AddHittable(&Sphere{center, 0.2, &Metal{color, roughness, metallic, 0.9}})
/*
			selection := rand.Float32()
			if selection < 0.4 {
				world.AddHittable(&Sphere{center, 0.2, &Lambertian{color}})
			} else if selection < 0.8 {
				roughness := rand.Float32()*0.5 + 0.1
				world.AddHittable(&Sphere{center, 0.2, &Metal{color, roughness, 0.9}})
			} else {
				world.AddHittable(&Sphere{center, 0.2, &Dielectric{color, rand.Float32()}})
			}
*/
		}
	}
	//world.AddHittable(&Sphere{Vector3{0.0, 1.0, 0.0}, 1.0, &Dielectric{Vector3{1.0, 1.0, 1.0}, 1.5}})
	//world.AddHittable(&Sphere{Vector3{-4.0, 1.0, 0.0}, 1.0, &Lambertian{Vector3{0.4, 0.2, 0.1}}})
	world.AddHittable(&Sphere{Vector3{4.0, 1.0, 0.0}, 1.0, &Metal{Vector3{0.7, 0.6, 0.5}, 0.05, 0.5, 0.9}})
	return world
}

func render(name string, width, height, spp, maxDepth int32, world *HittableList) {
	fmt.Printf("start render %v ...\n", name)
	start := time.Now()

	var envMap SphereMap
	envMap.Load("uffizi_probe.hdr")

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
				c := radiance(ray, world, maxDepth, &envMap)
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

func SampleSpecularEnvMap(roughness float32, direction Vector3, specularMaps []SphereMap) Vector3 {
	maxLevels := int32(len(specularMaps))
	r := float32(maxLevels-1) * roughness
	l0 := int32(r)
	d0 := r - float32(l0)
	var l1 int32
	if maxLevels<=l0 {
		l0 = maxLevels - 1
		d0 = 0.0
		l1 = l0
	} else {
		l1 = l0 + 1
		if maxLevels<=l1 {
			l1 = l0
		}
	}
	c0 := specularMaps[l0].Sample(direction)
	c1 := specularMaps[l1].Sample(direction)
	return LerpVector3(c0, c1, d0)
}

func SampleIrradiance(N Vector3, irradianceMap *SphereMap) Vector3 {
	return irradianceMap.Sample(N)
}

func SampleSpecularAsIrradiance(N Vector3, specularMaps []SphereMap) Vector3 {
	maxLevels := int32(len(specularMaps))
	return specularMaps[maxLevels-3].Sample(N)
}

func SampleBRDF(roughness, NV float32, brdfMap *SphereMap) Vector3 {
	return brdfMap.Pixel(roughness, NV)
}

func FresnelSchlickRoughness(cosTheta, roughness float32, F0 Vector3) Vector3 {
	r := 1.0-roughness
	F1 := Vector3{math32.Max(r, F0.X), math32.Max(r, F0.Y), math32.Max(r, F0.Z)}
	return AddVector3(F0, MulVector3(math32.Pow(Clamp0132(1-cosTheta), 5.0), SubVector3(F1,F0)))
}

func FresnelF0(albedo Vector3, metallic float32) Vector3 {
	return Vector3{Lerp32(0.04, albedo.X, metallic), Lerp32(0.04, albedo.Y, metallic), Lerp32(0.04, albedo.Z, metallic)}
}

func DistributionGGX(NH, roughness float32) float32 {
	a := roughness * roughness
	a2 := a*a
	NH2 := NH*NH
	denom := NH2 * (a2-1.0) + 1.0
	denom = math32.Pi * denom * denom
	return a2/denom
}

func GeometrySchlickGGX(NV, roughness float32) float32 {
	r := roughness + 1.0
	k := (r*r)/8.0
	denom := NV * (1.0-k) + k
	return NV/denom
}

func GeometrySmith(NV, NL, roughness float32) float32 {
	ggx0 := GeometrySchlickGGX(NV, roughness)
	ggx1 := GeometrySchlickGGX(NL, roughness)
	return ggx0 * ggx1
}

func FresnelSchlick(cosTheta float32, F0 Vector3) Vector3 {
	F1 := Vector3{1.0-F0.X, 1.0-F0.Y, 1.0-F0.Z}
	return AddVector3(F0, MulVector3(math32.Pow(Clamp0132(1-cosTheta), 5.0), F1))
}

func radiance_direct(ray Ray, world *HittableList, envMap, irradianceMap, brdfMap *SphereMap, specularMaps []SphereMap, useAsIrradiance bool) Color32 {
	li := Vector3{}
	hitRecord := HitRecord{}
	if !world.Hit(ray, 0.001, Infinity32, &hitRecord) {
		unitDirection := NormalizeVector3(ray.Direction)
		li = envMap.Sample(unitDirection)
		return Color32{li.X, li.Y, li.Z, 1.0}
	}
	L := Vector3{0.0, 1.0, 0.0}
	V := NormalizeVector3(SubVector3(ray.Origin, hitRecord.Position))
	N := hitRecord.Normal
	H := NormalizeVector3(AddVector3(V, L))
	R := Reflect(V.Minus(), N)

	NL := math32.Max(DotVector3(N, L), 0.0)
	NV := math32.Max(DotVector3(N, V), 0.0)
	HV := math32.Max(DotVector3(H, V), 0.0)
	NH := math32.Max(DotVector3(N, H), 0.0)

	roughness := hitRecord.Material.GetRoughness()
	metallic := hitRecord.Material.GetMetallic()
	albedo := hitRecord.Material.GetAlbedo()

	F0 := FresnelF0(albedo, metallic)

	NDF := DistributionGGX(NH, roughness)
	G := GeometrySmith(NV, NL, roughness)
	F := FresnelSchlick(HV, F0)
	kS := F
	kD := MulVector3(1.0-metallic, Vector3{1.0-kS.X, 1.0-kS.Y, 1.0-kS.Z})
	denom := 4.0 * NV * NL + 0.0001
	specular := MulVector3(1.0/denom, MulVector3(NDF * G, F))
	Lo := AddVector3(MulVector3(1.0/math32.Pi, HadamardDotVector3(kD, albedo)), specular)
	Lo = MulVector3(NL*0.1, Lo)

	envF := FresnelSchlickRoughness(NV, roughness, F0)
	aS := envF
	aD := MulVector3(1.0-metallic, Vector3{1.0-aS.X, 1.0-aS.Y, 1.0-aS.Z})
	var irradiance Vector3
	if useAsIrradiance {
		irradiance = SampleSpecularAsIrradiance(N, specularMaps)
	}else {
		irradiance = SampleIrradiance(N, irradianceMap)
	}
	ambientD := HadamardDotVector3(irradiance, albedo)

	reflection := SampleSpecularEnvMap(roughness, R, specularMaps)
	BRDF := SampleBRDF(roughness, NV, brdfMap)

	ambientS := HadamardDotVector3(AddVector3(MulVector3(BRDF.X, envF), Vector3{BRDF.Y, BRDF.Y, BRDF.Y}), reflection)
	ambient := AddVector3(HadamardDotVector3(aD, ambientD), ambientS)
	Lo = AddVector3(Lo, MulVector3(0.9, ambient))
	//Lo = AddVector3(ambientS, MulVector3(0.0, Lo))

	return Color32{Lo.X, Lo.Y, Lo.Z, 1.0}
}

func render_direct(name string, width, height int32, world *HittableList, useAsIrradiance bool) {
	fmt.Printf("start render %v ...\n", name)
	start := time.Now()

	var envMap SphereMap
	envMap.Load("uffizi_probe.hdr")
	irradianceMap := envMap.GenIrradiance(128, 128)
	//irradianceMap.SavePng("irradiance.png")
	specularMaps := envMap.GenSpecular(128, 128, 6)

	//for i:=0; i<6; i++ {
	//	specularMaps[i].SavePng(fmt.Sprintf("specular%d.png", i))
	//}

	brdfMap := envMap.GenBRDF(256, 256)
	//brdfMap.SavePng("brdf.png")

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(width), int(height)}})

	camera := NewCameraPerspectiveLens(uint32(width), uint32(height), DegToRad32*45.0, 0.01)
	camera.LookAt(Vector3{9.0, 1.2, 2.5}, Vector3{0.0, 0.0, 0.0}, Vector3{0.0, 1.0, 0.0})
	sample := Sample2{0.0, 0.0}
	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			ray := camera.GenerateRay(uint32(x), uint32(y), sample, sample)
			c := radiance_direct(ray, world, &envMap, &irradianceMap, &brdfMap, specularMaps, useAsIrradiance)
			c = LinearToSRGB(c)
			img.Set(int(x), int(height-y-1), color32ToRGBA(c))
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
	world := generateScene()
	var width int32 = 400
	var height int32 = 300
	var numSamples int32 = 512
	var maxDepth int32 = 16
	render("out_path.png", width, height, numSamples, maxDepth, &world)
	render_direct("out_ibl.png", width, height, &world, false)
	render_direct("out_ibl_pseudo.png", width, height, &world, true)
}

