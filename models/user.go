package models

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStruct struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstname"`
	LastName  string             `bson:"lastname"`
	Age       int                `bson:"age"`
	Email     string             `bson:"email"`
}

func NewUserCollection(db *mongo.Client) (mongo *mongo.Collection) {
	mongo = db.Database("users").Collection("users")
	return
}

func InitDB() (*mongo.Client, context.Context, error) {
	fmt.Println("Connecting to MongoDB...")

	mongoURI := "mongodb://localhost:27017/"
	if len(os.Getenv("DB_HOST")) > 0 {
		mongoURI = os.Getenv("DB_HOST")
	}
	mongoCtx := context.Background()
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, mongoCtx, err
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		return nil, mongoCtx, err
	}
	fmt.Println("Connected to Mongodb")
	return db, mongoCtx, nil
}
