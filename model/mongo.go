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

func FindOneLogEntry(criteria primitive.D) LogEntry {
	// create a value into which the result can be decoded
	var result LogEntry
	err := GetLogCollection().FindOne(context.TODO(), criteria).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
	return result
}

func FindManyLogEntries(criteria bson.M, limit int64) []*LogEntry {
	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(limit)
	// Here's an array in which you can store the decoded documents
	var results []*LogEntry
	// Passing nil as the filter matches all documents in the collection
	cur, err := GetLogCollection().Find(context.TODO(), criteria, options)
	if err != nil {
		fmt.Println(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem LogEntry
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results
}

func DeleteManyLogEntries(criteria bson.M) *mongo.DeleteResult {
	deleteResult, err := GetLogCollection().DeleteMany(context.TODO(), criteria)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Deleted %v documents in the logEntries collection\n", deleteResult.DeletedCount)
	return deleteResult
}

