package imageparse

import (
	"image"
	"image/jpeg"
	"os"
	"strings"
)

//util func that returns the contents of an image file
func ReadFrom(filePath string) (image.Image, error) {

	//todo: add handling for non image files
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	image, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return image, err
}

func SaveTo(filePath, fileName string, img image.Image) (bool, error) {

	fd, err := os.Create(strings.Join([]string{filePath, fileName}, "/"))
	if err != nil {
		return false, err
	}
	defer fd.Close()

	if err = jpeg.Encode(fd, img, nil); err != nil {
		return false, err
	}

	return true, nil
}

func GenerateNew(width, height int) image.Image {

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	if img != nil {
		return img
	}
	return nil
}
