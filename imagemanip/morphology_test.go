package imagemanip

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

var testImage image.Image

func init() {
	testImage = image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(testImage.(*image.RGBA), testImage.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)
	draw.Draw(testImage.(*image.RGBA), image.Rect(25, 25, 75, 75), &image.Uniform{green}, image.Point{}, draw.Src)
	draw.Draw(testImage.(*image.RGBA), image.Rect(40, 40, 60, 60), &image.Uniform{blue}, image.Point{}, draw.Src)
}

func TestMorphGradient(t *testing.T) {
	result := MorphGradient(testImage)

	if result == nil {
		t.Errorf("MorphGradient returned nil")
	}

	if result.Bounds() != testImage.Bounds() {
		t.Errorf("MorphGradient did not preserve image dimensions")
	}
}

func TestDilate(t *testing.T) {
	result := dilate(testImage, 3)

	if result == nil {
		t.Errorf("dilate returned nil")
	}

	if result.Bounds() != testImage.Bounds() {
		t.Errorf("dilate did not preserve image dimensions")
	}
}

func TestErode(t *testing.T) {
	result := erode(testImage, 3)

	if result == nil {
		t.Errorf("erode returned nil")
	}

	if result.Bounds() != testImage.Bounds() {
		t.Errorf("erode did not preserve image dimensions")
	}
}
