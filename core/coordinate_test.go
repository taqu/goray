package core
import (
	"testing"
	"math/rand"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestWorldToLocal(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<100; i++ {
		n := NormalizeVector3(Vector3{rand.Float32(), rand.Float32(), rand.Float32()})
		coordinate := NewCoordinate(n)
		x := coordinate.WorldToLocal(n)
		d := DotVector3(Vector3{0.0, 0.0, 1.0}, x)
		l := x.Length()
		assert.Truef(Equal32(1.0, l), "Length should be %v", l)
		assert.Truef(Equal32(1.0, d), "Dot of local %v and %v should be 0.0 (%v)", n, x, DotVector3(n,x))
	}
}

func TestLocalToWorld(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<100; i++ {
		n := NormalizeVector3(Vector3{rand.Float32(), rand.Float32(), rand.Float32()})
		coordinate := NewCoordinate(n)
		x := coordinate.LocalToWorld(Vector3{0.0, 0.0, 1.0})
		d := DotVector3(n, x)
		l := x.Length()
		assert.Truef(Equal32(1.0, l), "Length should be %v", l)
		assert.Truef(Equal32(1.0, d), "Dot of local %v and %v should be 0.0 (%v)", n, x, DotVector3(n,x))
	}
}


func TestRandomHemiSphereCoordinate(t *testing.T) {
	assert := assert.New(t)
	normal := NormalizeVector3(Vector3{0.5, -0.5, 0.5})
	coordinate := NewCoordinate(normal)
	for i:=0; i<100; i++ {
		n0 := RandomOnHemiSphere(rand.Float32(), rand.Float32())
		n1 := coordinate.LocalToWorld(n0)
		l := n1.Length()
		d := DotVector3(n0, n1)
		assert.Truef(Equal32(1.0, l), "Length shoudle be %v", l)
		assert.Truef(0.0<=d && d<=1.0, "Dot of %v and %v sholud be (0 1)", n0, n1)
	}
}
