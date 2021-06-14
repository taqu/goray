package core
import (
	"testing"
	"math/rand"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestLengthSqrVector3(t *testing.T) {
	assert := assert.New(t)

	v := Vector3{2.0, 2.0, 2.0}
	l := v.LengthSqr()
	assert.Truef(Equal32(12.0, l), "Length of %v should be 12.0", v)
}

func TestLengthVector3(t *testing.T) {
	assert := assert.New(t)
	v := Vector3{1.0, 1.0, 1.0}
	l := v.Length()
	assert.Truef(Equal32(1.7320508, l), "Length of %v should be 1.7320508", v)
}

func TestEqualVector3(t *testing.T) {
	assert := assert.New(t)
	v0 := Vector3{1.0, 2.0, 3.0}
	v1 := Vector3{1.0, 2.0, 3.0}
	assert.Truef(EqualVector3(v0, v1), "%v should equal to %v", v0, v1)

	v2 := Vector3{3.0, 2.0, 1.0}
	assert.Falsef(EqualVector3(v1, v2), "%v should not equal to %v", v1, v2)
}

func TestAddVector3(t *testing.T) {
	assert := assert.New(t)
	var result Vector3
	var expected Vector3
	v0 := Vector3{1.0, 1.0, 1.0}
	v1 := Vector3{1.0, 2.0, 3.0}
	result = AddVector3(v0, v1)
	expected = Vector3{2.0, 3.0, 4.0}
	assert.Truef(EqualVector3(result, expected), "%v + %v = %v", v0, v1, result)
}

func TestSubVector3(t *testing.T) {
	assert := assert.New(t)
	var result Vector3
	var expected Vector3
	v0 := Vector3{}
	v1 := Vector3{1.0, 1.0, 1.0}
	result = SubVector3(v0, v1)
	expected = Vector3{-1.0, -1.0, -1.0}
	assert.Truef(EqualVector3(result, expected), "%v - %v = %v", v0, v1, result)
}

func TestMulVector3(t *testing.T) {
	assert := assert.New(t)
	var result Vector3
	var expected Vector3
	v := Vector3{1.0, 2.0, 3.0}
	x := float32(2.0)
	result = MulVector3(x, v)
	expected = Vector3{2.0, 4.0, 6.0}
	assert.Truef(EqualVector3(result, expected), "%v * %v = %v", x, v, result)
}

func TestDivVector3(t *testing.T) {
	assert := assert.New(t)
	var result Vector3
	var expected Vector3
	v := Vector3{1.0, 2.0, 3.0}
	x := float32(2.0)
	result = DivVector3(v, x)
	expected = Vector3{0.5, 1.0, 1.5}
	assert.Truef(EqualVector3(result, expected), "%v / %v = %v", v, x, result)
}

func TestDotVector3(t *testing.T) {
	assert := assert.New(t)
	v0 := Vector3{}
	v1 := Vector3{1.0, 1.0, 1.0}
	assert.Truef(Equal32(0.0, DotVector3(v0, v1)), "Dot of %v and %v should be 0.0", v0, v1)

	v2 := Vector3{2.0, 2.0, 2.0}
	assert.Truef(Equal32(6.0, DotVector3(v1, v2)), "Dot of %v and %v should be 6.0", v1, v2)
}

func TestCrossVector3(t *testing.T) {
	assert := assert.New(t)
	v0 := Vector3{1.0, 0.0, 0.0}
	v1 := Vector3{0.0, 1.0, 0.0}
	result := CrossVector3(v0, v1)
	expected := Vector3{0.0, 0.0, 1.0}
	assert.Truef(EqualVector3(result, expected), "Cross of %v and %v should %v", v0, v1, result)
}

func TestNormalizeVector3(t *testing.T) {
	assert := assert.New(t)
	v := Vector3{1.0, 1.0, 1.0}
	result := NormalizeVector3(v)
	expected := Vector3{0.57735026, 0.57735026, 0.57735026}
	assert.Truef(EqualVector3(result, expected), "Normalized %v should be %v", v, result)
}

func TestSaturateVector3(t *testing.T) {
	assert := assert.New(t)
	v := Vector3{1.1, -1.0, 0.0}
	result := SaturateVector3(v)
	expected := Vector3{1.0, 0.0, 0.0}
	assert.Truef(EqualVector3(result, expected), "Saturated %v should be %v", v, result)
}

func TestRandomOnHemiSphere(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<100; i++ {
		v := RandomOnHemiSphere(rand.Float32(), rand.Float32())
		l := v.Length()
		assert.Truef(Equal32(l, 1.0) && 0.0<=v.Z && v.Z<=1.0, "%v (%v) should be on hemisphere", v, l)
	}
}
/*
func TestReflect(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<100; i++ {
		x := RandomOnHemiSphere(rand.Float32(), rand.Float32())
		x = SubVector3(Vector3{}, x)
		n := Vector3{0.0, 0.0, 1.0}
		r := Reflect(x, n)
		d := DotVector3(x, r)
		assert.Truef(0.0<=d, "%v reflected to %v (%v)", r, x, d)
	}
}
*/

