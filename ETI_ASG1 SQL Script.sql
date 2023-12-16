CREATE DATABASE IF NOT EXISTS carpooling_db;
USE carpooling_db;

-- Drop tables if they exist
DROP TABLE IF EXISTS user_trips;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    mobile_number VARCHAR(15) NOT NULL,
    email_address VARCHAR(100) NOT NULL,
    account_password VARCHAR(255) NOT NULL, -- Change the length to a suitable value
    is_car_owner BOOLEAN DEFAULT FALSE,
    driver_license_number VARCHAR(20),
    car_plate_number VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trips table
CREATE TABLE trips (
    tripId INT AUTO_INCREMENT PRIMARY KEY,
    ownerId INT,
    pickupLocation VARCHAR(255) NOT NULL,
    alternatePickupLocation VARCHAR(255),
    startTime TIME NOT NULL,
    destination VARCHAR(255) NOT NULL,
    seatsAvailable INT NOT NULL,
    published BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (ownerId) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Create user_trips table after users and trips
CREATE TABLE user_trips (
    user_id INT,
    tripId INT,
    PRIMARY KEY (user_id, tripId),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (tripId) REFERENCES trips(tripId)
);

-- Reset AUTO_INCREMENT values
ALTER TABLE users AUTO_INCREMENT = 1;
ALTER TABLE trips AUTO_INCREMENT = 1;

-- Create user if it does not exist
CREATE USER IF NOT EXISTS 'user'@'localhost' IDENTIFIED BY
'password';
GRANT ALL ON *.* TO 'user'@'localhost';

-- Insert sample data
INSERT INTO users (first_name, last_name, mobile_number, email_address, account_password, created_at)
VALUES ('John', 'Doe', '+123456789', 'john.doe@example.com', 'password', '2022-11-14 20:00:00');

INSERT INTO users (first_name, last_name, mobile_number, email_address, account_password, is_car_owner, driver_license_number, car_plate_number, created_at)
VALUES ('Jack', 'Doe', '34567', 'jack.doe@example.com', 'password', 1, '12345', '5151G', '2023-11-14 20:00:00');
