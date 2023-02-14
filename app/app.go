package app

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"

	"image/color"

	"github.com/l-pavlova/image-master/imagemanip"
	"github.com/l-pavlova/image-master/logging"
	"github.com/l-pavlova/image-master/tensorflowAPI"
)

const (
	BOUND_PATH          string = "./images/"
	OUT_PATH            string = "./output/" //todo: customize
	DEFAULT_CONCURRENCY int    = 10
)

// type changeable used for assertion on parsed images
type Changeable interface {
	Set(x, y int, c color.Color)
}

type ImageMaster struct {
	//	tfClient    TensorFlowClient
	logger      *logging.ImageMasterLogger
	imageList   []string
	mu          sync.Mutex
	concurrency int
}

type TensorFlowClient interface {
	ClassifyImage(image image.Image) error
}

func NewImageMaster() *ImageMaster {
	imagemaster := &ImageMaster{
		//tfClient:    nil,
		imageList:   make([]string, 0, 5),
		logger:      logging.NewImageMasterLogger(),
		concurrency: DEFAULT_CONCURRENCY,
	}

	//imagemaster.tfClient = &*tensorflowAPI.NewTensorFlowClient(*imagemaster.logger)
	return imagemaster
}

// scan directory iterates over the directory of passed images, and saves all of their paths in imageList to later convert
func (i *ImageMaster) scanDirectory(directory string) {
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				i.imageList = append(i.imageList, path)
				i.logger.Log("info", "adding image to process from path: ", path, info.Size())
			}

			return nil
		})
	if err != nil {
		i.logger.Log("error", err.Error())
	}
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

// apply sharpening via morphological operations
func (i *ImageMaster) Sharpen(inPath, outPath string, timesToRepeat int) error {

	img, err := imagemanip.ReadFrom(inPath)
	if err != nil {
		return fmt.Errorf("Error occurred during img parsing %w", err)
	}

	res := imagemanip.MorphGradient(img)
	// for i := 0; i < timesToRepeat-1; i++ {
	// 	tempRes := imagemanip.MorphGradient(res)
	// 	if err != nil {
	// 		return fmt.Errorf("Error occurred during img smoothing with gaussian filter on iteration %d %w", i, err)
	// 	}
	// 	res = tempRes
	// }

	_, err = imagemanip.SaveTo(outPath, "testSharpen.jpg", res)
	if err != nil {
		return fmt.Errorf("Error occurred during img saving %w", err)
	}

	return nil
}

// the folder passed to the docker image is mounted to the //images folder inside the container, so we perform our operations inside there
func (i *ImageMaster) Find(object string) error {

	//if cached return from cache
	i.execute(func(img image.Image) {
		//initialize a new client for each, but read the graph and pass it once to all?
		tfClient := tensorflowAPI.NewTensorFlowClient(*logging.NewImageMasterLogger())
		err := tfClient.ClassifyImage(img)
		if err != nil {
			fmt.Println("Error ocurred %w", err)
		}
	})

	return nil
}

func (i *ImageMaster) execute(operation func(image.Image)) error {
	//if cached return from cache
	i.scanDirectory(BOUND_PATH)
	limitChan := make(chan struct{}, i.concurrency)

	var wg sync.WaitGroup
	for j := 0; j < len(i.imageList); j++ {
		wg.Add(1)
		limitChan <- struct{}{}
		go func() {
			defer func() {
				<-limitChan
			}()

			i.mu.Lock()
			if len(i.imageList) == 0 {
				i.logger.Log("info", "no more images to process")
				wg.Done()
				i.mu.Unlock()
				close(limitChan)
				return
			}

			imagePath := i.imageList[0]
			i.imageList = i.imageList[1:]
			i.mu.Unlock()

			img, err := imagemanip.ReadFrom(imagePath)
			if err != nil {
				i.logger.Log("error", "Error occurred during img parsing %w", err)
			}

			operation(img)
			//initialize a new client for each, but read the graph and pass it once to all?
			wg.Done()
		}()
	}

	return nil
}
