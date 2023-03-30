package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DBname string `json:"dbname"`
	DBuser string `json:"dbuser"`
	DBpass string `json:"dbpass"`
}

var database *sql.DB //database connection

func main() {
	//Read the config file
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	//Connect to the database
	db, err := sql.Open("postgres", "dbname="+config.DBname+" user="+config.DBuser+" password="+config.DBpass+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	database = db
	if err := createTable(); err != nil {
		log.Fatal(err)
	}

	password := generatePassword(10, true, true)

	//Checks if password already exists and creates a new password if it already exists
	for checkPassword(db, password) {
		password = generatePassword(10, true, true)
	}

	addPassword(db, password)

	fmt.Println(password)
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

func checkPassword(db *sql.DB, password string) bool {
	//Checks if the password is in the database
	query := "SELECT COUNT(*) FROM passwords WHERE password = ?"
	var count int
	err := db.QueryRow(query, password).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	//Password already exists if count > 0
	return count > 0
}

func addPassword(db *sql.DB, password string) {
	//Adds the password to the database
	query := "INSERT INTO passwords(password) VALUES(?)"
	_, err := db.Exec(query, password)
	if err != nil {
		log.Fatal(err)
	}
}
