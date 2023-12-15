# ETI-ASG1
Overview
This carpooling system consists of two microservices: User API and Trip API, along with a console application that interacts with these APIs. Below are the design considerations for each component:

User API
Endpoints:

/api/v1/register: Allows users to register by providing their personal information.
/api/v1/delete/{userID}: Deletes a user account by specifying the user ID.
/api/v1/update/{userID}: Updates user details by specifying the user ID.
Security:

Passwords are stored securely using a hash-based encryption method to enhance security.
The system enforces a one-year minimum account age before allowing deletion to prevent accidental deletions.
Validation:

User registration checks for duplicate email addresses and mobile numbers to ensure uniqueness.
During user updates, the system validates that the required information is provided and ensures that car-related details are not blank when the user is a car owner.
Trip API
Endpoints:

/api/v1/MakeTrip: Creates a new trip by providing trip details.
/api/v1/EditTrip/{tripId}: Updates an existing trip by specifying the trip ID.
Security:

The API checks if the owner's email exists in the Users table before allowing trip creation or update.
Trip updates are restricted to the owner of the trip to prevent unauthorized modifications.
Validation:

Trip creation validates the owner's email and ensures it exists in the Users table before proceeding.
Updates to trips perform necessary validations, and further validations can be added based on specific business rules.
Console Application
User Interaction:

The console application provides a menu-driven interface for users to log in, register, create/edit trips, and update their user details.
User Authentication:

Users must log in with a valid email and password, and the system checks the provided credentials against the User API for authentication.
Trip Creation/Editing:

The console allows users to create new trips or edit existing trips by interacting with the Trip API.
User Details Update/Deletion:

Users can update their details and delete their accounts through the console, which interfaces with the User API.

This microservice is designed to handle two main functionalities: trip management and user management for a carpooling application. The microservice is implemented in Go, utilizing the Gorilla Mux router for handling HTTP requests and a MySQL database for data storage.

Trip API
Endpoints
Create Trip: /api/v1/MakeTrip (POST)

Creates a new carpooling trip based on the provided JSON payload.
Validates the owner's email by checking against the Users table.
Inserts trip details into the Trips table in the database.
Update Trip: /api/v1/EditTrip/{tripId} (PUT)

Updates an existing carpooling trip based on the provided trip ID and JSON payload.
Performs necessary validation on the updated trip data.
Updates trip details in the Trips table in the database.
View Published Trips: /api/v1/ViewTrips (GET)

Retrieves a list of published trips from the Trips table in the database.
Returns the list of trips as JSON.
View Past Trips: /api/v1/ViewPastTrips/{user_id} (GET)

Retrieves a list of past trips for a specific user from the Trips table.
Joins the user_trips table to identify trips associated with the user.
Filters trips based on the current time compared to the trip start time.
Returns the list of past trips as JSON.
Database Schema
Trips Table:

tripId (Primary Key)
ownerId (Foreign Key referencing Users table)
pickupLocation
alternatePickupLocation
startTime
destination
seatsAvailable
published
Users Table:

user_id (Primary Key)
first_name
last_name
mobile_number
email_address
account_password
is_car_owner
driver_license_number (Nullable)
car_plate_number (Nullable)
created_at
User_Trips Table:

user_id (Foreign Key referencing Users table)
tripId (Foreign Key referencing Trips table)
User API
Endpoints
Register User: /api/v1/register (POST)

Registers a new user based on form data.
Checks for duplicate email addresses and mobile numbers before registration.
Inserts user details into the Users table in the database.
Delete User: /api/v1/delete/{userID} (DELETE)

Deletes a user based on the provided user ID.
Checks if the user exists and if the account is more than 1 year old.
Deletes the user account from the Users table in the database.
Update User: /api/v1/update/{userID} (PUT)

Updates user details based on the provided user ID and form data.
Validates and updates user information in the Users table in the database.
Supports additional fields for car owners, including driver's license and car plate number.
