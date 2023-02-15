package imagemanip

import (
	"image"
	"image/color"
	"testing"
)

func TestHLine(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	col := color.RGBA{255, 0, 0, 255}
	HLine(img, 2, 5, 7, col)
	for x := 2; x <= 7; x++ {
		if img.At(x, 5) != col {
			t.Errorf("HLine failed to draw a red line at x=%d, y=5", x)
		}
	}
}

func TestVLine(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	col := color.RGBA{0, 255, 0, 255}
	VLine(img, 5, 2, 7, col)
	for y := 2; y <= 7; y++ {
		if img.At(5, y) != col {
			t.Errorf("VLine failed to draw a green line at x=5, y=%d", y)
		}
	}
}

func TestRect(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	col := color.RGBA{0, 0, 255, 255}
	start := 2
	end := 7
	Rect(img, start, start, end, end, 1, col)
	for x := start; x <= end; x++ {
		if img.At(x, 2) != col {
			t.Errorf("Rect failed to draw a blue line at x=%d, y=2", x)
		}
		if img.At(x, 7) != col {
			t.Errorf("Rect failed to draw a blue line at x=%d, y=7", x)
		}
	}
	for y := start; y <= end; y++ {
		if img.At(2, y) != col {
			t.Errorf("Rect failed to draw a blue line at x=2, y=%d", y)
		}
		if img.At(7, y) != col {
			t.Errorf("Rect failed to draw a blue line at x=7, y=%d", y)
		}
	}
}
