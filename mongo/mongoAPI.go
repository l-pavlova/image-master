package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MONGO_HOST      = "mongodb://localhost:27017"
	MONGO_DBNAME    = "Classificator"
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

func (m *Mongo) GetImageClassification(imagePath string) (ic ImageClassification, ok bool) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
		ok = false
	}

	defer Close(client, ctx, cancel)

	filter := bson.M{
		"imagepath": imagePath,
	}

	imClassification := FindOne(client, ctx, MONGO_DBNAME, COLLECTION_NAME, filter, nil)
	err = imClassification.Decode(&ic)

	if err != nil {
		ok = false
		return
	}

	ok = true
	return ic, ok
}

func GetAllImageClassifications() (ics []ImageClassification, err error) {
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
