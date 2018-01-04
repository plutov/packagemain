package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

const (
	graphFile  = "/model/tensorflow_inception_graph.pb"
	labelsFile = "/model/imagenet_comp_graph_label_strings.txt"
)

// Label type
type Label struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

// Labels type
type Labels []Label

func (a Labels) Len() int           { return len(a) }
func (a Labels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Labels) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func main() {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")
	if len(os.Args) < 2 {
		log.Fatalf("usage: imgrecognition <image_url>")
	}
	fmt.Printf("url: %s\n", os.Args[1])

	// Get image from URL
	response, e := http.Get(os.Args[1])
	if e != nil {
		log.Fatalf("unable to get image from url: %v", e)
	}
	defer response.Body.Close()

	modelGraph, labels, err := loadModel()
	if err != nil {
		log.Fatalf("unable to load model: %v", err)
	}

	// Make tensor
	tensor, err := normalizeImage(response.Body)
	if err != nil {
		log.Fatalf("unable to make a tensor from image: %v", err)
	}

	// Create a session for inference over modelGraph.
	session, err := tf.NewSession(modelGraph, nil)
	if err != nil {
		log.Fatalf("could not init session: %v", err)
	}
	defer session.Close()

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			modelGraph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			modelGraph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatalf("could not run inference: %v", err)
	}

	res := findBestLabels(labels, output[0].Value().([][]float32)[0])
	for _, l := range res {
		fmt.Printf("label: %s, probability: %d%%\n", l.Label, int(l.Probability*100))
	}
}

func loadModel() (*tf.Graph, []string, error) {
	// Load inception model
	model, err := ioutil.ReadFile(graphFile)
	if err != nil {
		return nil, nil, err
	}
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return nil, nil, err
	}

	// Load labels
	labelsFile, err := os.Open(labelsFile)
	if err != nil {
		return nil, nil, err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	return graph, labels, scanner.Err()
}

func findBestLabels(labels []string, probabilities []float32) []Label {
	var resultLabels []Label
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, Label{Label: labels[i], Probability: p})
	}

	sort.Sort(Labels(resultLabels))
	return resultLabels[:5]
}

func normalizeImage(body io.ReadCloser) (*tf.Tensor, error) {
	var buf bytes.Buffer
	io.Copy(&buf, body)

	tensor, err := tf.NewTensor(buf.String())
	if err != nil {
		return nil, err
	}

	graph, input, output, err := makeTransformImageGraph()
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
		return nil, err
	}

	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func makeTransformImageGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	// 3 return RGB image
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	// Div and Sub perform (value-Mean)/Scale for each pixel
	output = op.Div(s,
		op.Sub(s,
			// Resize to 224x224 with bilinear interpolation
			op.ResizeBilinear(s,
				// Create a batch containing a single image
				op.ExpandDims(s,
					// Casts image to float type
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{224, 224})),
			op.Const(s.SubScope("mean"), float32(117))),
		op.Const(s.SubScope("scale"), float32(1)))
	graph, err = s.Finalize()
	return graph, input, output, err
}
