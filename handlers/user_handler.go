package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type User struct {
	ID          int    `json:"id"`
    FirstName   string `json:"firstname"`
    LastName    string `json:"lastname"`
    DOB         string `json:"dob"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phonenumber"`
}

func CreateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var newUser User
        err := json.NewDecoder(r.Body).Decode(&newUser)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Insert user into the database
        _, err = db.Exec("INSERT INTO users (firstname, lastname, dob, email, phonenumber) VALUES (?, ?, ?, ?, ?)",
            newUser.FirstName, newUser.LastName, newUser.DOB, newUser.Email, newUser.PhoneNumber)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newUser)
    }
}

func GetUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse user ID from request parameters
        userID := r.URL.Query().Get("id")
        if userID == "" {
            http.Error(w, "Missing user ID parameter", http.StatusBadRequest)
            return
        }

        // Query the database to get the user by ID
        var user User
        err := db.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
            &user.ID, &user.FirstName, &user.LastName, &user.DOB, &user.Email, &user.PhoneNumber)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Encode user as JSON and send response
        json.NewEncoder(w).Encode(user)
    }
}

func GetAllUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Query the database to get all users
        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Create a slice to store all users
        var users []User

        // Iterate over the rows and scan each user into a struct
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DOB, &user.Email, &user.PhoneNumber); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            users = append(users, user)
        }

        // Check for errors during iteration
        if err := rows.Err(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Encode users as JSON and send response
        json.NewEncoder(w).Encode(users)
    }
}

func EditUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse user ID from request parameters
        userID := r.URL.Query().Get("id")
        if userID == "" {
            http.Error(w, "Missing user ID parameter", http.StatusBadRequest)
            return
        }

        // Decode request body into a User struct
        var updatedUser User
        err := json.NewDecoder(r.Body).Decode(&updatedUser)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Update the user in the database
        _, err = db.Exec("UPDATE users SET firstname=?, lastname=?, dob=?, email=?, phonenumber=? WHERE id=?",
            updatedUser.FirstName, updatedUser.LastName, updatedUser.DOB, updatedUser.Email, updatedUser.PhoneNumber, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Fetch the updated user from the database
        var editedUser User
        err = db.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
            &editedUser.ID, &editedUser.FirstName, &editedUser.LastName, &editedUser.DOB, &editedUser.Email, &editedUser.PhoneNumber)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Encode the updated user as JSON and send response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(editedUser)
    }
}