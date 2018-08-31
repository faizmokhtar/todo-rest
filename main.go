package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_at"`
}

var todoStore = make(map[string]Todo)
var id int = 0

// GET - /api/todos
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	for _, v := range todoStore {
		todos = append(todos, v)
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(todos)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// POST - /api/todos
func PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		panic(err)
	}

	todo.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	todoStore[k] = todo

	j, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GET - /api/notes/{id}
func GetTodoShowHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	todo := todoStore[k]

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// PUT - /api/notes/{id}
func PutTodoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	var updatedTodo Todo

	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		panic(err)
	}
	if todo, ok := todoStore[k]; ok {
		updatedTodo.CreatedOn = todo.CreatedOn

		delete(todoStore, k)
		todoStore[k] = updatedTodo
	} else {
		log.Printf("could not find key of Todo %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE - /api/notes/{id}
func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := todoStore[k]; ok {
		delete(todoStore, k)
	} else {
		log.Printf("could not find key of Todo %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/todos", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/todos/{id:[0-9]+}", GetTodoShowHandler).Methods("GET")
	r.HandleFunc("/api/todos", PostTodoHandler).Methods("POST")
	r.HandleFunc("/api/todos/{id:[0-9]+}", PutTodoHandler).Methods("PUT")
	r.HandleFunc("/api/todos/{id:[0-9]+}", DeleteTodoHandler).Methods("DELETE")

	server := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
