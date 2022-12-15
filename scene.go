package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene struct {
	name string

	Init   func()
	Update func()
	Draw   func(screen *ebiten.Image)

	//g      *Game
}

// ------- //

type SceneStack []*Scene

func (s *SceneStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *SceneStack) Push(ss *Scene) {
	GameInstance.sprites = nil
	ss.Init()
	*s = append(*s, ss)
	fmt.Println("Pushed Scene to SceneStack. Scene count:", len(*s))
}

func (s *SceneStack) Pop() (*Scene, bool) {
	if s.IsEmpty() {
		fmt.Println("Warning: No scenes to pop.")
		return nil, false
	} else {
		GameInstance.sprites = nil
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		fmt.Println("Popped Scene from SceneStack. Scene count:", len(*s))
		return element, true
	}
}

// returns the scene on top of the stack
func (s *SceneStack) Active() *Scene {
	if !s.IsEmpty() {
		return (*s)[len(*s)-1]
	} else {
		return nil
	}
}
