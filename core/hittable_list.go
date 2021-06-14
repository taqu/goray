package core

type HittableList struct {
	hittables []Hittable
}

func (hittableList *HittableList) AddHittable(hittable Hittable) {
	hittableList.hittables = append(hittableList.hittables, hittable)
}

func (hittableList *HittableList) AddHittables(hittables []Hittable) {
	hittableList.hittables = append(hittables, hittableList.hittables...)
}

func (hittableList *HittableList) Hit(ray Ray, tmin float32, tmax float32, record *HitRecord) bool {
	tmp := HitRecord{}
	hitAnything := false
	closestSoFar := tmax
	for i:=0; i<len(hittableList.hittables); i++ {
		if !hittableList.hittables[i].Hit(ray, tmin, closestSoFar, &tmp) {
			continue
		}
		hitAnything = true
		closestSoFar = tmp.T
		*record = tmp
	}
	return hitAnything
}

func NewHittableList() HittableList {
	return HittableList{}
}

