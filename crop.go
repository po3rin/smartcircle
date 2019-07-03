package smartcircle

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/pkg/errors"
)

// Cropper for crop circle.
type Cropper interface {
	CropCircle() (*image.RGBA, error)
	setSrc(src image.Image) error
}

// SubImager for smartcrop
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type cropper struct {
	src    image.Image
	radius int
}

// Params is parameters for NewDrawer functio
type Params struct {
	Src image.Image
}

// NewCropper init cropper from Params
func NewCropper(params Params) (Cropper, error) {
	d := &cropper{}
	err := d.setSrc(params.Src)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (c *cropper) setSrc(src image.Image) error {
	b := src.Bounds()
	srcWidth := b.Max.X
	srcHeight := b.Max.Y

	var radius int
	if srcWidth <= srcHeight {
		radius = srcWidth / 2
	} else {
		radius = srcHeight / 2
	}

	c.src = src
	c.radius = radius

	return nil
}

// CropCircle crop a circle image out of image.
func (c *cropper) CropCircle() (*image.RGBA, error) {
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, err := analyzer.FindBestCrop(c.src, 2*c.radius, 2*c.radius)
	if err != nil {
		return nil, errors.Wrap(err, "smartcircle: failed to smart crop")
	}

	// prepare src for draw
	src := c.src.(SubImager).SubImage(topCrop)

	// prepare src for mask
	mask := &circle{p: image.Point{c.radius, c.radius}, r: c.radius}

	// prepare dst for draw
	rect := image.Rect(0, 0, c.radius*2, c.radius*2)
	dst := image.NewRGBA(rect)
	fillRect(dst, color.RGBA{0, 0, 0, 0})

	draw.DrawMask(dst, dst.Bounds(), src, src.Bounds().Min, mask, image.ZP, draw.Over)
	return dst, nil
}

func fillRect(img *image.RGBA, col color.Color) {
	rect := img.Rect
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			img.Set(v, h, col)
		}
	}
}
