package smartcircle

import (
	"image"
	"image/color"
)

// Circle for crop circle.
type circle struct {
	p image.Point
	r int
}

// ColorModel return color.AlphaModel
func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

// Bounds return image.Rect
func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

// At retrun color.Alpaha
func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
