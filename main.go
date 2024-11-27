package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// TodoItem Todo T TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
type TodoItem struct {
	Item   string `json:"item"`
	Id     string `json:"id"`
	Status bool   `json:"status"`
}

func main() {
	var todos = make([]string, 0)
	var todoItem = make(map[string]TodoItem)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(todos)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	})

	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		var t TodoItem
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todos = append(todos, t.Item)
		w.WriteHeader(http.StatusCreated)
		return
	})

	mux.HandleFunc("DELETE /todo", func(w http.ResponseWriter, r *http.Request) {
		var t TodoItem
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 3 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		todoID := pathParts[2]

		// Checking if todo exists
		if _, exists := todoItem[todoID]; !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		//Deleting the todo
		delete(todoItem, todoID)

		w.WriteHeader(http.StatusOK)
		response := map[string]string{
			"message": "Todo item deleted successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
