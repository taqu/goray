package core

import (
	"git.maze.io/go/math32"
	"math/rand"
)

const (
	R2_G0     float32 = float32(1.0 / 1.61803398874989484820458683436563)
	R2_G1     float32 = float32(1.0 / 1.32471795724474602596090885447809)
	R2_G2     float32 = float32(1.0 / 1.22074408460575947536168534910883)
	R2_Delta  float32 = float32(0.76)
	R2_I0     float32 = float32(0.3)
	R2_SqrtPi float32 = float32(1.77245385091)
	Magic     float32 = float32(0.618033988749894)
)

type Sample2 struct {
	X float32
	Y float32
}

func (sample2 Sample2) Mul(x float32) Sample2 {
	return Sample2{sample2.X * x, sample2.Y * x}
}

type Sampler interface {
	Generate(n int32) float32
	Generate2(n int32) Sample2
}

type SamplerRandom struct {
	random *rand.Rand
}

func NewSamplerRandom(seed int64) *SamplerRandom {
	return &SamplerRandom{rand.New(rand.NewSource(seed))}
}

// Generate 1d sample
// [0.0 1.0)
func (sampler *SamplerRandom) Generate(_ int32) float32 {
	return sampler.random.Float32()
}

// Generate2 2d sample
// [0.0 1.0)
func (sampler *SamplerRandom) Generate2(_ int32) Sample2 {
	x := sampler.random.Float32()
	y := sampler.random.Float32()
	return Sample2{x, y}
}

// SamplerR2
//
// http://extremelearning.com.au/unreasonable-effectiveness-of-quasirandom-sequences/
type SamplerR2 struct {
}

func NewSamplerR2(seed int64) *SamplerR2 {
	return &SamplerR2{}
}

func (sampler *SamplerR2) Generate(n int32) float32 {
	a1 := R2_G0
	x := 0.5 + a1*float32(n)
	return x - math32.Floor(x)
}

// Generate2 2d sample
// [-0.5 0.5)
func (sampler *SamplerR2) Generate2(n int32) Sample2 {
	fn := float32(n)
	a1 := R2_G1
	a2 := R2_G1 * R2_G1
	x := 0.5 + a1*fn
	y := 0.5 + a2*fn
	return Sample2{x - math32.Floor(x), y - math32.Floor(y)}
}

// SamplerJitteredR2
type SamplerJitteredR2 struct {
	lambda float32
	random *rand.Rand
}

func NewSamplerJitteredR2(delta float32, seed int64) *SamplerJitteredR2 {
	random := rand.New(rand.NewSource(seed))
	return &SamplerJitteredR2{delta, random}
}

func (sampler *SamplerJitteredR2) Generate(n int32) float32 {
	fn := float32(n)
	x := 0.5 + R2_G0*fn
	p := x - math32.Floor(x)
	u := sampler.random.Float32()
	k := sampler.lambda * R2_Delta * R2_SqrtPi / (4.0 * math32.Sqrt(fn+R2_I0))
	p += k * u
	p -= math32.Floor(p)
	return p
}

// Generate2 2d sample
// [-0.5 0.5)
func (sampler *SamplerJitteredR2) Generate2(n int32) Sample2 {
	fn := float32(n)
	a1 := R2_G1
	a2 := R2_G1 * R2_G1
	x := 0.5 + a1*fn
	y := 0.5 + a2*fn
	p := Sample2{x - math32.Floor(x), y - math32.Floor(y)}
	u := Sample2{sampler.random.Float32(), sampler.random.Float32()}
	k := sampler.lambda * R2_Delta * R2_SqrtPi / (4.0 * math32.Sqrt(fn+R2_I0))
	p.X += k * u.X
	p.Y += k * u.Y
	p.X -= math32.Floor(p.X)
	p.Y -= math32.Floor(p.Y)
	return p
}

func RandomSet(numSamples int, random *rand.Rand) []Sample2 {
	points := make([]Sample2, numSamples)
	for i := 0; i< numSamples; i++ {
		points[i].X = random.Float32()
		points[i].Y = random.Float32()
	}
	return points
}

func SamplerSet(numSamples int, sampler Sampler) []Sample2 {
	points := make([]Sample2, numSamples)
	for i := int32(0); i<int32(numSamples); i++ {
		s := sampler.Generate2(i)
		points[i].X = s.X
		points[i].Y = s.Y
	}
	return points
}

// Golden
//
// Schretter Colas, Kobbelt Leif, Dehaye Paul-Olivier, "Golden Ratio Sequences for Low-Discrepancy Sampling", JCGT 2012
// https://www.graphics.rwth-aachen.de/publication/2/jgt.pdf
func GoldenSet(numSamples int, random *rand.Rand) []Sample2 {
	points := make([]Sample2, numSamples)
	x := random.Float32()
	min := x
	index := 0
	for i := 0; i < numSamples; i++ {
		points[i].Y = x
		if x < min {
			min = x
			index = i
		}
		x += Magic
		if 1.0 <= x {
			x -= 1.0
		}
	}
	f := 1
	fp := 1
	parity := 0
	for (f + fp) < numSamples {
		tmp := f
		f += fp
		fp = tmp
		parity++
	}
	var inc, dec int
	if 0 != (parity & 0x01) {
		inc = f
		dec = fp
	} else {
		inc = fp
		dec = f
	}

	points[0].X = points[index].Y
	for i := 1; i < numSamples; i++ {
		if index < dec {
			index += inc
			if numSamples <= index {
				index -= dec
			}
		} else {
			index -= dec
		}
		points[i].X = points[index].Y
	}
	y := random.Float32()
	for i := 0; i < numSamples; i++ {
		points[i].Y = y
		y += Magic
		if 1.0 <= y {
			y -= 1.0
		}
	}
	return points
}

/*
func RandomOnDisk(x0, x1 float32) Sample2 {
	r := math32.Sqrt(1.0 - x0*x0)
	phi := 2.0 * math32.Pi * x1
	sn := math32.Sin(phi)
	cs := math32.Cos(phi)
	return Sample2{r * cs, r * sn}
}
*/

func RandomOnDisk(x0, x1 float32) Sample2 {
	// http://psgraphics.blogspot.ch/2011/01/improved-code-for-concentric-map.html
	r0 := 2.0*x0 - 1.0
	r1 := 2.0*x1 - 1.0

	absR0 := math32.Abs(r0)
	absR1 := math32.Abs(r1)
	var phi float32
	var r float32
	if absR0 <= Epsilon32 && absR1 <= Epsilon32 {
		phi = 0.0
		r = 0.0
	} else if absR1 < absR0 {
		phi = (math32.Pi / 4.0) * (r1 / r0)
		r = r0
	} else {
		r = r1
		phi = (math32.Pi / 2.0) - (r0/r1)*(math32.Pi/4.0)
	}

	return Sample2{r * math32.Cos(phi), r * math32.Sin(phi)}
}
