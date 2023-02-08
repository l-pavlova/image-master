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

// The GrayScale function uses a standart method from the image library  to convert an image to GrayScale and saves it to the outPath directory
// For a realistic RGB -> grayscale conversion, the following weights have to be used: Y = 0.299 * R +  0.587 * G + 0.114 * B
// inPath is the path from which an image is taken
// outPath is the path where the image is saved
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

// The smoothen function uses a gaussian blur with kernel of size 5x5 to smoothen an image
// inPath is the path from which an image is taken
// outPath is the path where the image is saved
// timesToRepeat is an integer value, signaling how many times the blur would be applied to the image a.k.a filter strength
func (i *ImageMaster) Smoothen(inPath, outPath string, timesToRepeat int) error {
	img, err := imagemanip.ReadFrom(inPath)
	if err != nil {
		return fmt.Errorf("Error occurred during img parsing %w", err)
	}

	res, err := imagemanip.ApplyGaussian(img)
	if err != nil {
		return fmt.Errorf("Error occurred during img smoothing with gaussian filter %w", err)
	}

	for i := 0; i < timesToRepeat-1; i++ {
		tempRes, err := imagemanip.ApplyGaussian(res)
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
