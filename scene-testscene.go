package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func testScene() *Scene {
	s := &Scene{}
	s.name = "Movement and camera test scene"

	var playerRef *Sprite
	var otherRef *Sprite

	var playerHealth int

	s.Init = func() {
		GameInstance.sprites = nil

		playerHealth = 100

		blueDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-blue.png")
		if err != nil {
			log.Fatal(err)
		}

		redDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-red.png")

		playerSprite := Sprite{
			img:     blueDino,
			x:       32,
			y:       32,
			w:       24,
			h:       24,
			centerx: 12,
			centery: 12,
			subx:    0,
			suby:    0,
			subw:    24,
			subh:    24,
		}
		playerRef = &playerSprite
		GameInstance.sprites = append(GameInstance.sprites, &playerSprite)

		otherSprite := Sprite{
			img:  redDino,
			x:    64,
			y:    64,
			w:    24,
			h:    24,
			subx: 0,
			suby: 0,
			subw: 24,
			subh: 24,
		}
		otherRef = &otherSprite
		GameInstance.sprites = append(GameInstance.sprites, &otherSprite)

		fmt.Println("Initialized " + s.name)
		s.initialized = true
	}

	s.Update = func() {
		if !s.initialized {
			s.Init()
		}

		updateCamera(*playerRef)
		movePlayer(playerRef)
		updateEnemy(*playerRef, *&otherRef)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			otherRef.x, otherRef.y = GameInstance.cursor.worldx, GameInstance.cursor.worldy
		}

		if playerHealth <= 0 {
			GameInstance.scenes.Pop()
		}

		// detect collision
		if RectCollision(
			playerRef.x,
			playerRef.y,
			playerRef.w,
			playerRef.h,
			otherRef.x,
			otherRef.y,
			otherRef.w,
			otherRef.h) {
			playerHealth -= 1
		}

	}

	//s.Draw = func() *ebiten.Image {
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Health: %d", playerHealth), windowcenterx-33, 5)
	//}

	return s
}

func movePlayer(player *Sprite) {

	moveSpeed := float64(2)

	var inputX float64
	var inputY float64

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		inputX = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		inputX = 1
	} else {
		inputX = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		inputY = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		inputY = 1
	} else {
		inputY = 0
	}

	player.x -= inputX * moveSpeed
	player.y += inputY * moveSpeed

}

func updateEnemy(player Sprite, other *Sprite) {
	ms := 0.2
	if other.x > player.x {
		other.x -= ms
	}
	if other.x <= player.x {
		other.x += ms
	}
	if other.y > player.y {
		other.y -= ms
	}
	if other.y <= player.y {
		other.y += ms
	}
}

func updateCamera(p Sprite) {
	_, mY := ebiten.Wheel()
	GameInstance.cam.zoom += mY
	minzoom := 0.5
	if GameInstance.cam.zoom < minzoom {
		GameInstance.cam.zoom = minzoom
	}
	GameInstance.cam.SlowlyMove(p.x+p.centerx, p.y+p.centery, 0.4)
}
