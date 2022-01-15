package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ebobo/sqlite3_go/pkg/utility"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Hello Qi")
	os.Remove("sqlite-database.db")

	log.Println("Creating sqlite-database.db...")

	// make data dir if it is not exit
	err := utility.MakeDirIfNotExists("../data")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create SQLite file
	file, err := os.Create("../data/sqlite-database.db")

	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	log.Println("sqlite-database.db created")

	// Open the created SQLite File
	sqliteDatabase, _ := sql.Open("sqlite3", "../data/sqlite-database.db")

	// Defer Closing the database
	defer sqliteDatabase.Close()

	createTable(sqliteDatabase) // Create Database Tables

	// INSERT RECORDS
	insertCar(sqliteDatabase, "BB 00007", "Model X", "Tesla")
	insertCar(sqliteDatabase, "BB 00006", "MX5", "Mazda")
	insertCar(sqliteDatabase, "BB 00005", "X3", "BMW")
	insertCar(sqliteDatabase, "BB 00004", "MX5", "Mazda")
	insertCar(sqliteDatabase, "BB 00003", "Fortwo", "Smart")
	insertCar(sqliteDatabase, "BB 00002", "X1", "BMW")
	insertCar(sqliteDatabase, "BB 00001", "3 Series", "BMW")

	displayGarageCars(sqliteDatabase)
	log.Println("---------------------------------")
	displayGarageCarByBrand(sqliteDatabase, "BMW")
}

func createTable(db *sql.DB) {
	createGarageTableSQL := `CREATE TABLE garage (
		"idCar" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"license" TEXT,
		"model" TEXT,
		"brand" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create garage table...")
	statement, err := db.Prepare(createGarageTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("garage table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertCar(db *sql.DB, license string, model string, brand string) {
	log.Println("Inserting car info ...")
	insertGarageSQL := `INSERT INTO garage(license, model, brand) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertGarageSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(license, model, brand)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayGarageCars(db *sql.DB) {
	row, err := db.Query("SELECT * FROM garage ORDER BY brand")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var license string
		var model string
		var brand string
		row.Scan(&id, &license, &model, &brand)
		log.Println("Car: ", id, " ", license, " ", model, " ", brand)
	}
}

func displayGarageCarByBrand(db *sql.DB, brand string) {
	row, err := db.Query("SELECT * FROM garage WHERE brand = $1", brand)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var license string
		var model string
		var brand string
		row.Scan(&id, &license, &model, &brand)
		log.Println("Car: ", id, " ", license, " ", model, " ", brand)
	}
}
