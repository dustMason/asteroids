package main

import (
	"fmt"
	"github.com/dustmason/asteroids/server"
	"github.com/dustmason/asteroids/world"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()
	width := 80
	height := 40
	w := world.NewWorld(width, height)
	s := server.NewServer(w)
	s.Listen()
}
