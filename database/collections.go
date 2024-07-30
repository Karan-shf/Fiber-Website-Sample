package database

import (
	"context"

	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var AdminCollection *mongo.Collection
var NewsCollection *mongo.Collection

func DeployCollections() {

	Connect_Client()

	UserCollection = Client.Database("tarafdari_sample_db").Collection("users")
	AdminCollection = Client.Database("tarafdari_sample_db").Collection("admins")
	NewsCollection = Client.Database("tarafdari_sample_db").Collection("news")

	log.Println("Collections Deployed Succesfully")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := AdminCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal("error while adding index: ", err)
	}
	log.Println("unique index added to admin collection")
}

// -- making username unique in user collection --
// indexModel := mongo.IndexModel{
// 	Keys:    bson.M{"username": 1},
// 	Options: options.Index().SetUnique(true),
// }
// _, err := UserCollection.Indexes().CreateOne(context.Background(), indexModel)
// if err != nil {
// 	log.Fatal("error while adding index: ", err)
// }
