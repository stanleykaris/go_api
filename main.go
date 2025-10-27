package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// TodoItem Todo T TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
type TodoItem struct {
	Item   string `json:"item"`
	Id     string `json:"id"`
	Status bool   `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	// start the server in goroutine
	go func() {
		log.Printf("Listening on port %s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v\nTry using a different port", err)
		}
	case sig := <-shutdown:
		log.Printf("Received signal %s, shutting down gracefully...\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Fatalf("Error shutting down server: %v", err)
		}
		log.Println("Server shut down gracefully")
		return
	}

}
