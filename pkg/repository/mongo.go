package repository

import (
	"context"
	"fmt"
	mo "github.com/mahjadan/go-mfa/pkg/models"
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
	err := m.db.Collection("clients").FindOne(ctx, bson.M{"username": username}).Decode(mClient)
	return err
}

func (m MongoRepo) Exists(ctx context.Context, username string) (bool, error) {
	singleResult := m.db.Collection("clients").FindOne(ctx, bson.M{"username": username})
	if singleResult.Err() != nil {
		if singleResult.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, singleResult.Err()
	}
	return true, nil
}

func (m MongoRepo) Set(ctx context.Context, mClient interface{}) error {
	_, err := m.db.Collection("clients").InsertOne(ctx, mClient)
	return err
}
func (m MongoRepo) Update(ctx context.Context, mClient mo.MongoClient) error {
	filter := bson.M{
		"username": mClient.Username,
	}
	update := bson.D{
		{"$push", bson.M{"MFA": mClient.MFA[0]}},
	}
	updateOne, err := m.db.Collection("clients").UpdateOne(ctx, filter, update)
	fmt.Println("updateOne result: ", updateOne)
	return err
}
