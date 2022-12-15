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
	s := &Scene{
		name: "Flappy bird clone scene",
	}

	// --- scene vars --- //

	// reset camera for this scene
	GameInstance.cam.x, GameInstance.cam.y = 0, 0
	GameInstance.cam.zoom = 1

	var paused = false

	// player reference
	var playerRef *Sprite
	const playerX = -80

	// player physics vars
	var fallSpeed float64    // player's current vertical velocity
	const maxFallSpeed = 2.7 // fallspeed will accelerate until this speed
	const fallAccel = 0.09   // acceleration multiplier
	const jumpForce = 6.9    // jump force to be added to player velocity

	// jump vars
	const jumpCooldown = 13      // ticks (60 per second)
	var jumpTimer = jumpCooldown // current timer for the jump cooldown
	// make the jumptimer = cooldown to start,
	// so the player can jump instantly

	// pipe vars
	var defaultPipeSprite Sprite // default pipe to have its' value copied into new pipes
	const pipeSpawnCooldown = 130
	var pipeSpawnTimer = pipeSpawnCooldown // make this the cooldown so pipe spawns instantly
	const pipeMoveSpeed = 0.8
	const pipeGap = 130
	var pipeDeadZoneX float64 // defined later since it relies on defaultPipeSprite.w

	// initialize the scene
	//s.Init
	s.Init = func() {

		// background sprite
		pyramids, _, err := ebitenutil.NewImageFromFile("sprites/pyramid.png")
		if err != nil {
			log.Fatal(err)
		}
		bgSprite := Sprite{
			img:     pyramids,
			x:       0 - viewportw/2,
			y:       0 - viewporth/2,
			w:       viewportw,
			h:       viewporth,
			centerx: 0,
			centery: 0,
			subx:    0,
			suby:    0,
			subw:    900,
			subh:    725,
		}
		bgSprite.op.GeoM.Scale(0.5, 0.5)
		GameInstance.AddSprite(&bgSprite)

		// player sprite
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

		// pipe sprite
		pipeImage, _, err := ebitenutil.NewImageFromFile("sprites/pipe2.png")
		if err != nil {
			log.Fatal(err)
		}
		defaultPipeSprite = Sprite{
			img:     pipeImage,
			x:       viewportw / 2,
			y:       0,
			w:       24,
			h:       160,
			centerx: 12,
			centery: 12,
			subx:    0, // shift to get flipped sprite
			suby:    0,
			subw:    24,
			subh:    160,
		}

		pipeDeadZoneX = -viewportw/2.0 - defaultPipeSprite.w
		fmt.Println(pipeDeadZoneX)

		fmt.Println("Initialized " + s.name)
	}

	// (optional) draw function
	s.Draw = func(screen *ebiten.Image) {
		// loop through the game sprites and draw each one
		for _, s := range GameInstance.sprites {
			s.Draw(screen)
		}

		// if the scene is paused, display the pause message
		if paused {
			ebitenutil.DebugPrintAt(screen, "Paused", windowcenterx, windowcentery)
		}
	}

	// (optional) declaration of nested functions (which are defined below)
	var PlayerJump func()
	var ApplyGravity func()
	var CreatePipes func()
	var MovePipes func()
	var ClearPipes func()
	var HandleCollision func()

	// update every tick
	s.Update = func() {
		// check if the scene is already initialized

		// pause functinality
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			paused = !paused
		}
		if paused {
			return // if the scene is paused, skip everything below
		}

		// clear pipes which are off screen
		ClearPipes()

		// create new pipes with a timer
		CreatePipes()

		// move pipes to the left
		MovePipes()

		// apply gravity to the player
		ApplyGravity()

		// allow player to jump with timer
		PlayerJump()

		// apply gravity to player
		playerRef.y += fallSpeed

		// lose when the player collides with a pipe
		HandleCollision()

	}

	// bulk of the game logic below, separated into functions

	// give the player the ability to jump
	PlayerJump = func() {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && jumpTimer >= jumpCooldown {
			jumpTimer = 0
			fallSpeed -= jumpForce
		}
		if jumpTimer <= jumpCooldown {
			jumpTimer++
		}
	}

	// apply gravity to the player
	ApplyGravity = func() {
		if fallSpeed < maxFallSpeed {
			delta := math.Abs(fallSpeed - maxFallSpeed)
			fallSpeed += delta * fallAccel
		} else {
			fallSpeed = maxFallSpeed
		}
	}

	// create new pipes with the spawn timer
	CreatePipes = func() {

		if pipeSpawnTimer >= pipeSpawnCooldown {

			// randomly determine the y position of the new pipe
			thisy := -30 + float64(rand.Intn(30*2.0))

			// reset the timer
			pipeSpawnTimer = 0

			// copy the defaultpipesprite into a new 'bottom pip' sprite
			bottomPipe := defaultPipeSprite
			bottomPipe.y = 0 + pipeGap/2
			bottomPipe.y += thisy
			GameInstance.AddSprite(&bottomPipe)

			// do the same for the top sprite, however with the subx offset for the next frame
			topPipe := defaultPipeSprite
			topPipe.subx = 48 // todo: figure out why this is works as 48 and not 24
			topPipe.y = -96 - pipeGap/2
			topPipe.y += thisy
			GameInstance.AddSprite(&topPipe)

		} else {
			// increment the timer
			pipeSpawnTimer++
		}
	}

	MovePipes = func() {
		// loop through the pipes and move them toward the player
		// the loop starts at index=2,
		// because the sprites before index 2 are the background & player
		for i := 2; i < len(GameInstance.sprites); i++ {
			GameInstance.sprites[i].x -= pipeMoveSpeed
		}
	}

	ClearPipes = func() {
		// detect if any pipes are beyond the dead zone,
		// if they are, remove them

		// avoid looping through every pipe
		/*a := len(GameInstance.sprites)
		if a >= 4 {
			a = 4
		}*/
		// loop from 2 (easy way to skip non-pipes)
		//for i := 2; i < a; i++ {
		for i := 2; i < len(GameInstance.sprites); i++ {
			if GameInstance.sprites[i].x < pipeDeadZoneX {
				GameInstance.RemoveSprite(i)
			}
		}
	}

	HandleCollision = func() {

		// handle a collision between the player and any pipe

		// note: this may be optimized slightly by
		// only checking if the player is colliding with
		// the only the possible pipes.
		// i.e the player can only possibly collide with
		// the pipes which are aligned

		pipeCollision := func() bool {
			g := 6.0 // generosity
			for i := 2; i < len(GameInstance.sprites); i++ {
				cur := GameInstance.sprites[i]
				// if there is a collision, return true
				if RectCollision(cur.x, cur.y, cur.w, cur.h,
					playerRef.x+g, playerRef.y+g, playerRef.w-(g*2), playerRef.h-(g*2)) {
					return true
				}
			}
			return false
		}

		// check if the player has hit the top or bottom of the screen
		outOfBounds := func() bool {
			g := 12.0 // generosity for the player
			// use viewporth/2 because the center of the screen=viewport/2
			if playerRef.y+playerRef.h > viewporth/2+g {
				return true
			}
			if playerRef.y < -viewporth/2-g {
				return true
			}
			return false
		}

		// if either collision is true, pop this scene from the scene stack
		// this results in the 'main-menu' scene taking over
		if pipeCollision() || outOfBounds() {
			GameInstance.scenes.Pop()
		}
	}

	return s
}
