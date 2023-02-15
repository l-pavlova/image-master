package imagemanip

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestApplyGaussian(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			c := color.RGBA{uint8(x * 25), uint8(y * 25), 0, 255}
			testImg.Set(x, y, c)
		}
	}

	result, err := ApplyGaussian(testImg)
	if err != nil {
		t.Errorf("ApplyGaussian returned error: %v", err)
	}

	if result.Bounds() != testImg.Bounds() {
		t.Errorf("Result has different dimensions than input image")
	}

	if reflect.DeepEqual(result, testImg) {
		t.Errorf("Result is identical to input image")
	}
}
