package main

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var showdebuginfo = true

var (
	clrRed    = color.RGBA{255, 0, 0, 255}
	clrCyan   = color.RGBA{0, 255, 255, 255}
	clrYellow = color.RGBA{255, 255, 0, 255}
)

const (
	windowcenterx   = 320
	windowcentery   = 240.0
	originalcamzoom = 2.5
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
		return errors.New("Error: No scene available. Quitting")
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

	if !g.scenes.IsEmpty() && GameInstance.scenes.Active().Draw != nil {
		screen.DrawImage(GameInstance.scenes.Active().Draw(), nil)
	}

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

		ebitenutil.DrawRect(screen, 0, 0, 47, 15, color.Black)

		debugString := fmt.Sprintf(
			"FPS: %d\n"+
				"\n"+
				"Memory\n"+
				"Sprite Count: %d\n"+
				"\n"+
				"Scene\n"+
				"Index: %d (%d)\n"+
				"Name: %s\n"+
				"\n"+
				"Camera\n"+
				"Pos: (%d, %d)\n"+
				"Zoom: %0.1f -> %0.1f\n"+
				"\n"+
				"Mouse\n"+
				"ScreenPos: (%d, %d)\n"+
				"WorldPos: (%d, %d)",
			int(ebiten.ActualFPS()),
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
	return 640, 480
}

// instance of 'game' with global access
var GameInstance = &Game{
	cam: Camera{zoom: originalcamzoom},
}

func (g *Game) AddSprite(s *Sprite) {
	g.sprites = append(g.sprites, s)
}

func (g *Game) RemoveSprite(index int) {
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
}

func main() {

	s := 1.5
	w, h := int(640*s), int(480*s)

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Game Window")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	defaultScene := *testScene2()
	GameInstance.scenes.Push(&defaultScene)

	if err := ebiten.RunGame(GameInstance); err != nil {
		log.Fatal(err)
	}
}
