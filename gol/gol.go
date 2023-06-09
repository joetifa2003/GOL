package gol

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720

const CELL_SIZE = 4

const ROWS = WINDOW_HEIGHT / CELL_SIZE
const COLS = WINDOW_WIDTH / CELL_SIZE

const FPS = 60

type GameOfLife struct {
	cells [ROWS][COLS]bool

	mouseCellX int32
	mouseCellY int32

	running  bool
	hideMenu bool

	FPS    int
	camera rl.Camera2D
}

func NewGame() *GameOfLife {
	cells := [ROWS][COLS]bool{}

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			cells[y][x] = rand.Intn(2) == 1
		}
	}

	centerOfScreen := rl.NewVector2(float32(WINDOW_WIDTH/2), float32(WINDOW_HEIGHT/2))

	return &GameOfLife{
		cells:    cells,
		running:  false,
		hideMenu: false,
		FPS:      FPS,
		camera: rl.Camera2D{
			Target:   centerOfScreen,
			Offset:   centerOfScreen,
			Zoom:     1,
			Rotation: 0,
		},
	}
}

func (g *GameOfLife) Init() {
	rand.Seed(time.Now().UnixNano())
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Game of life")
	rl.SetTargetFPS(FPS)
}

func (g *GameOfLife) Draw() {
	rl.DrawRectangle(g.mouseCellX*CELL_SIZE, g.mouseCellY*CELL_SIZE, CELL_SIZE, CELL_SIZE, rl.Fade(rl.Black, 0.75))

	cameraStart := rl.GetScreenToWorld2D(rl.NewVector2(0, 0), g.camera)
	cameraEnd := rl.GetScreenToWorld2D(rl.NewVector2(float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())), g.camera)
	cameraSize := rl.Vector2Subtract(cameraEnd, cameraStart)
	cameraRect := rl.NewRectangle(cameraStart.X, cameraStart.Y, cameraSize.X, cameraSize.Y)

	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			if g.cells[y][x] == true {
				cellX := int32(x) * CELL_SIZE
				cellY := int32(y) * CELL_SIZE

				if float32(cellX) >= cameraRect.X &&
					float32(cellX) <= cameraRect.X+cameraRect.Width &&
					float32(cellY) >= cameraRect.Y &&
					float32(cellY) <= cameraRect.Y+cameraRect.Height {
					rl.DrawRectangle(cellX, cellY, CELL_SIZE, CELL_SIZE, rl.Black)
				}
			}
		}
	}

	borderThickness := 6
	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			-float32(borderThickness),
			-float32(borderThickness),
			float32(int32(rl.GetScreenWidth())+int32(borderThickness*2)),
			float32(int32(rl.GetScreenHeight())+int32(borderThickness*2)),
		),
		float32(borderThickness),
		rl.Fade(rl.Red, 0.5),
	)
}

func (g *GameOfLife) DrawGUI() {
	rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 10, 20, rl.Orange)

	windowWidth := int32(rl.GetScreenWidth())
	windowHeight := int32(rl.GetScreenHeight())

	if !g.running {
		if !g.hideMenu {
			rl.DrawRectangle(0, 0, windowWidth, windowHeight, rl.Fade(rl.Black, 0.8))
			pressStartWidth := rl.MeasureText("Press SPACE to start", 50)
			rl.DrawText("Press SPACE to start", windowWidth/2-pressStartWidth/2, windowHeight/2-150, 50, rl.White)

			pressCWidth := rl.MeasureText("Press C to clear", 50)
			rl.DrawText("Press C to clear", windowWidth/2-pressCWidth/2, windowHeight/2-100, 50, rl.White)

			pressRWidth := rl.MeasureText("Press R to randomize", 50)
			rl.DrawText("Press R to randomize", windowWidth/2-pressRWidth/2, windowHeight/2-50, 50, rl.White)

			pressHWidth := rl.MeasureText("Press H to toggle menu", 50)
			rl.DrawText("Press H to toggle menu", windowWidth/2-pressHWidth/2, windowHeight/2, 50, rl.White)

			pressWWidth := rl.MeasureText("Press W to increase FPS", 50)
			rl.DrawText("Press W to increase FPS", windowWidth/2-pressWWidth/2, windowHeight/2+100, 50, rl.White)

			pressSWidth := rl.MeasureText("Press S to decrease FPS", 50)
			rl.DrawText("Press S to decrease FPS", windowWidth/2-pressSWidth/2, windowHeight/2+150, 50, rl.White)
		}
	}
}

func (g *GameOfLife) GameLoop() {
	for !rl.WindowShouldClose() {
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
	newCells := [ROWS][COLS]bool{}

	var wg sync.WaitGroup
	for y := 0; y < ROWS; y++ {
		wg.Add(1)
		go g.updateCell(y, &newCells, &wg)
	}
	wg.Wait()

	g.cells = newCells
}

func (g *GameOfLife) Input() {
	if rl.IsKeyDown(rl.KeyLeftShift) && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		g.camera.Target = rl.Vector2Add(
			g.camera.Target,
			rl.Vector2Scale(
				rl.GetMouseDelta(),
				-1/g.camera.Zoom,
			),
		)
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
	g.camera.Zoom += float32(rl.GetMouseWheelMove()) * 1
	g.camera.Zoom = MaxFloat32(g.camera.Zoom, 1)

	// Mouse position
	mousePos := rl.GetScreenToWorld2D(rl.GetMousePosition(), g.camera)
	g.mouseCellX = int32(Max(Min(int(mousePos.X/float32(CELL_SIZE)), COLS-1), 0))
	g.mouseCellY = int32(Max(Min(int(mousePos.Y/float32(CELL_SIZE)), ROWS-1), 0))

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && !rl.IsKeyDown(rl.KeyLeftShift) {
		g.cells[g.mouseCellY][g.mouseCellX] = true
	} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.cells[g.mouseCellY][g.mouseCellX] = false
	}
}

func (g *GameOfLife) updateCell(y int, newCells *[ROWS][COLS]bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := 0; x < COLS; x++ {
		var neighborCount int

		for dirX := -1; dirX <= 1; dirX++ {
			for dirY := -1; dirY <= 1; dirY++ {
				if dirX == 0 && dirY == 0 {
					continue
				}

				neighborX := x + dirX
				neighborY := y + dirY

				if neighborX < 0 {
					neighborX += COLS
				}

				if neighborX >= COLS {
					neighborX -= COLS
				}

				if neighborY < 0 {
					neighborY += ROWS
				}

				if neighborY >= ROWS {
					neighborY -= ROWS
				}

				if g.cells[neighborY][neighborX] == true {
					neighborCount++
				}
			}
		}

		if g.cells[y][x] == true {
			if neighborCount < 2 || neighborCount > 3 {
				(*newCells)[y][x] = false
			} else {
				(*newCells)[y][x] = true
			}
		} else {
			if neighborCount == 3 {
				(*newCells)[y][x] = true
			} else {
				(*newCells)[y][x] = false
			}
		}
	}
}

func (g *GameOfLife) ClearCells() {
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			g.cells[y][x] = false
		}
	}
}

func (g *GameOfLife) RandomCells() {
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			g.cells[y][x] = rand.Intn(2) == 1
		}
	}
}
