package main

func RectCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	// credit: ChatGPT

	// Check if any of the sides of the first rectangle are outside the second rectangle
	if x1+w1 < x2 || x1 > x2+w2 ||
		y1+h1 < y2 || y1 > y2+h2 {
		return false
	}

	// If none of the sides of the first rectangle are outside the second rectangle,
	// then they must be colliding
	return true
}
