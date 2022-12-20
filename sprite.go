package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	img     *ebiten.Image
	x       float64 // world pos x
	y       float64 // world pos y
	w       float64 // image width (collisions?)
	h       float64 // image height (collisions?)
	centerx float64 // center of image x offset (w/2)
	centery float64 // center of image y offset (y/2)

	// sub image 'frame'
	subx int // sub image x pos
	suby int // sub image y pos
	subw int // sub image width
	subh int // sub image height

	op ebiten.DrawImageOptions
}

func (g *Game) AddSprite(s *Sprite) {
	g.sprites = append(g.sprites, s)
}

func (g *Game) RemoveSprite(index int) {
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
}

// draw the sprite to a screen
func (s *Sprite) Draw(screen *ebiten.Image) {
	op := s.op
	op.GeoM.Scale(GameInstance.cam.zoom, GameInstance.cam.zoom)

	scaledx, scaledy := s.x*GameInstance.cam.zoom, s.y*GameInstance.cam.zoom

	// adjust draw positions based on camera positions
	tx := -GameInstance.cam.x + windowcenterx + scaledx
	ty := -GameInstance.cam.y + windowcentery + scaledy
	op.GeoM.Translate(tx, ty)

	// render image
	img := (s.img).SubImage(image.Rect(s.subx, s.suby, s.subw, s.subh)).(*ebiten.Image)
	//screen.DrawImage(img, &op)
	screen.DrawImage(img, &op)

	if showdebuginfo {
		// debug lines
		DrawOutlineRect(screen, tx, ty, s.w*GameInstance.cam.zoom, s.h*GameInstance.cam.zoom)
		ebitenutil.DrawRect(screen, tx-3, ty-3, 6, 6, clrRed)
		ebitenutil.DrawRect(screen, windowcenterx-3, windowcentery-3, 6, 6, clrYellow)

		// debug position text
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%1.0f,%1.0f)", s.x, s.y), int(tx), int(ty))
	}
}
