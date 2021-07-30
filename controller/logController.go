package controller

import (
	"ex-011-go-web-jpa-postgres/model"
	u "ex-011-go-web-jpa-postgres/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

var GetLastLogEntries = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	limit, err := strconv.Atoi(params["limit"])

	if err != nil {
		//Переданное количество записей лога не является целым числом
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := model.FindLastLogEntries(bson.M{}, int64(limit))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}


func TestInsertSampleData(){
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
	model.FindOneLogEntry(updateTo)
	model.FindLastLogEntries(bson.M{}, 20)
	model.DeleteManyLogEntries(bson.M{"operation": "Create"})
}

