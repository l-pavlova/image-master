package app

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"image/color"
	"image/draw"

	"github.com/l-pavlova/image-master/imagemanip"
	"github.com/l-pavlova/image-master/logging"
	"github.com/l-pavlova/image-master/mongo"
	"github.com/l-pavlova/image-master/tensorflowAPI"
	"golang.org/x/image/colornames"
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

type MongoClient interface {
	GetAllImageClassifications() (ics []mongo.ImageClassification, err error)
	AddImageClassification(imagePath string, probabilities []string) error
	GetImageClassification(imagePath string) (ic mongo.ImageClassification, ok bool)
}

type ImageMaster struct {
	//	tfClient    TensorFlowClient
	logger      *logging.ImageMasterLogger
	imageList   []string
	mu          sync.Mutex
	concurrency int
	mongo       MongoClient
}

type TensorFlowClient interface {
	ClassifyImage(image image.Image) ([]tensorflowAPI.Label, []float32, [][]float32, error)
}

func NewImageMaster() *ImageMaster {
	imagemaster := &ImageMaster{
		//tfClient:    nil,
		imageList:   make([]string, 0, 5),
		logger:      logging.NewImageMasterLogger(),
		concurrency: DEFAULT_CONCURRENCY,
		//mongo:       &mongo.NewMongo(),
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
func (im *ImageMaster) GrayScale() error {

	im.execute(func(img image.Image, imagePath string) {
		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y
		generated := imagemanip.GenerateNew(width, height)
		target, ok := generated.(Changeable)
		if !ok {
			im.logger.Log("error", "Error occurred during image conversion, cannot filter this image.")
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				target.Set(x, y, color.Gray16Model.Convert(img.At(x, y)))
			}
		}
		_, err := imagemanip.SaveTo(path.Join(BOUND_PATH, "output"), filepath.Base(imagePath), target.(image.Image))
		if err != nil {
			im.logger.Log("error", "Error occurred during img saving %w", err)
		}
	})

	return nil
}

// The smoothen function uses a gaussian blur with kernel of size 5x5 to smoothen an image
// inPath is the path from which an image is taken
// outPath is the path where the image is saved
// timesToRepeat is an integer value, signaling how many times the blur would be applied to the image a.k.a filter strength
func (im *ImageMaster) Smoothen(timesToRepeat int) error {

	im.execute(func(img image.Image, imagePath string) {
		res, err := imagemanip.ApplyGaussian(img)
		if err != nil {
			im.logger.Log("error", "Error occurred during img smoothing with gaussian filter %w", err.Error())
		}

		for i := 0; i < timesToRepeat-1; i++ {
			tempRes, err := imagemanip.ApplyGaussian(res)
			if err != nil {
				im.logger.Log("error", "Error occurred during img smoothing with gaussian filter on iteration %d %w", i, err)
			}

			res = tempRes
		}
		_, err = imagemanip.SaveTo(path.Join(BOUND_PATH, "output"), filepath.Base(imagePath), res)
		if err != nil {
			im.logger.Log("error", "Error occurred during img saving %w", err)
		}
	})

	return nil
}

// apply sharpening via morphological operations
func (im *ImageMaster) Sharpen() error {

	im.execute(func(img image.Image, imagePath string) {
		res := imagemanip.MorphGradient(img)

		_, err := imagemanip.SaveTo(path.Join(BOUND_PATH, "output"), filepath.Base(imagePath), res)
		if err != nil {
			im.logger.Log("error", "Error occurred during img saving %w", err)
		}
	})
	return nil
}

// the folder passed to the docker image is mounted to the //images folder inside the container, so we perform our operations inside there
func (im *ImageMaster) Find(object string) error {

	//if cached return from cache
	im.execute(func(img image.Image, imagePath string) {
		//todo: optimize model to be loaded just once
		tfClient := tensorflowAPI.NewTensorFlowClient(*logging.NewImageMasterLogger())
		labels, _, _, err := tfClient.ClassifyImage(img)
		if err != nil {
			im.logger.Log("error", err.Error())
			return
		}

		//1. stash labels in db to use next time its called with this image path
		//2. foreach the labels and if even partial match to word, store image in outputs
		isMatch := false
		for _, label := range labels {
			if strings.Contains(label.Label, object) {
				isMatch = label.Probability >= 0.4 //if it's less it's probably not really the thing
			}
		}

		if isMatch {
			_, err := imagemanip.SaveTo(path.Join(BOUND_PATH, "output"), filepath.Base(imagePath), img)
			if err != nil {
				im.logger.Log("error", "Error occurred during img saving %w", err)
			}
		}
	})

	return nil
}

func pointOutMatches(img image.Image, labels []string, classes []float32, boxes [][]float32) {
	bounds := img.Bounds()

	res := image.NewRGBA(bounds)
	draw.Draw(res, bounds, img, bounds.Min, draw.Src)

	for ind, _ := range labels {
		x1 := float32(res.Bounds().Max.X) * boxes[ind][1]
		x2 := float32(res.Bounds().Max.X) * boxes[ind][3]
		y1 := float32(res.Bounds().Max.Y) * boxes[ind][0]
		y2 := float32(res.Bounds().Max.Y) * boxes[ind][2]

		imagemanip.Rect(res, int(x1), int(y1), int(x2), int(y2), 4, colornames.Map[colornames.Names[int(classes[ind])]])
		ind++
	}

}

func (im *ImageMaster) execute(operation func(image.Image, string)) error {
	//if cached return from cache
	im.scanDirectory(BOUND_PATH)
	limitChan := make(chan struct{}, im.concurrency)

	im.logger.Log("info", "images to process count: ", len(im.imageList))

	wg := &sync.WaitGroup{}
	imageCount := len(im.imageList)
	for j := 0; j < imageCount; j++ {
		wg.Add(1)
		limitChan <- struct{}{}
		go func() {
			defer func() {
				<-limitChan
			}()
			if len(im.imageList) == 0 {
				im.logger.Log("info", "no more images to process")
				return
			}

			im.mu.Lock()
			imagePath := im.imageList[0]
			im.imageList = im.imageList[1:]
			im.mu.Unlock()

			fmt.Println(imagePath)
			img, err := imagemanip.ReadFrom(imagePath)
			if err != nil {
				im.logger.Log("error", "Error occurred during img parsing %w", err)
			}

			operation(img, imagePath)
			wg.Done()
			//initialize a new client for each, but read the graph and pass it once to all?

		}()
	}

	wg.Wait()
	close(limitChan)

	return nil
}
