package main

import (
	"GOL/gol"
)

func main() {
	game := gol.NewGame()
	game.Init()
	game.GameLoop()
}
