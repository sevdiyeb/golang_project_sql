package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Car struct {
	ID           int
	Manufacturer string
	Model        string
	Year         int
}

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTable(db)

	newCarID := insert(db, Car{
		Manufacturer: "Toyota",
		Model:        "Corolla",
		Year:         2022,
	})
	fmt.Println("Inserted car with ID:", newCarID)

	err = update(db, 1, Car{
		Manufacturer: "Honda",
		Model:        "Civic",
		Year:         2023,
	})
	if err != nil {
		log.Fatal(err)
	}

	data := selectData(db)
	fmt.Println("Cars:", data)
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS cars (
		id SERIAL PRIMARY KEY,
		manufacturer VARCHAR(50) NOT NULL,
		model VARCHAR(50) NOT NULL,
		year INT NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insert(db *sql.DB, car Car) int {
	query := `INSERT INTO cars (manufacturer, model, year)
	VALUES ($1, $2, $3) RETURNING id`

	var pk int
	err := db.QueryRow(query, car.Manufacturer, car.Model, car.Year).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}

func update(db *sql.DB, id int, newCar Car) error {
	query := `UPDATE cars SET manufacturer = $1, model = $2, year = $3 WHERE id = $4`

	_, err := db.Exec(query, newCar.Manufacturer, newCar.Model, newCar.Year, id)
	if err != nil {
		return err
	}
	return nil
}

func selectData(db *sql.DB) []Car {
	data := []Car{}
	rows, err := db.Query("SELECT id, manufacturer, model, year FROM cars")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id, year int
	var manufacturer, model string

	for rows.Next() {
		err := rows.Scan(&id, &manufacturer, &model, &year)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Car{id, manufacturer, model, year})
	}
	return data
}
