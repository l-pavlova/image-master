package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MONGO_HOST      = "mongodb+srv://admin:{pass}@cluster0.nvd7u.mongodb.net/ImageMasterDB?retryWrites=true&w=majority"
	MONGO_DBNAME    = "ImageMasterDB"
	COLLECTION_NAME = "ImageClassifications"
)

type ImageClassification struct {
	ImagePath     string
	Probabilities []string
}

type Mongo struct {
}

func NewMongo() *Mongo {
	return &Mongo{}
}

// insert new image classification labels for a given imagePath
func (m *Mongo) AddImageClassification(imagePath string, probabilities []string) error {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Release resource when main function is returned.
	defer Close(client, ctx, cancel)

	imClassification := &ImageClassification{ImagePath: imagePath, Probabilities: probabilities}
	insertOneResult, err := InsertOne(client, ctx, MONGO_DBNAME, COLLECTION_NAME, imClassification)
	if err != nil {
		return err
	}

	fmt.Print(insertOneResult)
	return nil
}

// get image classification labels for a given imagePath
func (m *Mongo) GetImageClassification(imagePath string) (ic ImageClassification, err error) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
		return ImageClassification{}, err
	}

	defer Close(client, ctx, cancel)

	filter := bson.M{
		"imagepath": imagePath,
	}
	//fmt.Println("Filter: ", filter)

	imClassification := FindOne(client, ctx, MONGO_DBNAME, COLLECTION_NAME, filter, nil)
	//fmt.Print("im classification returned: ", imClassification)
	err = imClassification.Decode(&ic)

	fmt.Println("Result from db: ", ic)
	if err != nil {
		fmt.Println(err.Error())
		return ImageClassification{}, err
	}

	return ic, nil
}

// Get all image classifications
func (m *Mongo) GetAllImageClassifications() (ics []ImageClassification, err error) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer Close(client, ctx, cancel)

	filter := bson.D{{}}
	imClassifications, err := Query(client, ctx, MONGO_DBNAME, COLLECTION_NAME, filter, nil)
	if err != nil {
		return nil, err
	}

	imClassifications.Decode(&ics)
	return ics, nil
}
