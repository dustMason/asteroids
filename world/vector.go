package world

import (
	"math"
	"math/rand"
)

type vector struct {
	x, y float64
}

func (v *vector) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *vector) resize(l float64) vector {
	return v.multiply(l / v.length())
}

func (v *vector) multiply(n float64) vector {
	return vector{x: v.x * n, y: v.y * n}
}

func (v *vector) add(vv vector) vector {
	return vector{x: v.x + vv.x, y: v.y + vv.y}
}

func (v *vector) rotate(degrees float64) {
	radians := degrees * math.Pi / 180
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	v.x = v.x*cos - v.y*sin
	v.y = v.x*sin + v.y*cos
}

func randomVector(magnitude float64) vector {
	phi := rand.Float64() * math.Pi * 2.
	return vector{x: math.Cos(phi) * magnitude, y: math.Sin(phi) * magnitude}
}
