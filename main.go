package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	fractalMinX = -2.5
	fractalMaxX = 1.0
	fractalMinY = -1.0
	fractalMaxY = 1.0

	regionSize = 32
	numWorkers = 8
)

var (
	inColor  = RGBA{0, 0, 0, 255}
	outColor = RGBA{255, 255, 255, 255}
)

const (
	screenWidth  = 1366
	screenHeight = 768
)

type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Game struct {
	maxIterations int
	maxDistance   float64
	zoom          float64
	panX          float64
	panY          float64
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"Fractal",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	game := Game{1000, 2, 1, 0, 0}

	keyState := sdl.GetKeyboardState()

	GameLoop(game, keyState, renderer)

}

func GameLoop(game Game, keyState []uint8, renderer *sdl.Renderer) {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				os.Exit(1)
			}
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_PAGEUP] != 0 {
			game.maxIterations *= 2
			fmt.Println("iter up")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_PAGEDOWN] != 0 && game.maxIterations > 1 {
			game.maxIterations /= 2
			fmt.Println("iter down")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_KP_PLUS] != 0 {
			game.zoom *= 1.25
			fmt.Println("zoom in")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_KP_MINUS] != 0 && game.zoom > 1 {
			game.zoom /= 1.25
			fmt.Println("zoom out")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_UP] != 0 {
			game.panY -= 0.1 / game.zoom
			fmt.Println("up")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_DOWN] != 0 {
			game.panY += 0.1 / game.zoom
			fmt.Println("down")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_LEFT] != 0 {
			game.panX -= 0.1 / game.zoom
			fmt.Println("left")
		}
		if sdl.GetKeyboardState()[sdl.SCANCODE_RIGHT] != 0 {
			game.panX += 0.1 / game.zoom
			fmt.Println("right")
		}

		maxDistanceSquared := game.maxDistance * game.maxDistance

		for x := 0; x < screenWidth; x++ {
			for y := 0; y < screenHeight; y++ {

				x0, y0 := mapToFractalSpace(x, y, game)
				x1 := 0.0
				y1 := 0.0
				iteration := 0
				for x1*x1+y1*y1 <= maxDistanceSquared && iteration < game.maxIterations {
					xtemp := x1*x1 - y1*y1 + x0
					y1 = 2*x1*y1 + y0
					x1 = xtemp
					iteration++
				}

				drawPoint(x, y, mapColor(iteration, game.maxIterations), renderer)

			}
		}
		fmt.Println("Frame done")
		renderer.Present()
	}
}

func drawPoint(x, y int, rgba RGBA, renderer *sdl.Renderer) {
	renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, 255)
	renderer.DrawPoint(int32(x), int32(y))
}

func mapToFractalSpace(x, y int, game Game) (float64, float64) {
	scaleX := fractalMaxX - fractalMinX
	scaleY := fractalMaxY - fractalMinY
	mappedX := (float64(x)/float64(screenWidth))*scaleX + fractalMinX
	mappedY := (float64(y)/float64(screenHeight))*scaleY + fractalMinY
	mappedX = mappedX/game.zoom + game.panX
	mappedY = mappedY/game.zoom + game.panY
	return mappedX, mappedY
}

func mapColor(iteration, maxIterations int) RGBA {
	ratio := float64(iteration) / float64(maxIterations)
	invRatio := 1 - ratio
	return RGBA{
		uint8(float64(inColor.R)*ratio + float64(outColor.R)*invRatio),
		uint8(float64(inColor.G)*ratio + float64(outColor.G)*invRatio),
		uint8(float64(inColor.B)*ratio + float64(outColor.B)*invRatio),
		255,
	}
}
