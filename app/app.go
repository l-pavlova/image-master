package app

import (
	"fmt"
	"image"
	"image/color"

	"github.com/l-pavlova/image-master/imageparse"
)

//type changeable used for assertion on parsed images
type Changeable interface {
	Set(x, y int, c color.Color)
}

type ImageMaster struct {
}

func NewImageMaster() *ImageMaster {
	return &ImageMaster{}
}

//For a realistic RGB -> grayscale conversion, the following weights have to be used: Y = 0.299 * R +  0.587 * G + 0.114 * B
func (i *ImageMaster) GrayScale(inPath, outPath string) error {

	img, err := imageparse.ReadFrom(inPath)
	if err != nil {
		return fmt.Errorf("Error occurred during img parsing %w", err)
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	generated := imageparse.GenerateNew(width, height)
	target, ok := generated.(Changeable)
	if !ok {
		return fmt.Errorf("%s", "Error occurred during image conversion, cannot filter this image.")
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			/*oldPixel := img.At(x, y)

			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}*/
			target.Set(x, y, color.Gray16Model.Convert(img.At(x, y)))
		}
	}
	res, err := imageparse.SaveTo(outPath, "test.jpg", target.(image.Image))
	if err != nil {
		return fmt.Errorf("Error occurred during img saving %w", err)
	}

	fmt.Println(res)
	return nil
}
