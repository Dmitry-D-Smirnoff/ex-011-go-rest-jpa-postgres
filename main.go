package main

import (
	"ex-011-go-web-jpa-postgres/app"
	"ex-011-go-web-jpa-postgres/controller"
	"ex-011-go-web-jpa-postgres/model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	model.ConnectMongoDB()
	defer model.DisconnectMongoDB()
	//controller.TestInsertSampleData()

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // добавляем middleware проверки JWT-токена
	router.HandleFunc("/api/user/new", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controller.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts/{id}", controller.GetContactsFor).Methods("GET")
	router.HandleFunc("/api/log/recent/{limit}", controller.GetLastLogEntries).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(app.HandleNotFound)

	port := os.Getenv("PORT") //Получить порт из файла .env; мы не указали порт, поэтому при локальном тестировании должна возвращаться пустая строка
	if port == "" {
		port = "8000" //localhost
	}
	err := http.ListenAndServe(":" + port, router) //Запустите приложение, посетите localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}