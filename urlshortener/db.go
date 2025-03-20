package main

import (
	"context"
	"os"

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
	var err error
	m.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	return err
}

func (m *MongoDB) StoreURL(url string, key string) error {
	collection := m.client.Database("db").Collection("urls")

	_, err := collection.InsertOne(context.Background(), bson.D{{"url", url}, {"key", key}})
	return err
}

func (m *MongoDB) GetURL(key string) (string, error) {
	collection := m.client.Database("db").Collection("urls")

	var result struct {
		URL string
	}
	if err := collection.FindOne(context.Background(), bson.D{{"key", key}}).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}
