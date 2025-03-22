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

type doc struct {
	Url string
	Key string
}

func (m *MongoDB) Init() error {
	opts := options.Client().ApplyURI(os.Getenv("MONGO_ADDR"))
	var err error
	m.client, err = mongo.Connect(context.Background(), opts)
	return err
}

func (m *MongoDB) StoreURL(url string, key string) error {
	collection := m.client.Database("db").Collection("urls")

	_, err := collection.InsertOne(context.Background(), doc{
		Url: url,
		Key: key,
	})
	return err
}

func (m *MongoDB) GetURL(key string) (string, error) {
	collection := m.client.Database("db").Collection("urls")

	var res doc
	if err := collection.FindOne(context.Background(), bson.D{{"key", key}}).Decode(&res); err != nil {
		return "", err
	}

	return res.Url, nil
}
