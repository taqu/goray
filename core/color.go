package core

import (
	"git.maze.io/go/math32"
)

type Color32 struct {
	R float32
	G float32
	B float32
	A float32
}

func AddColor32(c0, c1 Color32) Color32 {
	c0.R += c1.R
	c0.G += c1.G
	c0.B += c1.B
	c0.A += c1.A
	return c0
}

func MulColor32(x float32, c Color32) Color32 {
	return Color32{x*c.R, x*c.G, x*c.B, x*c.A}
}

func DivColor32(c Color32, x float32) Color32 {
	return Color32{c.R / x, c.G / x, c.B / x, c.A / x}
}

func sRGBToLinear(x float32) float32 {
	if 0.0 <= x && x <= 0.04045 {
		return x / 12.92
	} else if 0.04045 < x && x <= 1.0 {
		return math32.Pow((x+0.055)/1.055, 2.4)
	} else {
		return x
	}
}

func linearToSRGB(x float32) float32 {
	if 0.0 <= x && x <= 0.0031308 {
		return x * 12.92
	} else if 0.0031308 < x && x <= 1.0 {
		return math32.Pow(1.055*x, 1.0/2.4) - 0.055
	} else {
		return x
	}
}

func SRGBToLinear(c Color32) Color32 {
	return Color32{
		sRGBToLinear(c.R),
		sRGBToLinear(c.G),
		sRGBToLinear(c.B),
		c.A}
}

func LinearToSRGB(c Color32) Color32 {
	return Color32{
		linearToSRGB(c.R),
		linearToSRGB(c.G),
		linearToSRGB(c.B),
		c.A}
}
