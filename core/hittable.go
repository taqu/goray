package core

type Hittable interface {
	Hit(ray Ray, tmin float32, tmax float32, record *HitRecord) bool
}

