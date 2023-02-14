package tensorflowAPI

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	tf "github.com/galeone/tensorflow/tensorflow/go"

	"github.com/galeone/tensorflow/tensorflow/go/op"
	"github.com/l-pavlova/image-master/logging"
)

const (
	inception_model_path   = "tensorflowAPI/model/tensorflow_inception_graph.pb"
	inception_model_labels = "tensorflowAPI/model/imagenet_comp_graph_label_strings.txt"
	coco_model_path        = "tensorflowAPI/model/frozen_inference_graph.pb"
	coco_labels_path       = "tensorflowAPI/model/coco_labels.txt"
)

var (
	graph  *tf.Graph
	labels []string
)

type Label struct {
	Label       string
	Probability float32
}

type TensorFlowClient struct {
	logger logging.ImageMasterLogger
}

func NewTensorFlowClient(logger logging.ImageMasterLogger) *TensorFlowClient {
	tfClient := &TensorFlowClient{
		logger: logger,
	}

	graph = initModel(coco_model_path, logger)
	labels = initLables(coco_labels_path, logger)
	return tfClient
}

func (t *TensorFlowClient) ClassifyImage(imagebuff string) error {

	imageTensor, err := tf.NewTensor(imagebuff)
	if err != nil {
		return err
	}
	t.logger.Log("info", "image tensor created")

	normalized, err := normalizeImage(imageTensor)
	if err != nil {
		return err
	}

	t.logger.Log("info", "image normalzied")

	probabilities, _, _, err := t.getProbabilities(normalized)
	if err != nil {
		return err
	}

	t.logger.Log("info", "probabilities gotten.")
	top5 := getTopLabel(labels, probabilities)
	t.logger.Log("info", "top 5 matches are", top5)
	return nil
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

	logger.Log("info", "retrieved labels count: ", len(labels))
	return graph
}

func initLables(labelsFilePath string, logger logging.ImageMasterLogger) []string {
	// Load labels
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

	fmt.Println("Starting new tf session")
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	fmt.Println("Normalizing image")
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

	fmt.Println("Normalized image is:")
	fmt.Print(normalized)
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
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	inputop := graph.Operation("image_tensor")
	// Output ops
	o1 := graph.Operation("detection_boxes")
	o2 := graph.Operation("detection_scores")
	o3 := graph.Operation("detection_classes")
	o4 := graph.Operation("num_detections")

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
		fmt.Print(err.Error())
		log.Fatal(err)
	}

	// Outputs
	probabilities := output[1].Value().([][]float32)[0]
	classes := output[2].Value().([][]float32)[0]
	boxes := output[0].Value().([][][]float32)[0]

	return probabilities, classes, boxes, nil
}

func getTopLabel(labels []string, probabilities []float32) []Label {
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

	fmt.Println(resultLabels[:5])
	return resultLabels[:5]
}
