package world

import (
	vv "golang.org/x/image/vector"
	"image"
	"image/draw"
	"math"
)

type body struct {
	form    []vector // a list of points around a closed shape
	pos     *vector
	heading *vector
	vec     *vector
	char    rune
}

func (b *body) Draw(r *vv.Rasterizer) {
	r.DrawOp = draw.Src
	for _, v := range b.form {
		r.LineTo(float32(v.x+b.pos.x), float32(v.y+b.pos.y))
	}
	r.ClosePath()
}

func (b *body) Char() rune {
	return b.char
}

func (b *body) rotate(degrees float64) {
	origin := b.form[0]
	radians := degrees * math.Pi / 180
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	newForm := make([]vector, len(b.form))
	newForm[0] = origin
	for i, fp := range b.form[1:] {
		x2 := fp.x - origin.x
		y2 := fp.y - origin.y
		newForm[i+1] = vector{
			x: x2*cos - y2*sin + origin.x,
			y: x2*sin + y2*cos + origin.y,
		}
	}
	b.form = newForm
	b.heading.rotate(degrees)
}

func (b *body) tick() {
	newPos := b.pos.add(*b.vec)
	b.pos = &newPos
}

func (b *body) collidesWith(ob *body) bool {
	// do these shapes intersect?
	return false
}

func pointsToVectors(points []image.Point) []vector {
	vecs := make([]vector, len(points))
	for i, p := range points {
		vecs[i] = vector{x: float64(p.X), y: float64(p.Y)}
	}
	return vecs
}
