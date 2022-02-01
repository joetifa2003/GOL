package main

import (
	"GOL/gol"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 500
const WINDOW_HEIGHT = 500
const BLOCK_SIZE = 20

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Game of life")
	rl.SetTargetFPS(30)

	game := gol.NewGame(WINDOW_WIDTH, WINDOW_HEIGHT, BLOCK_SIZE)
	game.GameLoop()
}
