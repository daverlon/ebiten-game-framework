package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func luaTestScene() *Scene {
	s := &Scene{name: "Lua test scene"}

	l := lua.NewState()
	defer l.Close()

	s.Init = func() {

		if err := l.DoFile("scenes/init.lua"); err != nil {
			fmt.Println("Error reading lua file:")
			fmt.Println(err)
		}

	}

	s.Draw = func(screen *ebiten.Image) {
		ebitenutil.DebugPrintAt(screen, "Press escape to go back", viewportw/2, viewporth/2)
	}

	s.Update = func() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			GameInstance.scenes.Pop()
		}

		if err := l.DoFile("scenes/update.lua"); err != nil {
			fmt.Println("Error reading lua file:")
			fmt.Println(err)

		}

	}

	return s
}
