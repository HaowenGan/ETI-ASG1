package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	var err error
	// Assign to the global db variable, not creating a new one
	db, err = sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/carpooling_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	// database operation
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/register", Register).Methods("POST")
	router.HandleFunc("/api/v1/delete/{userID:[0-9]+}", DeleteUser).Methods("DELETE")
	fmt.Println("Listening at Port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve user information from the form
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	mobileNumber := r.FormValue("mobile_number")
	emailAddress := r.FormValue("email_address")

	// Insert user information into the database
	_, err = db.Exec("INSERT INTO users (first_name, last_name, mobile_number, email_address) VALUES (?, ?, ?, ?)",
		firstName, lastName, mobileNumber, emailAddress)
	if err != nil {
		log.Println("Error inserting user into the database:", err)
		http.Error(w, "Error inserting user into the database", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	fmt.Fprintf(w, "User registered successfully!")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request URL
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Check if the user exists and get the creation date
	var createdAtStr string
	err := db.QueryRow("SELECT created_at FROM users WHERE user_id=?", userID).Scan(&createdAtStr)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error querying user:", err)
		http.Error(w, "Error querying user", http.StatusInternalServerError)
		return
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		log.Println("Error parsing created_at:", err)
		http.Error(w, "Error parsing created_at", http.StatusInternalServerError)
		return
	}

	// Calculate the duration between the creation date and the current date
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if createdAt.After(oneYearAgo) {
		http.Error(w, "User account is not more than 1 year old", http.StatusBadRequest)
		return
	}

	// Delete the user account
	_, err = db.Exec("DELETE FROM users WHERE user_id=?", userID)
	if err != nil {
		log.Println("Error deleting user:", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User deleted successfully!")
}
