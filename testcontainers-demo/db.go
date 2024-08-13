package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	Init() error
	StoreURL(url string, key string) error
	GetURL(key string) (string, error)
}

type MongoDB struct {
	client *mongo.Client
}

func (m *MongoDB) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	return err
}

func (m *MongoDB) StoreURL(url string, key string) error {
	collection := m.client.Database("db").Collection("urls")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bson.D{{"url", url}, {"key", key}})
	return err
}

func (m *MongoDB) GetURL(key string) (string, error) {
	collection := m.client.Database("db").Collection("urls")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var result struct {
		URL string
	}
	if err := collection.FindOne(ctx, bson.D{{"key", key}}).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}
