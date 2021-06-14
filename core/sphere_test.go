package core
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSphereHit(t *testing.T) {
	assert := assert.New(t)
	sphere := Sphere{Vector3{0.0, 0.0, 0.0}, 0.5, nil}
	var hitRecord HitRecord
	ray0 := Ray{Vector3{0.0, 0.51, 1.0}, NormalizeVector3(Vector3{0.0, 0.0, -1.0})}
	assert.Falsef(sphere.Hit(ray0, 0.0, 1.0, &hitRecord), "%v doesn't hit %v", ray0, sphere)
	ray1 := Ray{Vector3{0.0, 0.499, 1.0}, NormalizeVector3(Vector3{0.0, 0.0, -1.0})}
	assert.Truef(sphere.Hit(ray1, 0.0, 1.0, &hitRecord), "%v hit %v at %v", ray1, sphere, ray1.PointAt(hitRecord.T))
}

