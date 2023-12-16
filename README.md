# ETI-ASG1
## Overview
- This carpooling system consists of two microservices: User API and Trip API, along with a console application that interacts with these APIs.

## Instructions for Setting up and running of Microservices
- Run the "ETI_ASG1 SQL Script" in MySQL Workbench to setup database
- go run user.go and trip.go to start the API
- go run console.go to start the console

### User API Endpoints:

- `/api/v1/register`: Allows users to register by providing their personal information.
- `/api/v1/delete/{userID}`: Deletes a user account by specifying the user ID.
- `/api/v1/update/{userID}`: Updates user details by specifying the user ID.

### Security:

- The system enforces a one-year minimum account age before allowing deletion to prevent accidental deletions.

### Validation:

- User registration checks for duplicate email addresses and mobile numbers to ensure uniqueness.
- During user updates, the system validates that the required information is provided and ensures that car-related details are not blank when the user is a car owner.

### Trip API Endpoints:

- `/api/v1/MakeTrip`: Creates a new trip by providing trip details.
- `/api/v1/EditTrip/{tripId}`: Updates an existing trip by specifying the trip ID.
- `/api/v1/ViewTrips`: Lists all trips that are published so that user can enrol
- `/api/v1/ViewPastTrips/{user_id}`: List all trips that user have taken that are past the current time

### Security:

- The API checks if the owner's email exists in the Users table before allowing trip creation or update.
- Trip updates are restricted to the owner of the trip to prevent unauthorized modifications.

### Validation:

- Trip creation validates the owner's email and ensures it exists in the Users table before proceeding.

### Console Application User Interaction:

- The console application provides a menu-driven interface for users to log in, register, create/edit trips, and update their user details.

### User Authentication:

- Users must log in with a valid email and password, and the system checks the provided credentials against the User API for authentication.

### Trip Creation/Editing:

- The console allows users to create new trips or edit existing trips by interacting with the Trip API.

### User Details Update/Deletion:

- Users can update their details and delete their accounts through the console, which interfaces with the User API.

## Implementation Details

This microservice is designed to handle two main functionalities: trip management and user management for a carpooling application. The microservice is implemented in Go, utilizing the Gorilla Mux router for handling HTTP requests and a MySQL database for data storage.

### Database Schema

#### Trips Table:

- `tripId` (Primary Key)
- `ownerId` (Foreign Key referencing Users table)
- `pickupLocation`
- `alternatePickupLocation`
- `startTime`
- `destination`
- `seatsAvailable`
- `published`

#### Users Table:

- `user_id` (Primary Key)
- `first_name`
- `last_name`
- `mobile_number`
- `email_address`
- `account_password`
- `is_car_owner`
- `driver_license_number` (Nullable)
- `car_plate_number` (Nullable)
- `created_at`

#### User_Trips Table:

- `user_id` (Foreign Key referencing Users table)
- `tripId` (Foreign Key referencing Trips table)
