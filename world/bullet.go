package world

import "image"

type bullet struct {
	body
	playerID string // playerID that this bullet belongs to
}

func newBullet(pos vector, vec vector, playerID string) *bullet {
	// square shape
	points := []image.Point{
		{0, 0},
		{1, 0},
		{1, -1},
		{0, 1},
		{0, 0},
	}
	return &bullet{
		playerID: playerID,
		body: body{
			form:    pointsToVectors(points),
			pos:     &pos,
			char:    'o',
			heading: &vector{0, 1},
			vec:     &vec,
		},
	}
}
