package world

import "image"

type bullet struct {
	body
}

func newBullet(pos vector, vec vector) *bullet {
	// square shape
	points := []image.Point{
		{0, 0},
		{1, 0},
		{1, -1},
		{0, 1},
		{0, 0},
	}
	return &bullet{
		body: body{
			form:    pointsToVectors(points),
			pos:     &pos,
			char:    'o',
			heading: &vector{0, 1},
			vec:     &vec,
		},
	}
}
