package tensorflowAPI

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	"github.com/galeone/tensorflow/tensorflow/go/op"
)

var (
	graph  *tf.Graph
	labels []string
)

func loadModel() (*tf.Graph, []string, error) {
	// Load inception model
	model, err := ioutil.ReadFile("/model/tensorflow_inception_graph.pb")
	if err != nil {
		return nil, nil, err
	}
	graph = tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return nil, nil, err
	}
	// Load labels
	labelsFile, err := os.Open("/model/imagenet_comp_graph_label_strings.txt")
	if err != nil {
		return nil, nil, err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	fmt.Println(labels[0])

	return graph, labels, nil
}

func ClassifyImage(imagebuff []byte) error {

	imageTensor, err := tf.NewTensor(imagebuff)
	if err != nil {
		return err
	}

	normalized, err := normalizeImage(imageTensor)
	if err != nil {
		return err
	}

	model, labels, err := loadModel()
	if err != nil {
		return err
	}

	res, err := getProbabilities(model, normalized)
	if err != nil {
		return err
	}

	_ = labels
	fmt.Println(res.Value())
	return nil
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

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{
			input: tensor,
		},
		[]tf.Output{
			output,
		},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func getNormalizedGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)

	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	output = op.Sub(s,
		// make it 224x224: inception specific
		op.ResizeBilinear(s,
			op.ExpandDims(s,
				op.Cast(s, decode, tf.Float),
				op.Const(s.SubScope("make_batch"), int32(0))),
			op.Const(s.SubScope("size"), []int32{224, 224})),
		// mean = 117: inception specific
		op.Const(s.SubScope("mean"), float32(117)))
	graph, err = s.Finalize()

	return graph, input, output, err
}

func getProbabilities(model *tf.Graph, tensor *tf.Tensor) (tf.Tensor, error) {
	session, err := tf.NewSession(model, nil)
	if err != nil {
		return tf.Tensor{}, err
	}
	defer session.Close()

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			model.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			model.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		return tf.Tensor{}, err
	}

	return *output[0], nil
}
