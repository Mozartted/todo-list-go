package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mozartted/simple_todo_server/internal/config"
	"github.com/mozartted/simple_todo_server/internal/model"
)

type RouteHandler struct {
	// Db config.DbSync
	Db *sql.DB
}

// TODO API handling.
func (r *RouteHandler) TodoListAdd(w http.ResponseWriter, rq *http.Request) {
	defer rq.Body.Close()
	// input := make([]byte, 40)
	var generalString string

	inputData, err := io.ReadAll(rq.Body)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}
	generalString = string(inputData)
	// for {
	// 	n, err := rq.Body.Read(input)
	// 	// if err != nil {
	// 	// 	log.Fatalf("Error reading inputData")
	// 	// 	return
	// 	// }
	//
	// 	if n > 0 {
	// 		log.Print(n)
	// 		chunk := input[:n]
	// 		generalString += string(chunk)
	// 	}
	// 	// fmt.Printf("Read %d bytes: %s\n", n, string(chunk))
	//
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }

	// w.WriteHeader(200)
	var task model.TaskData

	if err := json.Unmarshal([]byte(generalString), &task); err != nil {
		log.Fatalf("error processing json data %v", err.Error())
	}

	// fmt.Printf("entered task data %v", task)

	if err := config.SaveTodoData(r.Db, task); err != nil {
		log.Fatalf("%s", err.Error())
	}

	data, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Something Went wrong %v", err.Error())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(data))

	}
}

func (r *RouteHandler) TodoListGet(w http.ResponseWriter, rq *http.Request) {
	resp := config.RetrieveAll(r.Db)

	response, err := json.Marshal(map[string][]model.TaskData{"todo": resp})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	w.WriteHeader(200)
	w.Write(response)
}

func (r *RouteHandler) TodoListUpdate(w http.ResponseWriter, rq *http.Request) {
	resp := mux.Vars(rq)
	category, err := strconv.ParseInt(resp["todoId"], 10, 32)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	updatedList := config.UpdateStatus(r.Db, uint(category))
	response, err := updatedList.ToJSON()
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	w.Write([]byte(response))
}

func (r *RouteHandler) TodoListDelete(w http.ResponseWriter, rq *http.Request) {
	resp := mux.Vars(rq)
	category, err := strconv.ParseInt(resp["todoId"], 10, 32)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	updatedList := config.DeleteTask(r.Db, uint(category))

	response, err := json.Marshal(map[string][]model.TaskData{"todo": updatedList})
	w.WriteHeader(200)
	w.Write(response)
}
