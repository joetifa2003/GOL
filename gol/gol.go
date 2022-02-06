package gol

import (
	"fmt"
	"math/rand"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOfLife struct {
	cells [][]int

	cols         int
	rows         int
	windowWidth  int32
	windowHeight int32
	blockSize    int32

	mouseCellX int32
	mouseCellY int32

	running  bool
	hideMenu bool

	FPS        int
	camera     rl.Camera2D
	mouseDelta rl.Vector2
}

func NewGame(windowWidth, windowsHeight, blockSize, FPS int) *GameOfLife {
	cells := [][]int{}
	rows := windowsHeight / blockSize
	cols := windowWidth / blockSize

	for i := 0; i < rows; i++ {
		row := []int{}
		for j := 0; j < cols; j++ {
			if rand.Intn(2) == 1 {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		cells = append(cells, row)
	}

	centerOfScreen := rl.NewVector2(float32(windowWidth/2), float32(windowsHeight/2))

	return &GameOfLife{
		cells:        cells,
		windowWidth:  int32(windowWidth),
		windowHeight: int32(windowsHeight),
		blockSize:    int32(blockSize),
		rows:         rows,
		cols:         cols,
		running:      false,
		hideMenu:     false,
		FPS:          FPS,
		camera: rl.Camera2D{
			Target:   centerOfScreen,
			Offset:   centerOfScreen,
			Zoom:     1,
			Rotation: 0,
		},
	}
}

func (g *GameOfLife) Draw() {
	rl.DrawRectangle(g.mouseCellX*g.blockSize, g.mouseCellY*g.blockSize, g.blockSize, g.blockSize, rl.Fade(rl.Black, 0.75))

	cameraWorld := rl.NewVector2(
		g.camera.Target.X-g.camera.Offset.X,
		g.camera.Target.Y-g.camera.Offset.Y,
	)

	arrayStartX := Max(int(int32(cameraWorld.X)/g.blockSize), 0)
	arrayStartY := Max(int(int32(cameraWorld.Y)/g.blockSize), 0)

	arrayEndX := Min(int(int32(cameraWorld.X+float32(g.windowWidth))/g.blockSize), g.cols)
	arrayEndY := Min(int(int32(cameraWorld.Y+float32(g.windowHeight))/g.blockSize), g.rows)

	for y := arrayStartY; y < arrayEndY; y++ {
		for x := arrayStartX; x < arrayEndX; x++ {
			if g.cells[y][x] == 1 {
				rl.DrawRectangle(int32(x)*g.blockSize, int32(y)*g.blockSize, g.blockSize, g.blockSize, rl.Black)
			}
		}
	}

	borderThickness := 6
	rl.DrawRectangleLinesEx(
		rl.NewRectangle(-float32(borderThickness),
			-float32(borderThickness),
			float32(g.windowWidth+int32(borderThickness*2)),
			float32(g.windowHeight+int32(borderThickness*2)),
		),
		float32(borderThickness),
		rl.Fade(rl.Red, 0.5),
	)
}

func (g *GameOfLife) DrawGUI() {
	rl.DrawText(fmt.Sprintf("FPS: %0.f", rl.GetFPS()), 10, 10, 20, rl.Orange)

	if !g.running {
		if !g.hideMenu {
			rl.DrawRectangle(0, 0, g.windowWidth, g.windowHeight, rl.Fade(rl.Black, 0.8))
			pressStartWidth := rl.MeasureText("Press SPACE to start", 50)
			rl.DrawText("Press SPACE to start", g.windowWidth/2-pressStartWidth/2, g.windowHeight/2-150, 50, rl.White)

			pressCWidth := rl.MeasureText("Press C to clear", 50)
			rl.DrawText("Press C to clear", g.windowWidth/2-pressCWidth/2, g.windowHeight/2-100, 50, rl.White)

			pressRWidth := rl.MeasureText("Press R to randomize", 50)
			rl.DrawText("Press R to randomize", g.windowWidth/2-pressRWidth/2, g.windowHeight/2-50, 50, rl.White)

			pressHWidth := rl.MeasureText("Press H to toggle menu", 50)
			rl.DrawText("Press H to toggle menu", g.windowWidth/2-pressHWidth/2, g.windowHeight/2, 50, rl.White)

			pressWWidth := rl.MeasureText("Press W to increase FPS", 50)
			rl.DrawText("Press W to increase FPS", g.windowWidth/2-pressWWidth/2, g.windowHeight/2+100, 50, rl.White)

			pressSWidth := rl.MeasureText("Press S to decrease FPS", 50)
			rl.DrawText("Press S to decrease FPS", g.windowWidth/2-pressSWidth/2, g.windowHeight/2+150, 50, rl.White)
		}
	}
}

func (g *GameOfLife) GameLoop() {
	prevMousePos := rl.GetMousePosition()
	for !rl.WindowShouldClose() {
		currentMousePos := rl.GetMousePosition()

		g.mouseDelta = rl.Vector2Subtract(prevMousePos, currentMousePos)
		prevMousePos = currentMousePos

		g.Input()
		if g.running {
			g.Update()
		}

		// ------------
		rl.BeginDrawing()

		rl.ClearBackground(rl.White)

		rl.BeginMode2D(g.camera)
		g.Draw()
		rl.EndMode2D()

		g.DrawGUI()

		rl.EndDrawing()
		// ------------
	}

}

func (g *GameOfLife) Update() {
	newCells := make([][]int, g.rows)
	for i := range newCells {
		newCells[i] = make([]int, g.cols)
	}

	var wg sync.WaitGroup

	for y := 0; y < g.rows; y++ {
		wg.Add(1)
		go g.updateCell(y, &newCells, &wg)
	}

	wg.Wait()

	g.cells = newCells
}

func (g *GameOfLife) Input() {
	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		g.camera.Target = rl.Vector2Add(g.camera.Target, g.mouseDelta)
	}

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

	// Handle camera zoom
	g.camera.Zoom += float32(rl.GetMouseWheelMove()) * 0.1
	g.camera.Zoom = MaxFloat32(g.camera.Zoom, 1)

	// Mouse position
	mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), g.camera)
	g.mouseCellX = int32(Max(Min(int(mousePos.X/float32(g.blockSize)), g.cols-1), 0))
	g.mouseCellY = int32(Max(Min(int(mousePos.Y/float32(g.blockSize)), g.rows-1), 0))

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && !rl.IsKeyDown(rl.KeyLeftShift) {
		g.cells[g.mouseCellY][g.mouseCellX] = 1
	} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.cells[g.mouseCellY][g.mouseCellX] = 0
	}
}

func (g *GameOfLife) updateCell(y int, newCells *[][]int, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := 0; x < g.cols; x++ {
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
				(*newCells)[y][x] = 0
			} else {
				(*newCells)[y][x] = 1
			}
		} else {
			if neighborCount == 3 {
				(*newCells)[y][x] = 1
			} else {
				(*newCells)[y][x] = 0
			}
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
