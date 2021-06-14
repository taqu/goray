package core

type Ray struct {
	Origin    Vector3
	Direction Vector3
}

func (ray *Ray) PointAt(t float32) Vector3 {
	return AddVector3(ray.Origin, MulVector3(t, ray.Direction))
}

