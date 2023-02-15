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

var ()

type Label struct {
	Label       string
	Probability float32
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

	tfClient.graph = initModel(COCO_MODEL_PATH, logger)
	tfClient.labels = initLables(COCO_LABELS_PATH, logger)
	return tfClient
}

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

func initModel(modelPath string, logger logging.ImageMasterLogger) *tf.Graph {
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

func initLables(labelsFilePath string, logger logging.ImageMasterLogger) []string {
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

	logger.Log("info", "retrieved labels count: ", len(labels))
	return labels
}

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
