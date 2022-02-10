package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongo(db *mongo.Database) *MongoRepo {
	return &MongoRepo{
		db: db,
	}
}

type MongoRepo struct {
	db *mongo.Database
}

func (m MongoRepo) Get(ctx context.Context, username string, mClient interface{}) error {
	err := m.db.Collection("clients").FindOne(ctx, bson.M{"username": username}).Decode(&mClient)
	return err
}

func (m MongoRepo) Set(ctx context.Context, mClient interface{}) error {
	_, err := m.db.Collection("clients").InsertOne(ctx, &mClient)
	return err
}
