package main

import (
	"bugsmirror/database"
	"bugsmirror/handlers"
	"bugsmirror/migrations"
	"fmt"
	"log"
	"net/http"
)

func main() {

	db, err := database.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    defer db.Close()

    // Run database migrations
    if err := migrations.Run(db); err != nil {
        log.Fatalf("Failed to run database migrations: %v", err)
    }

    http.HandleFunc("/users", handlers.CreateUser(db))
    http.HandleFunc("/user", handlers.GetUser(db))
    http.HandleFunc("/users/all", handlers.GetAllUsers(db))
    http.HandleFunc("/user/edit", handlers.EditUser(db))
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "App is working!")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
