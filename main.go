package main

import (
	"context"
	"ex-011-go-web-jpa-postgres/app"
	"ex-011-go-web-jpa-postgres/controller"
	"ex-011-go-web-jpa-postgres/model"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

func main() {

	model.ConnectMongoDB()
	defer model.DisconnectMongoDB()

	collection := model.GetLogCollection()

	logs := []model.LogEntry{
		{"Create", "Contact", "NewYork", primitive.NewDateTimeFromTime(time.Now())},
		{"Update", "Contact", "London", primitive.NewDateTimeFromTime(time.Now())},
		{"Create", "Account", "Jenny", primitive.NewDateTimeFromTime(time.Now())},
	}
	model.InsertManyLogEntries(logs)

	model.InsertOneLogEntry(model.LogEntry{
		Operation:  "Update",
		AppEntity:  "Account",
		EntityName: "Jenny",
		CreateDate: primitive.NewDateTimeFromTime(time.Now()),
	})

	updateTo := model.UpdateLogDetails("entity_name", "London", "York")

	// create a value into which the result can be decoded
	var result model.LogEntry
	err := collection.FindOne(context.TODO(), updateTo).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)


	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(20)
	filter := bson.M{}
	// Here's an array in which you can store the decoded documents
	var results []*model.LogEntry
	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		fmt.Println(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem model.LogEntry
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

	deleteFilter := bson.D{{"operation", "Create"}}
	deleteResult, err := collection.DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Deleted %v documents in the logEntries collection\n", deleteResult.DeletedCount)



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