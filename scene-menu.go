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

		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			GameInstance.scenes.Pop()
		}
	}

	s.Draw = func(screen *ebiten.Image) {
		posx := windowcenterx - 80
		posy := windowcentery - 25
		ebitenutil.DebugPrintAt(screen, "Press Q to quit.", posx, posy)
		ebitenutil.DebugPrintAt(screen, "Press 1 to play flappy game", posx, posy+20)
		ebitenutil.DebugPrintAt(screen, "Press 2 to play run away game", posx, posy+40)
	}

	return s
}
