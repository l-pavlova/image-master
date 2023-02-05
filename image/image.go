package imageParse

import (
	"image"
	"image/jpeg"
	"os"
)

func ReadFrom(filePath string) (image.Image, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	image, _, err := image.Decode(fd)
	return image, err
}

func SaveTo(filePath string, img image.Image) (bool, error) {

	fd, err := os.Create("test")
	if err != nil {
		return false, err
	}
	defer fd.Close()
	if err = jpeg.Encode(fd, img, nil); err != nil {
		return false, err
	}

	return true, nil
}
