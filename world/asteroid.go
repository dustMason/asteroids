package world

import "image"

type asteroid struct {
	body
}

func newAsteroid(pos vector, vec vector) *asteroid {
	// wonky random thing shape
	points := []image.Point{
		{0, 0},
		{-2, 2},
		{0, 4},
		{1, 2},
		{0, 0},
	}
	return &asteroid{
		body: body{
			form:    pointsToVectors(points),
			pos:     &pos,
			char:    '%',
			heading: &vector{0, 1},
			vec:     &vec,
		},
	}
}
