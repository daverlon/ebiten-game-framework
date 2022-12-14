package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
