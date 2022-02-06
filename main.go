package main

import (
	"GOL/gol"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Actual window width size
const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720

// Individual block size
const BLOCK_SIZE = 2

// Max fps
const FPS = 60

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Game of life")
	rl.SetTargetFPS(FPS)

	game := gol.NewGame(WINDOW_WIDTH, WINDOW_HEIGHT, BLOCK_SIZE, FPS)

	game.GameLoop()
}
