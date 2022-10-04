package world

import (
	"github.com/dustmason/asteroids/util"
	"image"
	"sync"
)

type player struct {
	sync.RWMutex
	body
	id      string
	name    string
	onDeath func()

	// counters
	health    int
	maxHealth int
}

func newPlayer(id, name string, pos vector) *player {
	// flying V shape, pointed down.
	points := []image.Point{
		{0, 0},
		{-2, -2},
		{0, 4},
		{2, -2},
		{0, 0},
	}
	return &player{
		id:        id,
		name:      name,
		health:    100,
		maxHealth: 100,
		body: body{
			form:    pointsToVectors(points),
			pos:     &pos,
			char:    '%',
			heading: &vector{0, 1},
			vec:     &vector{0, 0},
		},
	}
}

func (p *player) thrust(amount float64) {
	nv := p.body.vec.add(p.body.heading.multiply(amount))
	p.body.vec = &nv
	p.body.vec.x = util.ClampedFloat64(p.body.vec.x, -1, 1)
	p.body.vec.y = util.ClampedFloat64(p.body.vec.y, -1, 1)
}
