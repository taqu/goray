package core

type Coordinate struct {
	Normal    Vector3
	Binormal0 Vector3
	Binormal1 Vector3
}

// OrthonormalBasis
//
// Jeppe Revall Frisvad, "Building an Orthonormal Basis from a 3D Unit Vector Without Normalization", 2012 Jornal of Graphics Tools
func NewCoordinate(normal Vector3) Coordinate {
	if normal.Z < -0.9999999 {
		return Coordinate{normal, Vector3{0.0, -1.0, 0.0}, Vector3{-1.0, 0.0, 0.0}}
	}
	a := 1.0 / (1.0 + normal.Z)
	b := -normal.X * normal.Y * a
	binormal0 := Vector3{1.0 - normal.X*normal.X*a, b, -normal.X}
	binormal1 := Vector3{b, 1.0 - normal.Y*normal.Y*a, -normal.Y}
	return Coordinate{normal, binormal0, binormal1}
}

func (coordinate *Coordinate) WorldToLocal(v Vector3) Vector3 {
	x := DotVector3(v, coordinate.Binormal0)
	y := DotVector3(v, coordinate.Binormal1)
	z := DotVector3(v, coordinate.Normal)
	return NormalizeVector3(Vector3{x, y, z})
}

func (coordinate *Coordinate) LocalToWorld(v Vector3) Vector3 {
	x := coordinate.Binormal0.X*v.X + coordinate.Binormal1.X*v.Y + coordinate.Normal.X*v.Z
	y := coordinate.Binormal0.Y*v.X + coordinate.Binormal1.Y*v.Y + coordinate.Normal.Y*v.Z
	z := coordinate.Binormal0.Z*v.X + coordinate.Binormal1.Z*v.Y + coordinate.Normal.Z*v.Z
	return Vector3{x, y, z}
}
