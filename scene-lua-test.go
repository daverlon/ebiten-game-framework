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

	var l_update func()

	s.Init = func() {

		if err := l.DoFile("test.lua"); err != nil {
			fmt.Println("Error reading lua file:")
			fmt.Println(err)

			l.SetGlobal("update", l_update)
		}

	}

	s.Draw = func(screen *ebiten.Image) {
		ebitenutil.DebugPrintAt(screen, "Press escape to go back", viewportw/2, viewporth/2)
	}

	s.Update = func() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			GameInstance.scenes.Pop()
		}

		if l_update != nil {
			l_update()
		}
	}

	return s
}
