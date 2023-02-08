package imagemanip

import (
	"image"
	"image/color"
	"math"
)

func ApplyGaussian(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	kernelSize := 2
	kernel := generateGaussianKernel(1.0, kernelSize)
	target := image.NewRGBA(img.Bounds())

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a float64
			for i := 0; i < len(kernel); i++ {
				for j := 0; j < len(kernel); j++ {
					k := y + j - kernelSize
					if k >= 0 && k < height {
						r1, g1, b1, a1 := img.At(x, k).RGBA()
						//do this to make r g b values from uint8 to floats, then bring them back to uint8
						r += float64(r1>>8) * kernel[i][j]
						g += float64(g1>>8) * kernel[i][j]
						b += float64(b1>>8) * kernel[i][j]
						a += float64(a1>>8) * kernel[i][j]
					}
				}
			}
			target.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return target, nil
}

// a function that generates a kernel with normal distribution for gaussian blur
func generateGaussianKernel(sigma float64, kernelSize int) [][]float64 {
	s := 2.0 * sigma * sigma
	sum := 0.0
	kernel := make([][]float64, 2*kernelSize+1)
	for i := range kernel {
		kernel[i] = make([]float64, 2*kernelSize+1)
	}
	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			root := math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y))
			kernel[x+2][y+2] = math.Exp(-(root*root)/s) / (math.Pi * s)
			sum += kernel[x+2][y+2]
		}
	}

	//normalize it
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}
