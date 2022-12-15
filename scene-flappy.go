package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func flappyScene() *Scene {
	s := &Scene{}
	s.name = "Flappy bird clone scene"

	var drawTarget *ebiten.Image

	var paused = false

	var playerRef *Sprite
	const playerX = -80

	var fallSpeed float64 // vertical velocity

	const maxFallSpeed = 2.7
	const fallAccel = 0.09

	const jumpForce = 6.9

	const jumpCooldown = 13 // ticks (60 per second)
	var jumpTimer int64

	var defaultPipeSprite Sprite

	const pipeSpawnCooldown = 130
	var pipeSpawnTimer int64

	const pipeMoveSpeed = 0.8

	const pipeGap = 55

	s.Init = func() {
		GameInstance.sprites = nil

		pyramids, _, err := ebitenutil.NewImageFromFile("sprites/pyramid.png")
		if err != nil {
			log.Fatal(err)
		}
		bgSprite := Sprite{
			img:     pyramids,
			x:       -128,
			y:       -96,
			w:       128 * 2,
			h:       96 * 2,
			centerx: 0,
			centery: 0,
			subx:    0,
			suby:    0,
			subw:    900,
			subh:    725,
		}
		bgSprite.op.GeoM.Scale(0.3, 0.3)
		GameInstance.AddSprite(&bgSprite)

		blueDino, _, err := ebitenutil.NewImageFromFile("sprites/dino-blue.png")
		if err != nil {
			log.Fatal(err)
		}
		playerSprite := Sprite{
			img:     blueDino,
			x:       playerX,
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
		drawTarget = ebiten.NewImage(640, 480)
	}

	s.Draw = func() *ebiten.Image {
		drawTarget.Clear()
		for _, s := range GameInstance.sprites {
			s.Draw(drawTarget)
		}

		return drawTarget
	}

	var PlayerJump func()
	var ApplyGravity func()
	var CreatePipes func()
	var MovePipes func()
	var ClearPipes func()
	var HandleCollision func()

	s.Update = func() {
		if !s.initialized {
			s.Init()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			paused = !paused
		}
		if paused {
			return
		}

		// clear pipes which are off screen
		ClearPipes()

		// create new pipes with a timer
		CreatePipes()

		// apply gravity to the player
		ApplyGravity()

		// allow player to jump with timer
		PlayerJump()

		// apply gravity to player
		playerRef.y += fallSpeed

		// move pipes to the left
		MovePipes()

		// lose when the player collides with a pipe
		HandleCollision()
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

		thisy := -30 + float64(rand.Intn(30*2.0))

		if pipeSpawnTimer >= pipeSpawnCooldown {
			pipeSpawnTimer = 0

			bottomPipe := defaultPipeSprite
			bottomPipe.y = 0 + pipeGap/2
			bottomPipe.y += thisy
			GameInstance.AddSprite(&bottomPipe)

			topPipe := defaultPipeSprite
			topPipe.subx = 48 // todo: figure out why this is works as 48 and not 24
			topPipe.y = -96 - pipeGap/2
			topPipe.y += thisy
			GameInstance.AddSprite(&topPipe)

		} else {
			pipeSpawnTimer++
		}
	}

	MovePipes = func() {
		for i := 2; i < len(GameInstance.sprites); i++ {
			GameInstance.sprites[i].x -= pipeMoveSpeed
		}
	}

	ClearPipes = func() {
		deadzonex := -128.0 - 24
		for i := 2; i < len(GameInstance.sprites); i++ {
			if GameInstance.sprites[i].x < deadzonex {
				GameInstance.RemoveSprite(i)
			}
		}
	}

	HandleCollision = func() {

		// collision with pipe
		pipeCollision := func() bool {
			g := 6.0 // generosity
			for i := 2; i < len(GameInstance.sprites); i++ {
				cur := GameInstance.sprites[i]
				if RectCollision(cur.x, cur.y, cur.w, cur.h,
					playerRef.x+g, playerRef.y+g, playerRef.w-(g*2), playerRef.h-(g*2)) {
					return true
				}
			}
			return false
		}

		outOfBounds := func() bool {
			g := 12.0
			if playerRef.y+playerRef.h > 96+g {
				return true
			}
			if playerRef.y < -96-g {
				return true
			}
			return false
		}

		if pipeCollision() || outOfBounds() {
			GameInstance.scenes.Pop()
		}

	}

	return s
}
