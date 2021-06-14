package core

import (
	"git.maze.io/go/math32"
)

type Camera struct {
	Width   uint32
	Height  uint32
	Aspect  float32
	DX    float32
	DY    float32
	Origin  Vector3
	Forward Vector3
	Right   Vector3
	Up      Vector3
	LensRadius float32
}

func (camera *Camera) LookAt(eye, at, up Vector3) {
	forward := NormalizeVector3(SubVector3(at, eye))
	cs := DotVector3(forward, up)
	if 0.999 < math32.Abs(cs) {
		up = Vector3{forward.Z, forward.X, forward.Y}
	}
	right := NormalizeVector3(CrossVector3(forward, up))
	up = NormalizeVector3(CrossVector3(right, forward))

	camera.Origin = eye
	camera.Forward = forward
	camera.Right = right
	camera.Up = up
}

//Screen coordinate to Normalized Device Coordinate (NDC)
func screenToNDC(x, resolution uint32, jitter float32) float32 {
	return 2.0*((float32(x)+0.5+jitter)/float32(resolution)) - 1.0
}

func (camera *Camera) GenerateRay(x, y uint32, screenSample, lensSample Sample2) Ray {
	lensSample = RandomOnDisk(lensSample.X, lensSample.Y).Mul(camera.LensRadius)

	originUp := MulVector3(lensSample.X, camera.Up)
	originRight := MulVector3(lensSample.Y, camera.Right)
	origin := AddVector3(camera.Origin, AddVector3(originUp, originRight));

	dx := camera.DX * screenToNDC(x, camera.Width, screenSample.X-0.499)
	dy := camera.DY * screenToNDC(y, camera.Height, screenSample.Y-0.499)
	right := MulVector3(dx, camera.Right)
	up := MulVector3(dy, camera.Up)
	direction := NormalizeVector3(AddVector3(AddVector3(right, up), camera.Forward))

	return Ray{origin, direction}
}

func NewCameraPerspectiveFov(width uint32, height uint32, fovy float32) Camera {
	aspect := float32(width) / float32(height)
	fovy = math32.Tan(0.5 * fovy)
	fovx := fovy * aspect
	return Camera{width, height, aspect, fovx, fovy,
		Vector3{}, Vector3{0.0, 0.0, -1.0}, Vector3{1.0, 0.0, 0.0}, Vector3{0.0, 1.0, 0.0},
		0.0}
}

func NewCameraPerspectiveLens(width, height uint32, fovy, aperture float32) Camera {
	aspect := float32(width) / float32(height)
	fovy = math32.Tan(0.5 * fovy)
	fovx := fovy * aspect
	return Camera{width, height, aspect, fovx, fovy,
		Vector3{}, Vector3{0.0, 0.0, -1.0}, Vector3{1.0, 0.0, 0.0}, Vector3{0.0, 1.0, 0.0},
		aperture*0.5}
}

