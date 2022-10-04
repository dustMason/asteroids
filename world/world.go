package world

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	vv "golang.org/x/image/vector"
	"image"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type World struct {
	sync.RWMutex
	W, H      int
	players   map[string]*player // map of player id => entity that points to that player
	asteroids []*asteroid
	bullets   []*bullet
	// lastTick   time.Time
}

func (w *World) Tick(t time.Time) {
	// todo remove asteroids that are off the board
	// todo remove any bullets that are off the board
	// todo break any asteroids that have collided with a bullet
	// todo make any players that have hit the edge of the world wrap around

	for len(w.asteroids) < 6 {
		w.asteroids = append(w.asteroids, newAsteroid(w.randomPointOnEdge(), randomVector(0.2)))
	}
	for _, a := range w.asteroids {
		a.tick()
	}
	for _, p := range w.players {
		p.tick()
	}
	for i, b := range w.bullets {
		fmt.Println("bullet", i, b.pos, b.vec)
		b.tick()
	}
}

func (w *World) DisconnectPlayer(playerID string) {
	w.Lock()
	defer w.Unlock()
	delete(w.players, playerID)
}

func (w *World) PlayerJoin(playerID string, playerName string) {
	w.Lock()
	defer w.Unlock()
	p, ok := w.players[playerID]
	if !ok {
		// x, y, _ := w.randomAvailableCoord()
		x := float64(w.W) / 2.
		y := float64(w.H) / 2.
		p = newPlayer(playerID, playerName, vector{x, y})
		w.players[playerID] = p
	}
	p.name = playerName
}

func (w *World) OnPlayerDeath(key string, f func()) {
	// todo save `f` and run it when this player dies
}

func ansiGrey(a uint16) string {
	// return a value from 232 (black) to 255 (white)
	return strconv.Itoa(int(a)/2731 + 232)
}

func ansiBlue(a uint16) string {
	// return a value from 16 (black) to 21 (blue)
	return strconv.Itoa(int(a)/13107 + 16)
}

// idea: zoom the resolution based on terminal size

const asciiArt = ".++%"

func (w *World) Render(id string, name string, width int, height int) string {
	r := vv.NewRasterizer(width, height)
	for _, p := range w.players {
		p.Draw(r)
	}
	for _, a := range w.asteroids {
		a.Draw(r)
	}
	for _, b := range w.bullets {
		b.Draw(r)
	}
	dst := image.NewAlpha(image.Rect(0, 0, width, height))
	r.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})
	var out strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := dst.RGBA64At(x, y)
			char := asciiArt[uint8(c.A/256)>>6]
			if c.A > 0 {
				// clr, _ := colorful.MakeColor(c)
				out.WriteString(
					// todo fade down the grey slightly to allow the ascii to show better
					lipgloss.NewStyle().
						Background(lipgloss.Color(ansiGrey(c.A))).
						Foreground(lipgloss.Color("#ffffff")).
						Render(string(char)),
				)
			} else {
				out.WriteString(" ")
			}
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (w *World) AcceleratePlayer(i int, id string) {
	if p, ok := w.players[id]; ok {
		p.thrust(0.05 * float64(i))
	}
}

func (w *World) RotatePlayer(i int, id string) {
	if p, ok := w.players[id]; ok {
		p.rotate(float64(i) * 5)
	}
}

func (w *World) FirePlayer(playerID string) {
	if p, ok := w.players[playerID]; ok {
		vec := vector{x: p.heading.x, y: p.heading.y}
		vec = vec.resize(2.)
		vec = vec.add(*p.vec)
		w.bullets = append(w.bullets, newBullet(*p.pos, vec))
	}
}

func (w *World) RenderWorldStatus() string {
	// todo ? scores?
	return ""
}

func (w *World) RenderPosition(id string) string {
	// todo render a string of this users position?
	return ""
}

func (w *World) inBounds(x int, y int) bool {
	return x >= 0 && x < w.W && y >= 0 && y < w.H
}

func (w *World) randomPointOnEdge() vector {
	pick := rand.Intn(3)
	switch pick {
	case 0:
		return vector{x: rand.Float64() * float64(w.W), y: 0}
	case 1:
		return vector{x: float64(w.W), y: rand.Float64() * float64(w.H)}
	case 2:
		return vector{x: rand.Float64() * float64(w.W), y: float64(w.H)}
	}
	return vector{x: 0, y: rand.Float64() * float64(w.H)}
}

// func (w *World) randomAvailableCoord() (interface{}, interface{}, interface{}) {
// 	tries := 1000
// 	for tries > 0 {
// 		x := rand.Intn(w.W)
// 		y := rand.Intn(w.H)
// 		if w.walkable(x, y) && !w.occupied(x, y) {
// 			return x, y, nil
// 		}
// 		tries--
// 	}
// 	return 0, 0, errors.New("couldn't place random coord")
// }

func NewWorld(width, height int) *World {
	w := &World{
		W:       width,
		H:       height,
		players: make(map[string]*player),
	}
	go w.runTicker()
	return w
}

func (w *World) runTicker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for t := range ticker.C {
		w.Tick(t)
	}
}
