package main

import (
	_ "embed" // support embedding files in variables
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ebobo/sqlite3_go/pkg/model"
	"github.com/ebobo/sqlite3_go/pkg/utility"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

var mutex sync.RWMutex

// var schema = `
// DROP TABLE IF EXISTS lego;
// CREATE TABLE lego (
//     name text,
//     model integer,
//     catalog text
// );`

//go:embed schema.sql
var schema string

func main() {
	log.Println("Hello Qi, Sqlx")

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
	sqliteDatabase, _ := sqlx.Open("sqlite3", "../data/sqlite-database.db")

	// Defer Closing the database
	defer sqliteDatabase.Close()

	createSchema(sqliteDatabase) // Create Database Tables

	// INSERT RECORDS
	addLegoSet(sqliteDatabase, &model.LegoSet{Name: "Police Station", Model: 10278, Catalog: "Creator"})
	addLegoSet(sqliteDatabase, &model.LegoSet{Name: "Volkswagen T1 Camper Van", Model: 10220, Catalog: "Creator"})
	addLegoSet(sqliteDatabase, &model.LegoSet{Name: "Folkevognboble", Model: 10252, Catalog: "Creator"})
	addLegoSet(sqliteDatabase, &model.LegoSet{Name: "NASA Apollo Saturn V", Model: 21309, Catalog: "Ideas"})
	addLegoSet(sqliteDatabase, &model.LegoSet{Name: "Lamborghini", Model: 42115, Catalog: "Technic"})

	displayLegoSets(sqliteDatabase)

	log.Println("-------------Get set 42115----------------")

	set, err := getLegoSet(sqliteDatabase, 42115)

	if err != nil {
		log.Fatalln(err)
	}

	js, err := json.Marshal(set)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(js))

	// log.Println("-------------Get creator sets----------------")

	// sets, err := displayLegoSetsByCatalog(sqliteDatabase, "Creator")

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, set := range sets {
	// 	js, err := json.Marshal(set)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(string(js))
	// }

	log.Println("-------------Get creator sets----------------")

	func() {
		set, err := displayLegoSetsByCatalog2(sqliteDatabase, "Creator")
		if err != nil {
			log.Fatalln(err)
		}

		js, err := json.Marshal(set)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(js))
	}()

}

func createSchema(db *sqlx.DB) error {
	for n, statement := range strings.Split(schema, ";") {
		_, err := db.Exec(statement)
		if err != nil {
			return fmt.Errorf("statement %d failed: \"%s\" : %w", n+1, statement, err)
		}
	}
	return nil
}

func addLegoSet(db *sqlx.DB, set *model.LegoSet) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, err := db.NamedExec(
		`INSERT INTO lego (name, model, catalog) 
         VALUES (:name, :model, :catalog)`,
		set,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func getLegoSet(db *sqlx.DB, modelNum int) (*model.LegoSet, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var set model.LegoSet
	return &set, db.QueryRowx("SELECT * FROM lego WHERE model = ?", modelNum).StructScan(&set)
}

// func deleteSet(db *sqlx.DB, modelNum int) error {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	_, err := db.Exec("DELETE FROM lego WHERE id = ?", modelNum)
// 	return err
// }

// func displayLegoSetsByCatalog(db *sqlx.DB, catalog string) ([]model.LegoSet, error) {
// 	mutex.RLock()
// 	defer mutex.RUnlock()

// 	var sets []model.LegoSet
// 	return sets, db.Select(&sets, "SELECT * FROM lego WHERE catalog = ? ORDER BY name", catalog)
// }

func displayLegoSetsByCatalog2(db *sqlx.DB, catalog string) (*model.LegoSet, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var sets []model.LegoSet
	err := db.Select(&sets, "SELECT * FROM lego WHERE catalog = ? limit 1", catalog)

	if err == nil {
		log.Println("length of slice ", len(sets))
	}

	return &sets[0], err
}

func displayLegoSets(db *sqlx.DB) {
	row, err := db.Query("SELECT * FROM lego ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var name string
		var model int
		var catalog string
		row.Scan(&model, &catalog)
		log.Println("Lego: ", name, " ", model, " ", catalog)
	}
}

// func displayLegoSetsNameAndModelNumber(db *sqlx.DB) {
// 	row, err := db.Query("SELECT name, model FROM lego ORDER BY name")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer row.Close()
// 	for row.Next() { // Iterate and fetch the records from result cursor
// 		var name string
// 		var model int
// 		row.Scan(&name, &model)
// 		log.Println("Lego: ", name, " ", model)
// 	}
// }
