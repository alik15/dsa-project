package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open or create SQLite database file
	db, err := sql.Open("sqlite3", "./flights.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure that flights and seats tables exist, create if not
	if err := createFlightsTable(db); err != nil {
		log.Fatal(err)
	}
	if err := createSeatsTable(db); err != nil {
		log.Fatal(err)
	}

	// Example usage: insert a new flight
	flight := Flight{FlightNo: "ABC123", Departure: "New York", Destination: "Los Angeles"}
	if err := createFlight(db, flight); err != nil {
		log.Fatal(err)
	}

	// Book a seat on the flight
	if err := bookSeat(db, 1, 3); err != nil {
		log.Fatal(err)
	}

	// Check status of a seat
	status, err := checkSeatStatus(db, 1, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Seat 1 status:", status)
}

// Flight struct
type Flight struct {
	ID          int
	FlightNo    string
	Departure   string
	Destination string
}

// Function to create the flights table if it doesn't exist
func createFlightsTable(db *sql.DB) error {
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS flights (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            flight_no TEXT UNIQUE NOT NULL,
            departure TEXT NOT NULL,
            destination TEXT NOT NULL
        );
    `

	_, err := db.Exec(createTableSQL)
	return err
}

// Function to create the seats table if it doesn't exist
func createSeatsTable(db *sql.DB) error {
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS seats (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            flight_id INTEGER NOT NULL,
            seat_no INTEGER NOT NULL,
            booked BOOLEAN NOT NULL DEFAULT 0,
            FOREIGN KEY (flight_id) REFERENCES flights(id),
            CHECK (seat_no <= 500) -- Ensure seat number doesn't exceed 500
        );
    `

	_, err := db.Exec(createTableSQL)
	return err
}

// Function to create a new flight or update an existing one
func createFlight(db *sql.DB, flight Flight) error {
	// Check if the flight already exists
	var existingID int
	err := db.QueryRow("SELECT id FROM flights WHERE flight_no = ?", flight.FlightNo).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingID != 0 {
		// Flight already exists, update it
		_, err = db.Exec("UPDATE flights SET departure = ?, destination = ? WHERE id = ?", flight.Departure, flight.Destination, existingID)
		return err
	}

	// Flight doesn't exist, insert it
	_, err = db.Exec("INSERT INTO flights (flight_no, departure, destination) VALUES (?, ?, ?)",
		flight.FlightNo, flight.Departure, flight.Destination)
	return err
}

func bookSeat(db *sql.DB, flightID, seatNo int) error {
	// Check if the seat limit has been reached
	var bookedSeats int
	err := db.QueryRow("SELECT COUNT(*) FROM seats WHERE flight_id = ? AND booked = 1", flightID).Scan(&bookedSeats)
	if err != nil {
		return err
	}
	if bookedSeats >= 500 {
		return errors.New("seat limit reached for this flight")
	}

	// Check if the seat is already booked
	var booked bool
	err = db.QueryRow("SELECT booked FROM seats WHERE flight_id = ? AND seat_no = ? AND booked = 1", flightID, seatNo).Scan(&booked)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if booked {
		return fmt.Errorf("seat %d is already booked", seatNo)
	}

	// Update seat status to booked
	_, err = db.Exec("INSERT INTO seats (flight_id, seat_no, booked) VALUES (?, ?, 1)", flightID, seatNo)
	return err
}

// Function to check the status of a seat (if it is already booked or not)
func checkSeatStatus(db *sql.DB, flightID, seatNo int) (bool, error) {
	// Retrieve seat status
	var booked bool
	err := db.QueryRow("SELECT booked FROM seats WHERE flight_id = ? AND seat_no = ?", flightID, seatNo).Scan(&booked)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("seat %d is not available", seatNo)
		}
		return false, err
	}
	return booked, nil
}
