package main

import (
	"ex-011-go-web-jpa-postgres/app"
	"ex-011-go-web-jpa-postgres/controller"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)


func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // добавляем middleware проверки JWT-токена

	router.HandleFunc("/api/user/new", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controller.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts/{id}", controller.GetContactsFor).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(app.HandleNotFound)
	
	err := http.ListenAndServe(":8000", router) //Запустите приложение, посетите localhost:8000/api

	if err != nil {
		fmt.Print(err)
	}
}