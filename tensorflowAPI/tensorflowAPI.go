package tensorflowAPI

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	tf "github.com/galeone/tensorflow/tensorflow/go"

	"github.com/galeone/tensorflow/tensorflow/go/op"
	"github.com/l-pavlova/image-master/logging"
)

const (
	INCEPTION_MODEL_PATH          = "tensorflowAPI/model/tensorflow_inception_graph.pb"
	INCEPTION_MODEL_LABELS        = "tensorflowAPI/model/imagenet_comp_graph_label_strings.txt"
	COCO_MODEL_PATH               = "tensorflowAPI/model/frozen_inference_graph.pb"
	COCO_LABELS_PATH              = "tensorflowAPI/model/coco_labels.txt"
	DEPTH_OF_LABEL_CLASSIFICATION = 5
)

var lock = &sync.Mutex{}

type Label struct {
	Label       string
	Probability float32
}

type Model struct {
	graph  *tf.Graph
	labels []string
}

var singleModelInstance *Model

// a singleton is used for loading the heavy load data for the .pb tensorflow model and the labels for it, so all instances of the Tensorflow client are initialized with the same loaded data
// instead of all of them loading the same thing over and over
func getModelInstance() *Model {
	if singleModelInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleModelInstance == nil {
			fmt.Println("Creating single instance now.")
			singleModelInstance = &Model{
				graph:  initModel(COCO_MODEL_PATH),
				labels: initLables(COCO_LABELS_PATH),
			}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleModelInstance
}

type TensorFlowClient struct {
	logger logging.ImageMasterLogger
	graph  *tf.Graph
	labels []string
}

func NewTensorFlowClient(logger logging.ImageMasterLogger) *TensorFlowClient {
	tfClient := &TensorFlowClient{
		logger: logger,
	}

	model := getModelInstance()
	tfClient.graph = model.graph
	tfClient.labels = model.labels
	return tfClient
}

// function that takes an image and returns the top probable objects in it
func (t *TensorFlowClient) ClassifyImage(img image.Image) ([]Label, []float32, [][]float32, error) {

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error occurred during img buffering %w", err)
	}

	imageTensor, err := tf.NewTensor(buf.String())
	if err != nil {
		t.logger.Log("warning", "failed parsing imagebuffer to tensor", err)
		return nil, nil, nil, err
	}
	t.logger.Log("info", "image tensor created")

	normalized, err := normalizeImage(imageTensor)
	if err != nil {
		t.logger.Log("warning", "failed normalizing tensor", err)
		return nil, nil, nil, err
	}

	t.logger.Log("info", "image normalized")
	probabilities, classes, boxes, err := t.getProbabilities(normalized)
	if err != nil {
		t.logger.Log("warning", "failed calculating probabilities", err)
		return nil, nil, nil, err
	}

	t.logger.Log("info", "probabilities received.")

	topN := getTopLabels(t.labels, probabilities, DEPTH_OF_LABEL_CLASSIFICATION)
	t.logger.Log("info", "top matches are", DEPTH_OF_LABEL_CLASSIFICATION, topN)
	return topN, classes, boxes, nil
}

// function to load the model .pb file
func initModel(modelPath string) *tf.Graph {
	model, err := ioutil.ReadFile(modelPath)
	if err != nil {
		return nil
	}

	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return nil
	}

	return graph
}

// function to Load labels
func initLables(labelsFilePath string) []string {
	// Load labels
	labels := make([]string, 0, 10)
	labelsFile, err := os.Open(labelsFilePath)
	if err != nil {
		return nil
	}

	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)

	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	fmt.Println("info", "retrieved labels count: ", len(labels))
	return labels
}

// function that takes the initial tensor with the image and formats it so that it is usable by the model
func normalizeImage(tensor *tf.Tensor) (*tf.Tensor, error) {
	graph, input, output, err := getNormalizedGraph()
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{
			input: tensor,
		},
		[]tf.Output{
			output,
		},
		nil)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func getNormalizedGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.ExpandDims(s,
		op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)),
		op.Const(s.SubScope("make_batch"), int32(0)))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func (t *TensorFlowClient) getProbabilities(tensor *tf.Tensor) ([]float32, []float32, [][]float32, error) {
	session, err := tf.NewSession(t.graph, nil)
	if err != nil {
		t.logger.Log("error", err.Error())
	}
	defer session.Close()

	inputop := t.graph.Operation("image_tensor")
	// Output ops
	o1 := t.graph.Operation("detection_boxes")
	o2 := t.graph.Operation("detection_scores")
	o3 := t.graph.Operation("detection_classes")
	o4 := t.graph.Operation("num_detections")

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			inputop.Output(0): tensor,
		},
		[]tf.Output{
			o1.Output(0),
			o2.Output(0),
			o3.Output(0),
			o4.Output(0),
		},
		nil)

	if err != nil {
		t.logger.Log("error", err.Error())
	}

	// Outputs
	probabilities := output[1].Value().([][]float32)[0]
	classes := output[2].Value().([][]float32)[0]
	boxes := output[0].Value().([][][]float32)[0]

	return probabilities, classes, boxes, nil
}

func getTopLabels(labels []string, probabilities []float32, count int) []Label {
	var resultLabels []Label
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, Label{Label: labels[i], Probability: p})
	}

	sort.Slice(resultLabels, func(i, j int) bool {
		return resultLabels[i].Probability > resultLabels[j].Probability
	})

	return resultLabels[:count]
}
