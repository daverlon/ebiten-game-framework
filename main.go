package main

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var showdebuginfo = false

var (
	clrRed    = color.RGBA{255, 0, 0, 255}
	clrCyan   = color.RGBA{0, 255, 255, 255}
	clrYellow = color.RGBA{255, 255, 0, 255}
)

const (
	// 640 x 480
	viewportw       = 320 // how much of the 'world' the camera can see
	viewporth       = 240
	windowcenterx   = viewportw / 2 // 'screen' position
	windowcentery   = viewporth / 2
	originalcamzoom = 1.0
)

func DrawOutlineRect(img *ebiten.Image, x float64, y float64, w float64, h float64) {
	// draw red outline (debug)
	ebitenutil.DrawLine(img, x, y, x+w, y, clrRed)
	ebitenutil.DrawLine(img, x+w, y, x+w, y+h, clrRed)
	ebitenutil.DrawLine(img, x, y+h, x+w, y+h, clrRed)
	ebitenutil.DrawLine(img, x, y, x, y+h, clrRed)
}

func updateCursor(g *Game) {
	x, y := ebiten.CursorPosition()
	g.cursor.screenx, g.cursor.screeny = float64(x), float64(y)
	g.cursor.worldx, g.cursor.worldy = g.cam.ScreenToWorld(float64(g.cursor.screenx), float64(g.cursor.screeny))
}

type Game struct {
	scenes SceneStack
	cursor Cursor
	cam    Camera

	sprites []*Sprite
}

func (g *Game) Update() error {

	if GameInstance.scenes.IsEmpty() {
		fmt.Println("Warning: SceneStack is empty.")
		return errors.New("error: no scene available. quitting")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		showdebuginfo = !showdebuginfo
	}

	_, mY := ebiten.Wheel()
	GameInstance.cam.zoom += mY
	minzoom := 0.2
	if GameInstance.cam.zoom < minzoom {
		GameInstance.cam.zoom = minzoom
	}

	updateCursor(g)

	// update current scene
	if !GameInstance.scenes.IsEmpty() {
		GameInstance.scenes.Active().Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{50, 50, 50, 255})
	screen.Clear()

	if !g.scenes.IsEmpty() && GameInstance.scenes.Active().Draw != nil {
		GameInstance.scenes.Active().Draw(screen)
	}

	ebitenutil.DrawRect(screen, 0, 0, 47, 15, color.Black)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %1.2f", ebiten.ActualTPS()))

	// debug info
	if showdebuginfo {

		ebitenutil.DrawRect(screen, float64(g.cursor.screenx)-2, float64(g.cursor.screeny)-2, 4, 4, clrCyan)

		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("(%d, %d)", int(g.cursor.screenx), int(g.cursor.screeny)),
			int(g.cursor.screenx), int(g.cursor.screeny))
		// scene info (to avoid crashing)
		sceneName := "nil"
		sceneIndex := len(g.scenes) - 1
		sceneCount := len(g.scenes)
		if !g.scenes.IsEmpty() {
			sceneName = g.scenes[len(g.scenes)-1].name
		}

		debugString := fmt.Sprintf(
			"\n"+
				"Memory\n"+
				"Sprite Count: %d\n"+
				//"\n"+
				"Scene\n"+
				"Index: %d (%d)\n"+
				"Name: %s\n"+
				//"\n"+
				"Camera\n"+
				"Pos: (%d, %d)\n"+
				"Zoom: %0.1f -> %0.1f\n"+
				//"\n"+
				"Mouse\n"+
				"ScreenPos: (%d, %d)\n"+
				"WorldPos: (%d, %d)",
			len(GameInstance.sprites),
			sceneIndex, sceneCount,
			sceneName,
			int(g.cam.x), int(g.cam.y),
			originalcamzoom, g.cam.zoom,
			int(g.cursor.screenx), int(g.cursor.screeny),
			int(g.cursor.worldx), int(g.cursor.worldy),
		)
		ebitenutil.DebugPrint(screen, debugString)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return viewportw, viewporth
}

// instance of 'game' with global access
var GameInstance = &Game{
	cam: Camera{zoom: originalcamzoom},
}

func main() {

	s := 1.5
	w, h := int(640*s), int(480*s)

	// ebiten.SetVsyncEnabled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Game Window")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	GameInstance.scenes.Push(testScene2())

	if err := ebiten.RunGame(GameInstance); err != nil {
		log.Fatal(err)
	}
}
