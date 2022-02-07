package core

import (
	"math"
	"git.maze.io/go/math32"
)

const (
	Epsilon32  float32 = 1.0e-6
	Epsilon64  float64 = 1.0e-14
	Infinity32 float32 = 1.0e37
	DegToRad32 float32 = float32(1.57079632679489661923 / 90.0)
	RadToDeg32 float32 = float32(90.0 / 1.57079632679489661923)
)

func Equal32(x0, x1 float32) bool {
	return math32.Abs(x0-x1) <= Epsilon32
}

func EqualZero32(x float32) bool {
	return math32.Abs(x) <= Epsilon32
}

func Equal64(x0, x1 float64) bool {
	return math.Abs(x0-x1) <= Epsilon64
}

func EqualZero64(x float64) bool {
	return math.Abs(x) <= Epsilon64
}

func Saturate32(x float32) float32 {
	if x < 0.0 {
		return 0.0
	} else if 1.0 < x {
		return 1.0
	}
	return x
}

func Saturate64(x float64) float64 {
	if x < 0.0 {
		return 0.0
	} else if 1.0 < x {
		return 1.0
	}
	return x
}

func Lerp32(x0, x1, t float32) float32 {
	return x0 * (1.0-t) + x1 * t;
}

func Clamp0132(x float32) float32 {
	if x<0 {
		return 0
	}else if 1<x {
		return 1
	}else{
		return x
	}
}

func Schlick(cosine, refIndex float32) float32 {
	r0 := (1.0 - refIndex) / (1.0 + refIndex)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math32.Pow((1.0-cosine), 5.0)
}
