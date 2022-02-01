package gol

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOfLife struct {
	cells      [][]int
	generation int

	cols         int
	rows         int
	windowWidth  int
	windowHeight int
	blockSize    int
}

func NewGame(windowWidth, windowsHeight, blockSize int) *GameOfLife {
	cells := [][]int{}
	rows := windowsHeight / blockSize
	cols := windowWidth / blockSize

	for i := 0; i < rows; i++ {
		row := []int{}
		for j := 0; j < cols; j++ {
			row = append(row, rand.Intn(2))
		}
		cells = append(cells, row)
	}

	return &GameOfLife{
		cells:        cells,
		windowWidth:  windowWidth,
		windowHeight: windowsHeight,
		blockSize:    blockSize,
		rows:         rows,
		cols:         cols,
	}
}

func (g *GameOfLife) Render() {
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			if g.cells[y][x] == 1 {
				rl.DrawRectangle(int32(x*g.blockSize), int32(y*g.blockSize), int32(g.blockSize), int32(g.blockSize), rl.Green)
			}
		}
	}

	rl.DrawText(fmt.Sprintf("Generation: %d", g.generation), 10, 10, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("FPS: %.1f", rl.GetFPS()), 10, 30, 20, rl.Black)
}

func (g *GameOfLife) GameLoop() {
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Red)

		mousePos := rl.GetMousePosition()
		x := int(mousePos.X / float32(g.blockSize))
		y := int(mousePos.Y / float32(g.blockSize))

		rl.DrawRectangle(int32(x*g.blockSize), int32(y*g.blockSize), int32(g.blockSize), int32(g.blockSize), rl.Purple)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			g.cells[y][x] = 1
		}

		g.Render()
		g.Update()

		rl.EndDrawing()

		g.generation++
	}

}

func (g *GameOfLife) Update() {
	newCells := [][]int{}

	for y := 0; y < g.rows; y++ {
		row := []int{}

		for x := 0; x < g.cols; x++ {
			row = append(row, g.updateCell(y, x))
		}

		newCells = append(newCells, row)
	}

	g.cells = newCells
}

func (g *GameOfLife) updateCell(i, j int) int {
	var count int

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}

			if i+x < 0 || i+x >= g.rows {
				continue
			}

			if j+y < 0 || j+y >= g.cols {
				continue
			}

			if g.cells[i+x][j+y] == 1 {
				count++
			}
		}
	}

	if g.cells[i][j] == 1 {
		if count < 2 || count > 3 {
			return 0
		} else {
			return 1
		}
	} else {
		if count == 3 {
			return 1
		} else {
			return 0
		}
	}
}
