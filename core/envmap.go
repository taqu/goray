package core

import (
	"os"
	"git.maze.io/go/math32"
	"github.com/Opioid/rgbe"
	"image"
	"image/color"
	"image/png"
	"time"
)

type SphereMap struct {
	Width int32
	Height int32
	Image []Vector3
}

func (env *SphereMap) Load(path string) {
	fi, _ := os.Open(path)
	defer fi.Close()
	var width int
	var height int
	var tmp []float32
	width, height, tmp, _ = rgbe.Decode(fi)
	env.Width = int32(width)
	env.Height = int32(height)
	image := make([]Vector3, width*height)

	for i:=0; i<height; i++ {
		for j:=0; j<width; j++ {
			dst := i*width + j
			src := dst*3
			image[dst].X = tmp[src + 0]
			image[dst].Y = tmp[src + 1]
			image[dst].Z = tmp[src + 2]
		}
	}
	env.Image = image
}

func (env *SphereMap) Save(path string) {
	image := make([]float32, env.Width*env.Height*3)
	for i := int32(0); i<env.Height; i++ {
		for j := int32(0); j<env.Width; j++ {
			src := i*env.Width + j
			dst := src*3
			image[dst + 0] = env.Image[src].X
			image[dst + 1] = env.Image[src].Y
			image[dst + 2] = env.Image[src].Z
		}
	}
	fo, _ := os.Create(path)
	defer fo.Close()
	rgbe.Encode(fo, int(env.Width), int(env.Height), image)
}

func (env *SphereMap) SavePng(path string) {
	clamp := func (x float32) uint8 {
		x = x*255.99
		ix := int(x)
		if ix<0 {
			return 0
		}else if 255<ix {
			return 255
		}else {
			return uint8(ix)
		}
	}
	img := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{int(env.Width), int(env.Height)}})
	for i := int32(0); i<env.Height; i++ {
		for j := int32(0); j<env.Width; j++ {
			src := i*env.Width + j
			//luminance := 0.299 * env.Image[src].X + 0.587 * env.Image[src].Y + 0.114 * env.Image[src].Z
			tone := float32(1.0) //luminance/(luminance+1.0)
			r := clamp(env.Image[src].X * tone)
			g := clamp(env.Image[src].Y * tone)
			b := clamp(env.Image[src].Z * tone)
			img.Set(int(j), int(i), color.RGBA{r, g, b, 255})
		}
	}
	fo, _ := os.Create(path)
	defer fo.Close()
	png.Encode(fo, img)
}

func (env *SphereMap) Sample(n Vector3) Vector3 {
	r := (1/math32.Pi) * math32.Acos(n.Z)/math32.Sqrt(n.X*n.X + n.Y*n.Y)
	u := (n.X*r + 1)*0.5
	v := (1-n.Y*r)*0.5
	x := u*float32(env.Width)
	y := v*float32(env.Height)
	return env.Pixel(x, y)
}

func (env *SphereMap) Normal(x, y float32) Vector3 {
	u := x*2.0 - 1.0
	v := 1.0 - y*2.0
	theta := math32.Atan2(u,v)
	phi := math32.Pi * math32.Sqrt(u*u + v*v)
	rz := -math32.Cos(phi)
	rx := math32.Sin(phi) * math32.Sin(-theta)
	ry := math32.Sin(phi) * math32.Cos(-theta)
	return NormalizeVector3(Vector3{rx,ry,rz})
}

func (env *SphereMap) Pixel(x, y float32) Vector3 {
	clamp := func(x, minx, maxx int32) int32 {
		if x<minx {
			return minx
		} else if maxx<x {
			return maxx
		}
		return x
	}
	clamp01 := func(x float32) float32 {
		if(x<0){
			return 0
		}else if(1<x){
			return 1
		}
		return x
	}
	lerp := func(x0, x1 Vector3, t float32) Vector3 {
		it := 1.0-t
		x := x0.X * it + x1.X * t
		y := x0.Y * it + x1.Y * t
		z := x0.Z * it + x1.Z * t
		return Vector3{x,y,z}
	}

	ix := clamp(int32(x), 0, env.Width-1)
	iy := clamp(int32(y), 0, env.Height-1)
	dx := clamp01(x-float32(ix))
	dy := clamp01(y-float32(iy))
	ix2 := clamp(ix+1, 0, env.Width-1)
	iy2 := clamp(iy+1, 0, env.Height-1)
	c00 := env.Image[iy*env.Width + ix]
	c01 := env.Image[iy*env.Width + ix2]
	c10 := env.Image[(iy2)*env.Width + ix]
	c11 := env.Image[(iy2)*env.Width + ix2]
	c0 := lerp(c00, c01, dx)
	c1 := lerp(c10, c11, dx)
	return lerp(c0, c1, dy)
}

func (env *SphereMap) GenIrradiance(width, height int32) SphereMap {
	const samples int = 4096
	image := make([]Vector3, width*height)
	invw := 1.0/float32(width-1)
	invh := 1.0/float32(height-1)
	r2 := NewSamplerR2(time.Now().UnixNano())
	for i:=int32(0); i<height; i++ {
		y := invh * float32(i)
		for j:=int32(0); j<width; j++ {
			x := invw * float32(j)
			n := env.Normal(x, y)
			total := Vector3{0.0, 0.0, 0.0}
			for s:=0; s<samples; s++ {
				sample := r2.Generate2(int32(s))
				direction := RandomOnHemiSphereAround(sample.X, sample.Y, n)
				weight := DotVector3(direction, n)
				value := MulVector3(weight, env.Sample(direction))
				total = AddVector3(total, value)
			}
			image[i*width + j] = MulVector3(math32.Pi/float32(samples), total)
		}
	}
	return SphereMap{width, height, image}
}

func ImportanceSampleGGX(x, y, roughness float32, n Vector3) Vector3 {
	a := roughness * roughness
	phi := 2.0 * math32.Pi * x
	cosTheta := math32.Sqrt((1.0-y)/(1.0 + (a*a-1.0)*y))
	sinTheta := math32.Sqrt(1.0 - cosTheta*cosTheta)
	h := Vector3{math32.Cos(phi)*sinTheta, math32.Sin(phi)*sinTheta, cosTheta}
	var up Vector3
	if math32.Abs(n.Z)<0.999 {
		up = Vector3{0.0, 0.0, 1.0}
	}else {
		up = Vector3{1.0, 0.0, 0.0}
	}
	binormal0 := NormalizeVector3(CrossVector3(up, n))
	binormal1 := CrossVector3(n, binormal0)
	v0 := MulVector3(h.X, binormal0)
	v1 := MulVector3(h.Y, binormal1)
	v2 := MulVector3(h.Z, n)
	return NormalizeVector3(AddVector3(AddVector3(v0,v1),v2))
}

func (env *SphereMap) GetSpecularOne(width, height int32, roughness float32) SphereMap {
	const samples int = 4096
	r2 := NewSamplerR2(time.Now().UnixNano())
	image := make([]Vector3, width*height)
	invw := 1.0/float32(width-1)
	invh := 1.0/float32(height-1)
	for i:=int32(0); i<height; i++ {
		y := invh * float32(i)
		for j:=int32(0); j<width; j++ {
			x := invw * float32(j)
			n := env.Normal(x, y)
			total := Vector3{0.0, 0.0, 0.0}
			totalWeight := float32(0)
			for s:=0; s<samples; s++ {
				sample := r2.Generate2(int32(s))
				direction := ImportanceSampleGGX(sample.X, sample.Y, roughness, n)
				weight := DotVector3(direction, n)
				value := MulVector3(weight, env.Sample(direction))
				total = AddVector3(total, value)
				totalWeight = totalWeight + weight
			}
			image[i*width + j] = MulVector3(1.0/totalWeight, total)
		}
	}
	return SphereMap{width, height, image}
}

func (env *SphereMap) GenSpecular(width, height, miplevels int32) []SphereMap {
	maps := make([]SphereMap, miplevels)
	for i:=int32(0); i<miplevels; i++ {
		roughness := float32(i)/float32(miplevels-1)
		maps[i] = env.GetSpecularOne(width, height, roughness)
		width = width/2
		height = height/2
	}
	return maps
}

func GeometrySchlickGGX(NV, roughness float32) float32 {
	k := (roughness * roughness) * 0.5
	denom := NV * (1.0-k) + k
	return NV/denom
}

func GeometrySmith(N, V, L Vector3, roughness float32) float32 {
	NV := math32.Max(DotVector3(N,V), 0.0)
	NL := math32.Max(DotVector3(N,L), 0.0)
	return GeometrySchlickGGX(NL, roughness) * GeometrySchlickGGX(NV, roughness)
}

func (env *SphereMap) GenBRDF(width, height int32) SphereMap {
	const samples int = 1024
	r2 := NewSamplerR2(time.Now().UnixNano())
	image := make([]Vector3, width*height)
	invw := 1.0/float32(width-1)
	invh := 1.0/float32(height-1)
	N := Vector3{0.0, 0.0, 1.0}
	for i:=int32(0); i<height; i++ {
		y := invh * float32(i)
		for j:=int32(0); j<width; j++ {
			roughness := invw * float32(j)
			NV := y
			V := Vector3{math32.Sqrt(1.0-NV*NV), 0.0, NV}
			total := Vector3{0.0, 0.0, 0.0}
			for s:=0; s<samples; s++ {
				sample := r2.Generate2(int32(s))
				H := ImportanceSampleGGX(sample.X, sample.Y, roughness, N)
				L := NormalizeVector3(SubVector3(MulVector3(2.0*DotVector3(V,H), H), V))
				NL := math32.Max(L.Z, 0.0)
				NH := math32.Max(H.Z, 0.0)
				VH := math32.Max(DotVector3(V,H), 0.0)
				if 0.0<NL {
					G := GeometrySmith(N, V, L, roughness)
					Gv := (G*VH)/(NH*NV)
					Fc := math32.Pow(1.0-VH, 5.0)
					total.X += (1.0-Fc)*Gv
					total.Y += Fc*Gv
				}
			}
			image[i*width + j] = MulVector3(1.0/float32(samples), total)
		}
	}
	return SphereMap{width, height, image}
}

