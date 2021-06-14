package core
import "testing"

func TestEqual32(t *testing.T) {
	if !Equal32(0.0, Epsilon32) {
		t.Fatal("Epsilon should be equal to zero.")
	}
	if Equal32(0.0, Epsilon32*2.0) {
		t.Fatal("Epsilon shoudl be larger than zero")
	}
}

func TestEqual64(t *testing.T) {
	if !Equal64(0.0, Epsilon64) {
		t.Fatal("Epsilon should be equal to zero.")
	}
	if Equal64(0.0, Epsilon64*2.0) {
		t.Fatal("Epsilon shoudl be larger than zero")
	}
}

func TestSaturate32(t *testing.T) {
	x0 := Saturate32(-0.1)
	if x0<0.0 || 1.0<x0 {
		t.Fatal("Saturated should be between 0 and 1")
	}
	x1 := Saturate32(1.1)
	if x1<0.0 || 1.0<x1 {
		t.Fatal("Saturated should be between 0 and 1")
	}
	x2 := Saturate32(0.5)
	if x2<0.0 || 1.0<x2 {
		t.Fatal("Saturated should be between 0 and 1")
	}
}

func TestSaturate64(t *testing.T) {
	x0 := Saturate64(-0.1)
	if x0<0.0 || 1.0<x0 {
		t.Fatal("Saturated should be between 0 and 1")
	}
	x1 := Saturate64(1.1)
	if x1<0.0 || 1.0<x1 {
		t.Fatal("Saturated should be between 0 and 1")
	}
	x2 := Saturate64(0.5)
	if x2<0.0 || 1.0<x2 {
		t.Fatal("Saturated should be between 0 and 1")
	}
}
