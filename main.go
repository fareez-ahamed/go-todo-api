package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoPayload struct {
	Desc string `json:"description"`
}

var todosStore = TodoStore{
	todos: []Todo{{1, "Read a book", false},
		{2, "Exercise", true}},
}

func main() {
	fmt.Print("Server is running")
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/api/todos", getTodosHandler).Methods("GET")
	router.HandleFunc("/api/todos/{id}", getTodoDetailHandler).Methods("GET")
	router.HandleFunc("/api/todos", addTodoHandler).Methods("POST")
	router.HandleFunc("/api/todos/{id}/mark_completed", markCompletedHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}/mark_incomplete", markIncompleteHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", updateTodoHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodoHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:9000", router))
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("completed") == "true" {
		writeJson(w, http.StatusOK, todosStore.GetByStatus(true))
	} else if r.URL.Query().Get("completed") == "false" {
		writeJson(w, http.StatusOK, todosStore.GetByStatus(false))
	} else {
		writeJson(w, http.StatusOK, todosStore.GetAll())
	}
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	var payload TodoPayload
	json.NewDecoder(r.Body).Decode(&payload)
	error := todosStore.Update(id, payload.Desc)
	if err != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully updated")
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	err = todosStore.Delete(id)
	if err != nil {
		writeMessage(w, http.StatusNotFound, err.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully deleted")
}

func getTodoDetailHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	todo, err := todosStore.GetTodoDetail(id)
	if err != nil {
		writeMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJson(w, http.StatusOK, todo)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	var payload TodoPayload
	json.NewDecoder(r.Body).Decode(&payload)
	todosStore.AddTodo(payload.Desc)
	writeMessage(w, http.StatusOK, "Successfully added")
}

func markCompletedHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	error := todosStore.SetCompleted(id, true)
	if err != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully marked completed")
}

func markIncompleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	error := todosStore.SetCompleted(id, false)
	if err != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully marked incomplete")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]string{
		"message": "Hello World!",
	})
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeMessage(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{
		"message": message,
	})
}
