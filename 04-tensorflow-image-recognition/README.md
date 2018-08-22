Tensorflow is not a Machine Learning specific library, instead, is a general purpose computation library that represents computations with graphs. Its core is implemented in C++ and there are also bindings for different languages, including Go.

In the last few years the field of machine learning has made tremendous progress on addressing the difficult problem of image recognition.

One of the challenges with machine learning is figuring out how to deploy trained models into production environments. After training your model, you can "freeze" it and export it to be used in a production environment.

For some common-use-cases we're beginning to see organizations sharing their trained models, you can find some in the TensorFlow Models repo https://github.com/tensorflow/models.

In this video we'll use one of them, called Inception to recognize an image. https://github.com/tensorflow/models/tree/master/research/inception/inception

We'll build a small command line application that takes URL to an image as input and outputs labels in order.

First of all we need to install TensorFlow, and here Docker can be really helpful, because installation of Tensorflow may be complicated. There is a Docker image with Tensorflow, but without Go, so I found an image with Tensorflow plus Go to reduce my Dockerfile.

https://github.com/ctava/tensorflow-go

Download Inception data:
http://download.tensorflow.org/models/inception5h.zip

Let's start with simple main.go file to test if our Dockerfile works.

```
package main

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: imgrecognition <image_url>")
	}
	fmt.Printf("url: %s\n", os.Args[1])
}
```

```
docker build -t imgrecognition .
docker run imgrecognition https://www.iaspaper.net/wp-content/uploads/2017/10/Rabbit-Essay.jpg
```

Let's get our image from the provided URL:
```
// Get image from URL
response, e := http.Get(os.Args[1])
if e != nil {
	log.Fatalf("unable to get image from url: %v", e)
}
defer response.Body.Close()
```

Now we need to load our model. Model contains graph and labels in 2 files:
```
const (
	graphFile = "/model/imagenet_comp_graph_label_strings.txt"
	labelsFile = "/model/imagenet_comp_graph_label_strings.txt"
)

graph, labels, err := loadModel()
if err != nil {
	log.Fatalf("unable to load model: %v", err)
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
	labelsFile, err := os.Open()
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
```

Here finally we start using tensorflow Go package.
To be able to work with our image we need to normalize it, because Inception model expects it to be in a certain format, it uses images from ImageNet, and they are 224x224. But that's a bit tricky. Let's see:
 - NewTensor converts from a Go value to a Tensor
 - Build a graph of our image
 - Init a session, because all Graph operations in Tensorflow are done with sessions.
 - Run the session to normalize image, using input and output.
 - normalized[0] contains normalized Tensor.
 - In makeTransformImageGraph we define the rules of normalization.

```
func normalizeImage(body io.ReadCloser) (*tensorflow.Tensor, error) {
	var buf bytes.Buffer
	io.Copy(&buf, body)

	tensor, err := tensorflow.NewTensor(buf.String())
	if err != nil {
		return nil, err
	}

	graph, input, output, err := getNormalizedGraph()
	if err != nil {
		return nil, err
	}

	session, err := tensorflow.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	normalized, err := session.Run(
		map[tensorflow.Output]*tensorflow.Tensor{
			input: tensor,
		},
		[]tensorflow.Output{
			output,
		},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func getNormalizedGraph() (graph *tensorflow.Graph, input, output tensorflow.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tensorflow.String)
	// 3 return RGB image
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	// Sub: returns x - y element-wise
	output = op.Sub(s,
		// make it 224x224: inception specific
		op.ResizeBilinear(s,
			// inserts a dimension of 1 into a tensor's shape.
			op.ExpandDims(s,
				// cast image to float type
				op.Cast(s, decode, tensorflow.Float),
				op.Const(s.SubScope("make_batch"), int32(0))),
			op.Const(s.SubScope("size"), []int32{224, 224})),
		// mean = 117: inception specific
		op.Const(s.SubScope("mean"), float32(117)))
	graph, err = s.Finalize()

	return graph, input, output, err
}
```

Fnally we need to init one more session on our initial model graph to find matches.

```
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
```

It will return list of probabilities for each label. What we need now is to loop over all probabilities and find label in `labels` slice. And print top 5.

```
res := getTopFiveLabels(labels, output[0].Value().([][]float32)[0])
for _, l := range res {
	fmt.Printf("label: %s, probability: %.2f%%\n", l.Label, l.Probability*100)
}

func getTopFiveLabels(labels []string, probabilities []float32) []Label {
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
```

https://www.zoo-berlin.de/fileadmin/_processed_/4/4/csm_Meng_Meng_Baby_1_88cad0f74f.jpg

Also let's skip those warnings:
```
os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")
```

Here we worked with pre-trained model, let's try this program with something unusual, like ... Gopher.
```
docker run imgrecognition https://i.pinimg.com/736x/12/5c/e0/125ce0baff3271761ca61843eccf7985.jpg
```

Mouse?? no! But it's possible to train our models from Go in TensorFlow, and I will definitely do a video about it.