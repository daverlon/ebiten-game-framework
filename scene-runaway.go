package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func testScene() *Scene {
	s := &Scene{name: "Movement and camera test"}

	// vars
	var playerRef *Sprite
	var otherRef *Sprite

	var playerHealth int

	// functions

	s.Init = func() {

		playerHealth = 100

		blueDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-blue.png")
		if err != nil {
			log.Fatal(err)
		}

		redDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-red.png")
		if err != nil {
			log.Fatal(err)
		}

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
	}

	s.Draw = func(screen *ebiten.Image) {
		for _, spr := range GameInstance.sprites {
			spr.Draw(screen)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Health: %d", playerHealth), windowcenterx-33, 5)
	}

	var movePlayer func()
	var updateEnemy func()
	var updateCamera func()

	s.Update = func() {

		updateCamera()
		movePlayer()
		updateEnemy()

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

	movePlayer = func() {

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

		playerRef.x -= inputX * moveSpeed
		playerRef.y += inputY * moveSpeed

	}

	updateEnemy = func() {

		ms := 0.02

		deltax := playerRef.x - otherRef.x
		deltay := playerRef.y - otherRef.y
		//fmt.Println(deltax, deltay)

		otherRef.x += deltax * ms
		otherRef.y += deltay * ms
	}

	updateCamera = func() {
		GameInstance.cam.SlowlyMove(playerRef.x+playerRef.centerx, playerRef.y+playerRef.centery, 0.4)
	}

	return s
}
