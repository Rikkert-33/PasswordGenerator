package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
)

const (
	dbname = "mijndb"
	dbuser = "rik"
	dbpass = "SQLR1k"
)

var database *sql.DB //database connection

func main() {
	//Connect to the database
	db, err := sql.Open("postgres", "dbname="+dbname+" user="+dbuser+" password="+dbpass+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	database = db
	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	generatePassword(10, true, true)

}

func createTable() error {
	query := `CREATE TABLE IF NOT EXISTS passwords (
		id			SERIAL	PRIMARY KEY,
		password	TEXT	NOT NULL
	);`
	_, err := database.Exec(query)
	return err
}

func generatePassword(length int, includeNumbers bool, includeSymbols bool) string {
	password := ""

	characters := ""
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	symbols := "!@#$%^&*()_+-=[]{};:,./<>?"

	//Numbers and or symbols are included if true
	characters += letters
	if includeNumbers {
		characters += numbers
	}
	if includeSymbols {
		characters += symbols
	}

	//Generates a random password
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			log.Fatal(err)
		}
		password += string(characters[index.Int64()])
	}
	fmt.Println(password) //used to test the function
	return password
}
