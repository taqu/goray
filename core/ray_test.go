package core
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRayPointAt(t *testing.T){
	assert := assert.New(t)
	ray := Ray{Vector3{0.0, 0.0, 0.0}, NormalizeVector3(Vector3{1.0, 1.0, 1.0})}
	time := float32(2.0)
	point := ray.PointAt(time)
	expected := Vector3{1.1547005, 1.1547005, 1.1547005}
	assert.Truef(EqualVector3(point, expected), "point of %v at %v is %v", ray, time, point)
}

