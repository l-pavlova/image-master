package app

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/l-pavlova/image-master/logging"
	mock_app "github.com/l-pavlova/image-master/mocks"
)

var (
	tmpDir,
	img1Path,
	img2Path,
	img3Path string
)

func init() {
	var err error
	tmpDir, err = ioutil.TempDir("", "test")
	if err != nil {
		return
	}

	defer os.RemoveAll(tmpDir)
	img1Path = filepath.Join(tmpDir, "test1.jpg")
	img2Path = filepath.Join(tmpDir, "test2.jpg")
	img3Path = filepath.Join(tmpDir, "test3.png")
}

func TestImageMaster_Find(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockLogger := logging.NewImageMasterLogger()
	mockMongo := mock_app.NewMockMongoClient(ctrl)

	if err := createTestImageFile(img1Path); err != nil {
		t.Fatal(err)
	}
	if err := createTestImageFile(img2Path); err != nil {
		t.Fatal(err)
	}
	if err := createTestImageFile(img3Path); err != nil {
		t.Fatal(err)
	}

	im := &ImageMaster{
		logger: mockLogger,
		mongo:  mockMongo,
	}

	object := "person"
	err := im.Find(object)
	if err != nil {
		t.Fatalf("Find returned an error: %v", err)
	}

	expectedImages := []string{img1Path, img2Path}
	foundDir := filepath.Join(tmpDir, "found")
	foundImages, err := ioutil.ReadDir(foundDir)
	if err != nil {
		t.Fatalf("Error reading found directory: %v", err)
	}
	if len(foundImages) != len(expectedImages) {
		t.Fatalf("Unexpected number of found images. Expected %d, but found %d", len(expectedImages), len(foundImages))
	}
	for _, foundImg := range foundImages {
		foundImgPath := filepath.Join(foundDir, foundImg.Name())
		if !contains(expectedImages, foundImgPath) {
			t.Fatalf("Found unexpected image: %s", foundImgPath)
		}
	}
}

func createTestImageFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
