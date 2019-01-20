package api

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

const topicName = "randomNumbers"

// Send generates random integer and sends it to Cloud Pub/Sub
func Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	topic := client.Topic(topicName)

	rand.Seed(time.Now().UnixNano())

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strconv.Itoa(rand.Intn(1000))),
	})
	id, err := result.Get(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}
