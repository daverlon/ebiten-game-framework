package main

type Camera struct {
	x    float64
	y    float64
	zoom float64 // for debugging purposes
}

func (cam *Camera) ScreenToWorld(x float64, y float64) (float64, float64) {
	retx, rety := x, y

	retx += cam.x
	rety += cam.y

	retx -= windowcenterx
	rety -= windowcentery

	retx /= cam.zoom
	rety /= cam.zoom

	return retx, rety
}

func (c *Camera) SlowlyMove(x float64, y float64, speed float64) {
	targetx, targety := x*c.zoom, y*c.zoom
	deltax, deltay := targetx-c.x, targety-c.y
	c.x += deltax * speed
	c.y += deltay * speed
}
