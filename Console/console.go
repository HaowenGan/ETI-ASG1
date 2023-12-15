package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

	Loginmenu()
}

func Loginmenu() {
	var option int
	for {
		fmt.Println("===================\nCar Pooling Console\n 1. Login\n 2. Register\n 0. Quit")
		fmt.Print("Enter an option: ")
		fmt.Scan(&option)
		if option == 0 {
			fmt.Println("Thank you for using Car Pool!")
			break
		} else if option == 1 {
			// Call login function
			email, password := getUserCredentials()
			successfulLogin := Login(email, password)
			if successfulLogin {
				fmt.Println("Login successful!")
				UserMenu(email)

			} else {
				fmt.Println("Login failed. Invalid email or password.")
			}
		} else if option == 2 {
			registerUser()
		}
	}
}

func getUserCredentials() (string, string) {
	var email, password string
	fmt.Print("Enter email: ")
	fmt.Scan(&email)
	fmt.Print("Enter password: ")
	fmt.Scan(&password)
	return email, password
}

func Login(email, password string) bool {
	// Check if the email and password match a user in the database
	var storedPassword string
	err := db.QueryRow("SELECT account_password FROM users WHERE email_address=?", email).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		// User not found
		return false
	} else if err != nil {
		log.Println("Error querying user:", err)
		return false
	}

	// Check if the stored password matches the provided password
	return storedPassword == password
}

func registerUser() {
	// Collect user information
	var firstName, lastName, mobileNumber, emailAddress, password string

	fmt.Print("Enter first name: ")
	fmt.Scan(&firstName)
	fmt.Print("Enter last name: ")
	fmt.Scan(&lastName)
	fmt.Print("Enter mobile number: ")
	fmt.Scan(&mobileNumber)
	fmt.Print("Enter email address: ")
	fmt.Scan(&emailAddress)
	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	// You can add more input validation if needed

	// Call the Register function from the API
	err := Register(firstName, lastName, mobileNumber, emailAddress, password)
	if err != nil {
		fmt.Println("Error registering user:", err)
	} else {
		fmt.Println("User registered successfully!")
	}
}

func Register(firstName, lastName, mobileNumber, emailAddress, password string) error {
	// Construct the form data
	formData := url.Values{
		"first_name":       {firstName},
		"last_name":        {lastName},
		"mobile_number":    {mobileNumber},
		"email_address":    {emailAddress},
		"account_password": {password},
	}

	// Make the HTTP POST request
	resp, err := http.PostForm("http://localhost:5001/api/v1/register", formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status and handle it accordingly
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status code: %d", resp.StatusCode)
	}

	// If needed, you can read and process the response body here

	return nil
}

func UserMenu(email string) {
	var option int
	var password string
	var isCarOwner bool

	// Check if the user is a car owner
	err := db.QueryRow("SELECT is_car_owner FROM users WHERE email_address=?", email).Scan(&isCarOwner)
	if err != nil {
		fmt.Println("Error retrieving user information:", err)
		return
	}
	for {
		if isCarOwner {
			fmt.Println("============\nUser Console\n 1. Create Trip\n 2. Edit Trip\n 3. Update User Details\n 4. Delete Account\n 0. Logout")
			fmt.Print("Enter an option: ")
			fmt.Scan(&option)
			if option == 0 {
				break
			} else if option == 1 {
				CreateTrip(email)
			} else if option == 2 {
				EditTrip(email)
			} else if option == 3 {
				UpdateUserDetails(email, password)
			} else if option == 4 {
				fmt.Print("Enter your password to confirm account deletion: ")
				fmt.Scan(&password)
				err := DeleteUser(email, password)
				if err != nil {
					fmt.Println("Error deleting user:", err)
				} else {
					fmt.Println("Account deleted successfully!")
					break // exit the loop after deleting the account
				}
			}
		} else {
			fmt.Println("============\nUser Console\n 1. Enrol Trip\n 2. Update User Details\n 3. Delete Account\n 4. Past Trips\n 0. Logout")
			fmt.Print("Enter an option: ")
			fmt.Scan(&option)
			if option == 0 {
				break
			} else if option == 2 {
				UpdateUserDetails(email, password)
			} else if option == 3 {
				fmt.Print("Enter your password to confirm account deletion: ")
				fmt.Scan(&password)
				err := DeleteUser(email, password)
				if err != nil {
					fmt.Println("Error deleting user:", err)
				} else {
					fmt.Println("Account deleted successfully!")
					break // exit the loop after deleting the account
				}
			}
		}

	}
}

func UpdateUserDetails(email, password string) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email_address=?", email).Scan(&userID)
	if err != nil {
		fmt.Println("Error retrieving user ID:", err)
		return
	}
	// Retrieve updated user information
	var firstName, lastName, mobileNumber, newEmail, newPassword string
	var isCarOwner bool
	var driverLicenseNumber, carPlateNumber string

	fmt.Print("Enter new first name: ")
	fmt.Scan(&firstName)
	fmt.Print("Enter new last name: ")
	fmt.Scan(&lastName)
	fmt.Print("Enter new mobile number: ")
	fmt.Scan(&mobileNumber)
	fmt.Print("Enter new email address: ")
	fmt.Scan(&newEmail)
	fmt.Print("Enter new password: ")
	fmt.Scan(&newPassword)

	// Check if the user is a car owner and collect additional information
	fmt.Print("Are you a car owner? (true/false): ")
	fmt.Scan(&isCarOwner)
	if isCarOwner {
		fmt.Print("Enter driver license number: ")
		fmt.Scan(&driverLicenseNumber)
		fmt.Print("Enter car plate number: ")
		fmt.Scan(&carPlateNumber)
	}

	// Call the UpdateUser function from the API
	err = UpdateUser(userID, email, password, firstName, lastName, mobileNumber, newEmail, newPassword, isCarOwner, driverLicenseNumber, carPlateNumber)
	if err != nil {
		fmt.Println("Error updating user details:", err)
	} else {
		fmt.Println("User details updated successfully!")
	}
}

func UpdateUser(userID int, email, password, firstName, lastName, mobileNumber, newEmail, newPassword string, isCarOwner bool, driverLicenseNumber, carPlateNumber string) error {
	// Construct the form data
	formData := url.Values{
		"email":                 {email},
		"password":              {password},
		"first_name":            {firstName},
		"last_name":             {lastName},
		"mobile_number":         {mobileNumber},
		"email_address":         {newEmail},
		"account_password":      {newPassword},
		"is_car_owner":          {strconv.FormatBool(isCarOwner)},
		"driver_license_number": {driverLicenseNumber},
		"car_plate_number":      {carPlateNumber},
	}

	// Construct the URL with the user ID
	url := fmt.Sprintf("http://localhost:5001/api/v1/update/%d", userID)

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status and handle it accordingly
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update failed with status code: %d", resp.StatusCode)
	}

	// If needed, you can read and process the response body here

	return nil
}

func DeleteUser(email, password string) error {
	// Check if the email and password match a user in the database
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email_address=? AND account_password=?", email, password).Scan(&userID)
	if err == sql.ErrNoRows {
		// User not found
		return fmt.Errorf("user not found")
	} else if err != nil {
		return err
	}

	// Construct the URL with the user ID
	url := fmt.Sprintf("http://localhost:5001/api/v1/delete/%d", userID)

	// Create a new DELETE request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status and handle it accordingly
	if resp.StatusCode == http.StatusOK {
		fmt.Println("User deleted successfully!")
		return nil
	}

	// Read the response body for debugging purposes
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusForbidden {
		errorMessage := fmt.Sprintf("%s", string(body))
		return fmt.Errorf(errorMessage)
	}

	// Handle other unexpected errors
	errorMessage := fmt.Sprintf("Unexpected error occurred. Response: %s", string(body))
	fmt.Println("Error deleting user:", errorMessage)
	return fmt.Errorf(errorMessage)
}

func CreateTrip(email string) {
	var pickupLocation, alternatePickupLocation, startTimeStr, destination string
	var seatsAvailable int
	var published bool

	var ownerId int
	err := db.QueryRow("SELECT user_id FROM users WHERE email_address=?", email).Scan(&ownerId)
	if err != nil {
		fmt.Println("Error checking owner's email:", err)
		return
	}

	fmt.Print("Enter pickup location: ")
	fmt.Scan(&pickupLocation)
	fmt.Print("Enter alternate pickup location: ")
	fmt.Scan(&alternatePickupLocation)
	fmt.Print("Enter start time (HH:mm format): ")
	fmt.Scan(&startTimeStr)
	fmt.Print("Enter destination: ")
	fmt.Scan(&destination)
	fmt.Print("Enter seats available: ")
	fmt.Scan(&seatsAvailable)
	fmt.Print("Is the trip published? (true/false): ")
	fmt.Scan(&published)

	// Validate and parse the time string
	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return
	}

	// Construct the trip data
	tripData := Trip{
		OwnerEmail:              email,
		PickupLocation:          pickupLocation,
		AlternatePickupLocation: alternatePickupLocation,
		StartTime:               startTime.Format("15:04"),
		Destination:             destination,
		SeatsAvailable:          seatsAvailable,
		Published:               published,
	}

	// Convert the trip data to JSON
	tripJSON, err := json.Marshal(tripData)
	if err != nil {
		fmt.Println("Error marshaling trip data:", err)
		return
	}

	// Make the HTTP POST request to the trip API
	resp, err := http.Post("http://localhost:5000/api/v1/MakeTrip", "application/json", bytes.NewBuffer(tripJSON))
	if err != nil {
		fmt.Println("Error creating trip:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status and handle it accordingly
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Creating trip failed with status code:", resp.StatusCode)
		return
	}

	// If needed, you can read and process the response body here

	fmt.Println("Trip created successfully")
}

func EditTrip(email string) {
	var tripID int
	fmt.Print("Enter the trip ID to edit: ")
	fmt.Scan(&tripID)

	// Check if the trip belongs to the user
	var ownerID int
	err := db.QueryRow("SELECT ownerId FROM Trips WHERE tripId=?", tripID).Scan(&ownerID)
	if err != nil {
		fmt.Println("Error checking trip owner:", err)
		return
	}

	// Check if the trip belongs to the logged-in user
	var userID int
	err = db.QueryRow("SELECT user_id FROM users WHERE email_address=?", email).Scan(&userID)
	if err != nil {
		fmt.Println("Error checking user ID:", err)
		return
	}

	if ownerID != userID {
		fmt.Println("You do not have permission to edit this trip.")
		return
	}

	// Call the EditTrip function with the tripID
	EditTripDetails(tripID)
}

func EditTripDetails(tripID int) {
	// Fetch the existing trip details
	var existingTrip Trip
	err := db.QueryRow("SELECT tripId, ownerId, pickupLocation, alternatePickupLocation, startTime, destination, seatsAvailable, published FROM Trips WHERE tripId=?", tripID).Scan(
		&existingTrip.TripID,
		&existingTrip.OwnerID,
		&existingTrip.PickupLocation,
		&existingTrip.AlternatePickupLocation,
		&existingTrip.StartTime,
		&existingTrip.Destination,
		&existingTrip.SeatsAvailable,
		&existingTrip.Published,
	)
	if err != nil {
		fmt.Println("Error fetching trip details:", err)
		return
	}

	// Display existing details and prompt for new details
	fmt.Println("Existing Trip Details:")
	fmt.Printf("%+v\n", existingTrip)

	var updatedTrip Trip
	fmt.Print("Enter new pickup location (press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.PickupLocation)
	if updatedTrip.PickupLocation == "" {
		updatedTrip.PickupLocation = existingTrip.PickupLocation
	}

	fmt.Print("Enter new alternate pickup location (press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.AlternatePickupLocation)
	if updatedTrip.AlternatePickupLocation == "" {
		updatedTrip.AlternatePickupLocation = existingTrip.AlternatePickupLocation
	}

	fmt.Print("Enter new Start time (press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.StartTime)
	if updatedTrip.StartTime == "" {
		updatedTrip.StartTime = existingTrip.StartTime
	}

	fmt.Print("Enter new Destination (press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.Destination)
	if updatedTrip.Destination == "" {
		updatedTrip.Destination = existingTrip.Destination
	}

	fmt.Print("Enter new Seats Available (press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.SeatsAvailable)
	if updatedTrip.SeatsAvailable == 0 {
		updatedTrip.SeatsAvailable = existingTrip.SeatsAvailable
	}

	fmt.Print("Enter new Published value (true/false, press Enter to keep existing): ")
	fmt.Scan(&updatedTrip.Published)
	if !updatedTrip.Published {
		updatedTrip.Published = existingTrip.Published
	}

	// Add similar prompts for other fields you want to update

	// Call the API to update the trip details
	err = UpdateTripDetails(tripID, updatedTrip)
	if err != nil {
		fmt.Println("Error updating trip details:", err)
		return
	}

	fmt.Println("Trip details updated successfully!")
}

func UpdateTripDetails(tripID int, updatedTrip Trip) error {
	// Convert the updated trip data to JSON
	tripJSON, err := json.Marshal(updatedTrip)
	if err != nil {
		return fmt.Errorf("Error marshaling updated trip data: %v", err)
	}

	// Construct the URL with the trip ID
	url := fmt.Sprintf("http://localhost:5000/api/v1/EditTrip/%d", tripID)

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(tripJSON))
	if err != nil {
		return fmt.Errorf("Error creating PUT request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error updating trip: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status and handle it accordingly
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Updating trip failed with status code: %d", resp.StatusCode)
	}

	// If needed, you can read and process the response body here

	return nil
}
