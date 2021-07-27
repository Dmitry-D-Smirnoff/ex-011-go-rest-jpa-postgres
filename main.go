package main

import (
	"context"
	"ex-011-go-web-jpa-postgres/app"
	"ex-011-go-web-jpa-postgres/controller"
	"fmt"
	"github.com/gorilla/mux"
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

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("persons")

	fmt.Println(collection)



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