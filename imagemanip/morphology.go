package imagemanip

import (
	"image"
	"image/color"
)

// sharpening function applying morphological dilation and erosion in order to sharpen the image
func MorphGradient(img image.Image) image.Image {

	eroded := erode(img, 5)
	dilated := dilate(eroded, 5)
	secondDilated := dilate(dilated, 5)
	final := erode(secondDilated, 5)

	return final
}

// the dilate morphological operation
func dilate(img image.Image, structElementSize int) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			var maxColor color.Color
			for i := -structElementSize; i <= structElementSize; i++ {
				for j := -structElementSize; j <= structElementSize; j++ {
					if x+j >= bounds.Min.X && x+j < bounds.Max.X && y+i >= bounds.Min.Y && y+i < bounds.Max.Y {
						c := img.At(x+j, y+i)
						if maxColor == nil || func() bool {
							_, _, _, a := c.RGBA()
							return a == 65535
						}() {
							maxColor = c
						}
					}
				}
			}
			dst.Set(x, y, maxColor)
		}
	}

	return dst
}

// the erode morphological operation
func erode(img image.Image, structElementSize int) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			var minColor color.Color
			for i := -structElementSize; i <= structElementSize; i++ {
				for j := -structElementSize; j <= structElementSize; j++ {
					if x+j >= bounds.Min.X && x+j < bounds.Max.X && y+i >= bounds.Min.Y && y+i < bounds.Max.Y {
						c := img.At(x+j, y+i)
						if minColor == nil || func() bool {
							_, _, _, a := c.RGBA()
							return a == 0
						}() {
							minColor = c
						}
					}
				}
			}
			dst.Set(x, y, minColor)
		}
	}

	return dst
}
