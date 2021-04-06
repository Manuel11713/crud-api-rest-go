package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID      int    `json: ID`
	Name    string `json: Name`
	Content string `json: Content`
}
type allTasks []Task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task1",
		Content: "Content1",
	},
	{
		ID:      2,
		Name:    "Task2",
		Content: "Content2",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my api compile daemon")
}
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		fmt.Fprintf(w, "insert a valid task")
		return
	}
	if newTask.Content == "" || newTask.Name == "" {
		fmt.Fprintf(w, "content and name must be provided")
		return
	}

	newTask.ID = len(tasks) + 1

	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}
func getTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "id should be a number")
		return
	}
	for _, v := range tasks {
		if v.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(v)
			return
		}

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]struct{}{})

}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "id should be a number")
		return
	}
	for idx, v := range tasks {
		if v.ID == id {
			tasks = append(tasks[:idx], tasks[idx+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(v)
			return

		}
	}
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reqBody, err := ioutil.ReadAll(r.Body)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}
	for idx, t := range tasks {
		if t.ID == id {

			var newTask Task

			if err != nil {
				fmt.Fprintf(w, "insert a valid task")
			}
			json.Unmarshal(reqBody, &newTask)
			newTask.ID = id
			fmt.Println(newTask)
			tasks = append(tasks[:idx], tasks[idx+1:]...)
			tasks = append(tasks, newTask)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}

	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexRoute)
	r.HandleFunc("/tasks", getTasks)
	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task/{id}", getTaskById).Methods("GET")
	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", r))
}
