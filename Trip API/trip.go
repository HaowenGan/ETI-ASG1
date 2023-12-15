package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Trip struct {
	TripID                  int    `json:"tripId"`
	OwnerID                 int    `json:"ownerId"`
	OwnerEmail              string `json:"ownerEmail"`
	PickupLocation          string `json:"pickupLocation"`
	AlternatePickupLocation string `json:"alternatePickupLocation"`
	StartTime               string `json:"startTime"`
	Destination             string `json:"destination"`
	SeatsAvailable          int    `json:"seatsAvailable"`
	Published               bool   `json:"published"`
}

func main() {
	var err error
	// Assign to the global db variable, not creating a new one
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpooling_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	// database operation
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/MakeTrip", CreateTrip).Methods("POST")
	router.HandleFunc("/api/v1/EditTrip/{tripId}", UpdateTrip).Methods("PUT")

	fmt.Println("Listening at Port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func CreateTrip(w http.ResponseWriter, r *http.Request) {
	var trip Trip
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&trip); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Check if the owner's email exists in the Users table
	var ownerId int
	err := db.QueryRow("SELECT user_id FROM users WHERE email_address=?", trip.OwnerEmail).Scan(&ownerId)
	if err != nil {
		http.Error(w, "Owner's email not found", http.StatusBadRequest)
		return
	}

	// Continue with trip creation using the valid ownerId

	// Insert the trip data into the database with the valid ownerId
	result, err := db.Exec("INSERT INTO Trips (ownerId, pickupLocation, alternatePickupLocation, startTime, destination, seatsAvailable, published) VALUES (?, ?, ?, ?, ?, ?, ?)",
		ownerId, trip.PickupLocation, trip.AlternatePickupLocation, trip.StartTime, trip.Destination, trip.SeatsAvailable, trip.Published)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating trip: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the auto-generated trip ID
	tripID, _ := result.LastInsertId()

	// Set the trip ID in the response
	trip.TripID = int(tripID)

	// Return the created trip as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func UpdateTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID, err := strconv.Atoi(vars["tripId"])
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	var updatedTrip Trip
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&updatedTrip); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Perform validation on the updated trip data if necessary

	// Update the trip data in the database
	_, err = db.Exec("UPDATE Trips SET pickupLocation=?, alternatePickupLocation=?, startTime=?, destination=?, seatsAvailable=?, published=? WHERE tripId=?",
		updatedTrip.PickupLocation, updatedTrip.AlternatePickupLocation, updatedTrip.StartTime, updatedTrip.Destination, updatedTrip.SeatsAvailable, updatedTrip.Published, tripID)
	if err != nil {
		log.Printf("Error updating trip: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the updated trip as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTrip)
}
