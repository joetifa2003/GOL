package gol

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOfLife struct {
	cells [][]int

	cols         int
	rows         int
	windowWidth  int
	windowHeight int
	blockSize    int

	mouseCellX int
	mouseCellY int

	running  bool
	hideMenu bool

	FPS int
}

func NewGame(windowWidth, windowsHeight, blockSize, FPS int) *GameOfLife {
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
		running:      false,
		hideMenu:     false,
		FPS:          FPS,
	}
}

func (g *GameOfLife) Draw() {
	rl.DrawRectangle(int32(g.mouseCellX*g.blockSize), int32(g.mouseCellY*g.blockSize), int32(g.blockSize), int32(g.blockSize), rl.Fade(rl.Black, 0.75))

	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			if g.cells[y][x] == 1 {
				rl.DrawRectangle(int32(x*g.blockSize), int32(y*g.blockSize), int32(g.blockSize), int32(g.blockSize), rl.Black)
			}
		}
	}

	if !g.running {
		if !g.hideMenu {
			rl.DrawRectangle(0, 0, int32(g.windowWidth), int32(g.windowHeight), rl.Fade(rl.Black, 0.8))
			pressStartWidth := rl.MeasureText("Press SPACE to start", 50)
			rl.DrawText("Press SPACE to start", int32(g.windowWidth/2)-(pressStartWidth/2), int32(g.windowHeight/2)-150, 50, rl.White)

			pressCWidth := rl.MeasureText("Press C to clear", 50)
			rl.DrawText("Press C to clear", int32(g.windowWidth/2)-(pressCWidth/2), int32(g.windowHeight/2)-100, 50, rl.White)

			pressRWidth := rl.MeasureText("Press R to randomize", 50)
			rl.DrawText("Press R to randomize", int32(g.windowWidth/2)-(pressRWidth/2), int32(g.windowHeight/2)-50, 50, rl.White)

			pressHWidth := rl.MeasureText("Press H to toggle menu", 50)
			rl.DrawText("Press H to toggle menu", int32(g.windowWidth/2)-(pressHWidth/2), int32(g.windowHeight/2), 50, rl.White)

			pressWWidth := rl.MeasureText("Press W to increase FPS", 50)
			rl.DrawText("Press W to increase FPS", int32(g.windowWidth/2)-(pressWWidth/2), int32(g.windowHeight/2)+100, 50, rl.White)

			pressSWidth := rl.MeasureText("Press S to decrease FPS", 50)
			rl.DrawText("Press S to decrease FPS", int32(g.windowWidth/2)-(pressSWidth/2), int32(g.windowHeight/2)+150, 50, rl.White)
		}
	} else {
		rl.DrawText(fmt.Sprintf("FPS: %.1f", rl.GetFPS()), 10, 10, 20, rl.Orange)
	}
}

func (g *GameOfLife) GameLoop() {
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.White)

		g.Draw()
		if g.running {
			g.Update()
		}
		g.Input()

		rl.EndDrawing()
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

func (g *GameOfLife) Input() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.running = !g.running
	}

	if rl.IsKeyPressed(rl.KeyC) {
		g.ClearCells()
	}

	if rl.IsKeyPressed(rl.KeyR) {
		g.RandomCells()
	}

	if rl.IsKeyPressed(rl.KeyH) {
		g.hideMenu = !g.hideMenu
	}

	if rl.IsKeyPressed(rl.KeyW) {
		g.FPS += 10
		rl.SetTargetFPS(int32(g.FPS))
	}

	if rl.IsKeyPressed(rl.KeyS) {
		g.FPS -= 10
		g.FPS = Max(g.FPS, 10)
		rl.SetTargetFPS(int32(g.FPS))
	}

	// Mouse position
	mousePos := rl.GetMousePosition()
	g.mouseCellX = Max(Min(int(mousePos.X/float32(g.blockSize)), g.cols-1), 0)
	g.mouseCellY = Max(Min(int(mousePos.Y/float32(g.blockSize)), g.rows-1), 0)

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		g.cells[g.mouseCellY][g.mouseCellX] = 1
	} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.cells[g.mouseCellY][g.mouseCellX] = 0
	}
}

func (g *GameOfLife) updateCell(y, x int) int {
	var neighborCount int

	for dirX := -1; dirX <= 1; dirX++ {
		for dirY := -1; dirY <= 1; dirY++ {
			if dirX == 0 && dirY == 0 {
				continue
			}

			neighborX := x + dirX
			neighborY := y + dirY

			if neighborX < 0 {
				neighborX += g.cols
			}

			if neighborX >= g.cols {
				neighborX -= g.cols
			}

			if neighborY < 0 {
				neighborY += g.rows
			}

			if neighborY >= g.rows {
				neighborY -= g.rows
			}

			if g.cells[neighborY][neighborX] == 1 {
				neighborCount++
			}
		}
	}

	if g.cells[y][x] == 1 {
		if neighborCount < 2 || neighborCount > 3 {
			return 0
		} else {
			return 1
		}
	} else {
		if neighborCount == 3 {
			return 1
		} else {
			return 0
		}
	}
}

func (g *GameOfLife) ClearCells() {
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			g.cells[y][x] = 0
		}
	}
}

func (g *GameOfLife) RandomCells() {
	g.ClearCells()

	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			if rand.Intn(2) == 1 {
				g.cells[y][x] = 1
			}
		}
	}
}
