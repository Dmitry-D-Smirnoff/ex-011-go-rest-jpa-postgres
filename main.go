package main

import (
	"context"
	"ex-011-go-web-jpa-postgres/app"
	"ex-011-go-web-jpa-postgres/controller"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func disconnectMongoDB(client *mongo.Client){
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {

	// Create client
	client, err := mongo.NewClient(options.Client().
		ApplyURI(os.Getenv("mongodb_uri")))
	if err != nil {
		log.Fatal(err)
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer disconnectMongoDB(client)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("ex-011-database").Collection("trainers")
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 10, "Pewter City"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	trainers := []interface{}{misty, brock}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	updateFilter := bson.D{{"name", "Ash"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// create a value into which the result can be decoded
	var result Trainer
	err = collection.FindOne(context.TODO(), updateFilter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)


	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(20)
	filter := bson.M{}
	// Here's an array in which you can store the decoded documents
	var results []*Trainer
	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	deleteFilter := bson.D{{"age", 10}}
	deleteResult, err := collection.DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)



	/* MAIN CODE */

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // добавляем middleware проверки JWT-токена

	router.HandleFunc("/api/user/new", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controller.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts/{id}", controller.GetContactsFor).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(app.HandleNotFound)

	port := os.Getenv("PORT") //Получить порт из файла .env; мы не указали порт, поэтому при локальном тестировании должна возвращаться пустая строка
	if port == "" {
		port = "8000" //localhost
	}
	err = http.ListenAndServe(":" + port, router) //Запустите приложение, посетите localhost:8000/api

	if err != nil {
		fmt.Print(err)
	}
}