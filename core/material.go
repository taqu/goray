package core

import (
	"math/rand"
	"git.maze.io/go/math32"
)

type MaterialSample struct {
	Continue bool
	PDF float32
	Weight Vector3
	Scattered Vector3
}

type Material interface {
	Sample(wi Vector3, eta0, eta1 float32) MaterialSample
	Scatter(ray *Ray, hitRecord *HitRecord, attenuation *Vector3, scattered *Ray) bool
	GetRoughness() float32
	GetMetallic() float32
	GetAlbedo() Vector3
}

type Lambertian struct {
	Albedo Vector3
}

func NewLambertian(albedo Vector3) Lambertian {
	return Lambertian{albedo}
}

func (lambertian *Lambertian) Sample(wi Vector3, eta0, eta1 float32) MaterialSample {
	if wi.Z <= Epsilon32 {
		return MaterialSample{false, 0.0, Vector3{}, Vector3{}}
	}
	wm := RandomOnCosineHemiSphere(eta0, eta1)
	pdf := wi.Z / math32.Pi //PDF of cosine hemisphere
	return MaterialSample{true, pdf, lambertian.Albedo, wm}
}

func (material *Lambertian) Scatter(ray *Ray, hitRecord *HitRecord, attenuation *Vector3, scattered *Ray) bool {
	coordinate := NewCoordinate(hitRecord.Normal)
	n := RandomOnHemiSphere(rand.Float32(), rand.Float32())
	*scattered = Ray{hitRecord.Position, coordinate.LocalToWorld(n)}
	*attenuation = material.Albedo
	return true
}

func (material *Lambertian) GetRoughness() float32 {
	return 1.0
}

func (material *Lambertian) GetMetallic() float32 {
	return 0.0
}

func (material *Lambertian) GetAlbedo() Vector3 {
	return material.Albedo
}

func project2(x, y float32, v Vector3) float32 {
	if Equal32(x, y) {
		return x*y
	}
	sn := 1.0 - v.Z*v.Z
	if sn <= Epsilon32 {
		return x*y
	}
	invSn := 1.0/sn
	cs2 := v.X*v.X*invSn
	sn2 := v.Y*v.Y*invSn
	return cs2*x*x + sn2*y*y
}

func ggx_NDF(m Vector3, alpha2 float32) float32 {
	d := m.Z
	denom := (d*d) * (alpha2 - 1.0) + 1.0
	return alpha2/(math32.Pi * denom * denom)
}

func ggx_G1(v Vector3, alpha2 float32) float32 {
	dotNV := v.Z
	denom := math32.Sqrt(alpha2 + (1.0-alpha2)*dotNV*dotNV) + dotNV
	return 2.0 * dotNV/denom
}

func ggx_G2(wi, wo Vector3, alpha2 float32) float32 {
	dotNI := wi.Z
	dotNO := wo.Z
	denomI := dotNO * math32.Sqrt(alpha2 + (1.0 - alpha2) * dotNI*dotNI)
	denomO := dotNI * math32.Sqrt(alpha2 + (1.0 - alpha2) * dotNO*dotNO)
	return 2.0 * dotNI * dotNO / (denomI + denomO)
}

func ggx_VNDF(wo Vector3, roughness, eta0, eta1 float32) Vector3 {
	//1. Transform the view direction to the hemisphere configuration
	v := NormalizeVector3(Vector3{roughness*wo.X, roughness*wo.Y, wo.Z})

    //2. Construct orthonormal bais
	var t1 Vector3
	if v.Z < 0.999 {
		t1 = NormalizeVector3(CrossVector3(v, Vector3{0.0, 0.0, 1.0}))
	}else {
		t1 = Vector3{1.0, 0.0, 0.0}
	}
	t2 := CrossVector3(t1, v)

    //3. Make parameterization of the projected area
	/*
	a := 1.0/(1.0+v.Y)
	r := math32.Sqrt(eta0)
	var phi float32
	var b float32
	if eta1<a {
		phi = eta1/a * math32.Pi
		b = 1.0
	}else{
		phi = math32.Pi + (eta1-a)/(1.0-a) * math32.Pi;
		b = v.Z
	}
	p1 := r * math32.Cos(phi)
	p2 := r * math32.Sin(phi) * b
	*/
	r := math32.Sqrt(eta0)
	phi := (math32.Pi * 2.0) * eta1
	p1 := r * math32.Cos(phi)
	p2 := r * math32.Sin(phi)
	a := 0.5 * (1.0 + v.Z)
	p2 = (1.0-a)*math32.Sqrt(1.0-p1*p1) + a*p2

    //4. Reproject onto hemisphere
	t1 = MulVector3(p1, t1)
	t2 = MulVector3(p2, t2)
	v = MulVector3(math32.Sqrt(math32.Max(0.0, 1.0-p1*p1-p2*p2)), v)
	n := AddVector3(AddVector3(t1, t2), v)

    //5. Transform the normal back to the elipsoid configuration
	return NormalizeVector3(Vector3{roughness*n.X, roughness*n.Y, math32.Max(0.0, n.Z)})
}

// https://hal.archives-ouvertes.fr/hal-01509746/document
type Metal struct {
	Albedo Vector3
	Roughness float32
	Metallic float32
	RefIndex float32
}

func (metal *Metal) Sample(wi Vector3, eta0, eta1 float32) MaterialSample {
	wm := ggx_VNDF(wi, metal.Roughness, eta0, eta1)
	wo := SubVector3(MulVector3(2.0*DotVector3(wi, wm), wm), wi)
	if 0.0 < wo.Z {
		f := Schlick(DotVector3(wo, wm), metal.RefIndex)
		g1 := ggx_G1(wi, metal.Roughness)
		g2 := ggx_G2(wo, wi, metal.Roughness*metal.Roughness)
		pdf := f * (g2/g1)
		return MaterialSample{true, pdf, metal.Albedo, wm}
	}else{
		return MaterialSample{false, 0.0, metal.Albedo, wm}
	}

	/*
	alpha2 := metal.Roughness * metal.Roughness
	theta := math32.Acos(math32.Sqrt((1.0-eta0)/((alpha2-1.0)*eta0 + 1.0)))
	phi := 2.0 * math32.Pi * eta1
	sinTheta := math32.Sin(theta)
	cosTheta := math32.Cos(theta)
	sinPhi := math32.Sin(phi)
	cosPhi := math32.Cos(phi)
	wm := Vector3{sinTheta*cosPhi, sinTheta*sinPhi, cosTheta}
	wo := SubVector3(MulVector3(2.0*DotVector3(wi, wm), wm), wi)

	if 0.0 < wo.Z  && 0.0<DotVector3(wo, wm) {
		f := Schlick(DotVector3(wo, wm), metal.RefIndex)
		g := ggx_G2(wo, wi, alpha2)
		pdf := f * g * math32.Abs(DotVector3(wi,wm))/(wi.Z * wm.Z)
		return MaterialSample{true, pdf, metal.Albedo, wm}
	}else{
		return MaterialSample{false, 0.0, metal.Albedo, wm}
	}
	*/
}

func (metal *Metal) Scatter(ray *Ray, hitRecord *HitRecord, attenuation *Vector3, scattered *Ray) bool {
	reflected := NormalizeVector3(Reflect(ray.Direction, hitRecord.Normal))
	*scattered = Ray{hitRecord.Position, reflected}
	*attenuation = metal.Albedo
	return 0.0001 < DotVector3(scattered.Direction, hitRecord.Normal)
}

func (material *Metal) GetRoughness() float32 {
	return material.Roughness
}

func (material *Metal) GetMetallic() float32 {
	return material.Metallic
}

func (material *Metal) GetAlbedo() Vector3 {
	return material.Albedo
}

type Dielectric struct {
	Albedo   Vector3
	RefIndex float32
}

func (dielectric *Dielectric) Sample(wi Vector3, eta0, eta1 float32) MaterialSample {
	if wi.Z <= Epsilon32 {
		return MaterialSample{false, 0.0, Vector3{}, Vector3{}}
	}

	normal := Vector3{0.0, 0.0, 1.0}
	direction := wi.Minus()
	reflected := NormalizeVector3(Reflect(direction, normal))
	var niOverNt float32
	var cosine float32
	var n Vector3
	if 0.0 < direction.Z {
		n = Vector3{0.0, 0.0, -1.0}
		niOverNt = dielectric.RefIndex
		cosine = dielectric.RefIndex * direction.Z
	}else{
		n = Vector3{0.0, 0.0, 1.0}
		niOverNt = 1.0 / dielectric.RefIndex
		cosine = -direction.Z
	}

	var refracted Vector3
	if !Refract(&refracted, direction, n, niOverNt) {
		return MaterialSample{true, 1.0, Vector3{1.0, 1.0, 1.0}, reflected}
	}
	reflectProb := Schlick(cosine, dielectric.RefIndex)
	if rand.Float32() < reflectProb {
		return MaterialSample{true, 1.0, Vector3{1.0, 1.0, 1.0}, reflected}
	}else{
		return MaterialSample{true, 1.0, Vector3{1.0, 1.0, 1.0}, refracted}
	}
}

func (dielectric *Dielectric) Scatter(ray *Ray, hitRecord *HitRecord, attenuation *Vector3, scattered *Ray) bool {
	reflected := NormalizeVector3(Reflect(ray.Direction, hitRecord.Normal))
	*attenuation = dielectric.Albedo
	var niOverNt float32
	var cosine float32
	var normal Vector3
	if 0.0 < DotVector3(ray.Direction, hitRecord.Normal) {
		normal = hitRecord.Normal.Minus()
		niOverNt = dielectric.RefIndex
		cosine = dielectric.RefIndex * DotVector3(ray.Direction, hitRecord.Normal)

	} else {
		normal = hitRecord.Normal
		niOverNt = 1.0 / dielectric.RefIndex
		cosine = -DotVector3(ray.Direction, hitRecord.Normal)
	}

	var refracted Vector3
	if !Refract(&refracted, ray.Direction, normal, niOverNt) {
		*scattered = Ray{hitRecord.Position, reflected}
		return true
	}
	reflectProb := Schlick(cosine, dielectric.RefIndex)
	if rand.Float32() < reflectProb {
		*scattered = Ray{hitRecord.Position, reflected}
	}else{
		*scattered = Ray{hitRecord.Position, refracted}
	}
	return true
}

func (material *Dielectric) GetRoughness() float32 {
	return 0.0
}

func (material *Dielectric) GetMetallic() float32 {
	return 0.0
}

func (material *Dielectric) GetAlbedo() Vector3 {
	return material.Albedo
}

