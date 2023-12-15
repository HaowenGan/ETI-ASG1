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

