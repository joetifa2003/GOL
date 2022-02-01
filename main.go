package main

import (
	"GOL/gol"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720
const BLOCK_SIZE = 20
const FPS = 30

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Game of life")
	rl.SetTargetFPS(FPS)

	game := gol.NewGame(WINDOW_WIDTH, WINDOW_HEIGHT, BLOCK_SIZE)
	game.GameLoop()
}
