package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type LogEntry struct {
	Operation  string             `bson:"operation" json:"operation"`
	AppEntity  string             `bson:"app_entity" json:"appEntity"`
	EntityName string             `bson:"entity_name" json:"entityName"`
	CreateDate primitive.DateTime `bson:"create_date" json:"createDate"`
}

var currentClient *mongo.Client

func GetLogCollection() (collection *mongo.Collection){
	if currentClient == nil {
		ConnectMongoDB()
	}else{
		//TODO: Ping and if fails -> Disconnect/Connect again
		//currentClient.Ping(context.TODO(),???)
	}

	//TODO: Environmental variables for Database and Log collection
	return currentClient.Database("ex-011-database").Collection("logEntries")
}

func ConnectMongoDB(){
	// Create client
	client, err := mongo.NewClient(options.Client().
		ApplyURI(os.Getenv("mongodb_uri")))
	if err != nil {
		fmt.Println(err)
	}
	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB!")

	currentClient = client
}

func DisconnectMongoDB(){
	err := currentClient.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err)
	}else{
		currentClient = nil
		fmt.Println("Connection to MongoDB closed.")
	}

}

func InsertManyLogEntries(logEntries []LogEntry){
	logInterface := make([]interface{}, len(logEntries))
	for i, v := range logEntries {
		logInterface[i] = v
	}

	insertManyResult, err := GetLogCollection().InsertMany(context.TODO(), logInterface)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func InsertOneLogEntry(logEntry LogEntry){
	insertOneResult, err := GetLogCollection().InsertOne(context.TODO(), logEntry)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Inserted one document: ", insertOneResult.InsertedID)
}

func UpdateLogDetails(mongoKey string, oldValue string, newValue string) primitive.D {
	updateFilter := bson.D{{mongoKey, oldValue}}
	updateTo := bson.D{{mongoKey, newValue},{"create_date", primitive.NewDateTimeFromTime(time.Now())}}
	update := bson.D{{"$set", updateTo}}
	updateResult, err := GetLogCollection().UpdateMany(context.TODO(), updateFilter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return updateTo
}



