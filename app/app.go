package app

import (
	"fmt"
	"image"
	"image/color"

	"github.com/l-pavlova/image-master/imagemanip"
)

// type changeable used for assertion on parsed images
type Changeable interface {
	Set(x, y int, c color.Color)
}

type ImageMaster struct {
}

func NewImageMaster() *ImageMaster {
	return &ImageMaster{}
}

// For a realistic RGB -> grayscale conversion, the following weights have to be used: Y = 0.299 * R +  0.587 * G + 0.114 * B
func (i *ImageMaster) GrayScale(inPath, outPath string) error {

	img, err := imagemanip.ReadFrom(inPath)
	if err != nil {
		return fmt.Errorf("Error occurred during img parsing %w", err)
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	generated := imagemanip.GenerateNew(width, height)
	target, ok := generated.(Changeable)
	if !ok {
		return fmt.Errorf("%s", "Error occurred during image conversion, cannot filter this image.")
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			target.Set(x, y, color.Gray16Model.Convert(img.At(x, y)))
		}
	}
	res, err := imagemanip.SaveTo(outPath, "test.jpg", target.(image.Image))
	if err != nil {
		return fmt.Errorf("Error occurred during img saving %w", err)
	}

	fmt.Println(res)
	return nil
}

func (i *ImageMaster) Smoothen(inPath, outPath string, timesToRepeat int) error {
	img, err := imagemanip.ReadFrom(inPath)
	if err != nil {
		return fmt.Errorf("Error occurred during img parsing %w", err)
	}

	res, err := imagemanip.Gaussian(img)
	if err != nil {
		return fmt.Errorf("Error occurred during img smoothing with gaussian filter %w", err)
	}

	for i := 0; i < timesToRepeat; i++ {
		tempRes, err := imagemanip.Gaussian(res)
		if err != nil {
			return fmt.Errorf("Error occurred during img smoothing with gaussian filter on iteration %d %w", i, err)
		}

		res = tempRes
	}

	_, err = imagemanip.SaveTo(outPath, "testGaussian.jpg", res)
	if err != nil {
		return fmt.Errorf("Error occurred during img saving %w", err)
	}
	return nil
}
