package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func testScene2() *Scene {
	s := &Scene{name: "Main menu"}

	s.Init = func() {
		fmt.Println("Initialized " + s.name)
	}

	s.Update = func() {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			GameInstance.scenes.Push(flappyScene())
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			GameInstance.scenes.Push(testScene())
		}

		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			GameInstance.scenes.Push(luaTestScene())
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			GameInstance.scenes.Pop()
		}
	}

	s.Draw = func(screen *ebiten.Image) {
		posx := windowcenterx - 70
		ebitenutil.DebugPrintAt(screen, "Press Q to quit.", posx, windowcentery)
		ebitenutil.DebugPrintAt(screen, "Press 1 to play flappy game", posx, windowcentery+20)
		ebitenutil.DebugPrintAt(screen, "Press 2 to play run away game", posx, windowcentery+40)
		ebitenutil.DebugPrintAt(screen, "Press 3 to test lua scene", posx, windowcentery+60)
	}

	return s
}
