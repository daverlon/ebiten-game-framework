package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func flappyScene() *Scene {
	s := &Scene{}
	s.name = "Flappy bird clone scene"

	var playerRef *Sprite

	var fallSpeed float64 // vertical velocity

	const maxFallSpeed = 2
	const fallAccel = 0.09

	const jumpForce = 6

	const jumpCooldown = 13 // ticks (60 per second)
	var jumpTimer int64

	var defaultPipeSprite Sprite

	const pipeSpawnCooldown = 140
	var pipeSpawnTimer int64

	const pipeMoveSpeed = 0.6

	s.Init = func() {
		GameInstance.sprites = nil

		blueDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-blue.png")
		if err != nil {
			log.Fatal(err)
		}
		playerSprite := Sprite{
			img:     blueDino,
			x:       -80,
			y:       0,
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
		GameInstance.AddSprite(&playerSprite)

		pipeImage, _, err := ebitenutil.NewImageFromFile("sprites/pipe2.png")
		if err != nil {
			log.Fatal(err)
		}
		defaultPipeSprite = Sprite{
			img:     pipeImage,
			x:       128,
			y:       0,
			w:       24,
			h:       96,
			centerx: 12,
			centery: 12,
			subx:    0,
			suby:    0,
			subw:    24,
			subh:    96,
		}

		fmt.Println("Initialized " + s.name)
		s.initialized = true

		GameInstance.cam.x = 0
		GameInstance.cam.y = 0

		jumpTimer = jumpCooldown
		pipeSpawnTimer = pipeSpawnCooldown
	}

	var PlayerJump func()
	var ApplyGravity func()
	var CreatePipes func()
	var HandlePipes func()

	s.Update = func() {
		if !s.initialized {
			s.Init()
		}

		CreatePipes()

		ApplyGravity()

		// apply jump
		PlayerJump()

		// apply gravity to player
		playerRef.y += fallSpeed

		HandlePipes()

	}

	PlayerJump = func() {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && jumpTimer >= jumpCooldown {
			jumpTimer = 0
			fallSpeed -= jumpForce
		}
		if jumpTimer <= jumpCooldown {
			jumpTimer++
		}
	}

	ApplyGravity = func() {
		if fallSpeed < maxFallSpeed {
			delta := math.Abs(fallSpeed - maxFallSpeed)
			fallSpeed += delta * fallAccel
		} else {
			fallSpeed = maxFallSpeed
		}
	}

	CreatePipes = func() {
		if pipeSpawnTimer >= pipeSpawnCooldown {
			pipeSpawnTimer = 0

			bottomPipe := defaultPipeSprite
			GameInstance.AddSprite(&bottomPipe)

			topPipe := defaultPipeSprite
			topPipe.y += 50
			topPipe.op.GeoM.Scale(1, -1)
			GameInstance.AddSprite(&topPipe)

		} else {
			pipeSpawnTimer++
		}
	}

	HandlePipes = func() {
		// remove pipes when off screen
		// and add pipes with timer
		deadzonex := -128.0 - 24
		for i := 1; i < len(GameInstance.sprites); i++ {
			if GameInstance.sprites[i].x < deadzonex {
				GameInstance.RemoveSprite(i)
				continue
			}
			GameInstance.sprites[i].x -= pipeMoveSpeed
		}
	}

	return s
}
