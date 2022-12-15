package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func testScene2() *Scene {
	s := &Scene{}
	s.name = "Death menu test scene"

	s.Init = func() {
		GameInstance.sprites = nil
		fmt.Println("Initialized " + s.name)
		s.initialized = true
	}

	s.Update = func() {
		if !s.initialized {
			s.Init()
		}

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

	/*s.Draw = func(screen *ebiten.Image) {
		ebitenutil.DebugPrintAt(screen, "Press Q to quit.", windowcenterx, windowcentery)
		ebitenutil.DebugPrintAt(screen, "Press 1 to play flappy game", windowcenterx, windowcentery+20)
		ebitenutil.DebugPrintAt(screen, "Press 2 to play run away game", windowcenterx, windowcentery+40)
	}*/

	return s
}
