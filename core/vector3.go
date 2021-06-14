package core

import (
	"git.maze.io/go/math32"
)

type Vector3 struct {
	X, Y, Z float32
}

func (x *Vector3) LengthSqr() float32 {
	return DotVector3(*x, *x)
}

func (x *Vector3) Length() float32 {
	return math32.Sqrt(x.LengthSqr())
}

func (x *Vector3) Minus() Vector3 {
	return Vector3{-x.X, -x.Y, -x.Z}
}

func (x *Vector3) IsZero() bool {
	return EqualZero32(x.X) && EqualZero32(x.Y) && EqualZero32(x.Z)
}

func EqualVector3(x0, x1 Vector3) bool {
	return Equal32(x0.X, x1.X) && Equal32(x0.Y, x1.Y) && Equal32(x0.Z, x1.Z)
}

func AddVector3(x0, x1 Vector3) Vector3 {
	return Vector3{x0.X + x1.X, x0.Y + x1.Y, x0.Z + x1.Z}
}

func SubVector3(x0, x1 Vector3) Vector3 {
	return Vector3{x0.X - x1.X, x0.Y - x1.Y, x0.Z - x1.Z}
}

func MulVector3(x0 float32, x1 Vector3) Vector3 {
	return Vector3{x0 * x1.X, x0 * x1.Y, x0 * x1.Z}
}

func DivVector3(x0 Vector3, x1 float32) Vector3 {
	inv := 1.0 / x1
	return Vector3{x0.X * inv, x0.Y * inv, x0.Z * inv}
}

func DotVector3(x0, x1 Vector3) float32 {
	return x0.X*x1.X + x0.Y*x1.Y + x0.Z*x1.Z
}

func CrossVector3(x0, x1 Vector3) Vector3 {
	x := x0.Y*x1.Z - x0.Z*x1.Y
	y := x0.Z*x1.X - x0.X*x1.Z
	z := x0.X*x1.Y - x0.Y*x1.X
	return Vector3{x, y, z}
}

func NormalizeVector3(x Vector3) Vector3 {
	invL := 1.0 / math32.Sqrt(DotVector3(x, x))
	return MulVector3(float32(invL), x)
}

func HadamardDotVector3(x0 Vector3, x1 Vector3) Vector3 {
	return Vector3{x0.X * x1.X, x0.Y * x1.Y, x0.Z * x1.Z}
}

func SaturateVector3(x Vector3) Vector3 {
	x.X = Saturate32(x.X)
	x.Y = Saturate32(x.Y)
	x.Z = Saturate32(x.Z)
	return x
}

func RandomInSphere(x0, x1, x2 float32) Vector3 {
	theta := 2.0*x0 - 1.0
	r := math32.Sqrt(1.0 - theta*theta)
	phi := (math32.Pi * 2.0) * x1
	sn := math32.Sin(phi)
	cs := math32.Cos(phi)
	r *= x2
	return Vector3{r * cs, r * sn, x2 * theta}
}

func RandomOnSphere(x0, x1 float32) Vector3 {
	theta := 2.0*x0 - 1.0
	r := math32.Sqrt(1.0 - theta*theta)
	phi := (math32.Pi * 2.0) * x1
	sn := math32.Sin(phi)
	cs := math32.Cos(phi)
	return Vector3{r * cs, r * sn, theta}
}

func RandomOnHemiSphere(x0, x1 float32) Vector3 {
	theta := x0
	r := math32.Sqrt(1.0 - theta*theta)
	phi := (math32.Pi * 2.0) * x1
	sn := math32.Sin(phi)
	cs := math32.Cos(phi)
	return Vector3{r * cs, r * sn, theta}
}

func RandomOnCosineHemiSphere(x0, x1 float32) Vector3 {
	p := RandomOnDisk(x0, x1)
	z := math32.Max(Epsilon32, math32.Sqrt(math32.Max(Epsilon32, (1.0-p.X*p.X-p.Y*p.Y))))
	return Vector3{p.X, p.Y, z}
}

func RandomCone(x0, x1, cosCutoff float32) Vector3 {
	cosTheta := (1.0 - x0) + x0*cosCutoff
	sinTheta := math32.Sqrt(math32.Max(Epsilon32, (1.0 - cosTheta*cosTheta)))
	phi := 2.0 * math32.Pi * x1
	sinPhi := math32.Sin(phi)
	cosPhi := math32.Cos(phi)
	return Vector3{cosPhi * sinTheta, sinPhi * sinTheta, cosTheta}
}

func Reflect(x, n Vector3) Vector3 {
	return SubVector3(x, MulVector3(2.0*DotVector3(x, n), n))
}

func Refract(refracted *Vector3, x, n Vector3, niOverNt float32) bool {
	dt := DotVector3(x, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1.0-dt*dt)
	if 0.0 < discriminant {
		*refracted = SubVector3(MulVector3(niOverNt, SubVector3(x, MulVector3(dt, n))), MulVector3(math32.Sqrt(discriminant), n))
		return true
	}
	return false
}
