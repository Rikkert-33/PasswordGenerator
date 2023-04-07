package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	_ "github.com/go-sql-driver/mysql" //Used to initialize the mysql driver package, so you don't need to call it directly in your code.
	_ "github.com/lib/pq"
)

/*
ideas for future updates:
- user to choose length of password
- user to choose if they want numbers and/or symbols
- option to ask for stored password to show out of the database
- option to ask for password to delete from database
- encrypt/decrypt passwords
- db passwords connect to user?
*/

type Config struct {
	DBname string `json:"dbname"`
	DBuser string `json:"dbuser"`
	DBpass string `json:"dbpass"`
	DBhost string `json:"dbhost"`
	DBport string `json:"dbport"`
}

var database *sql.DB //database connection
var logger *log.Logger

// add log file
func init() {
	// Open the log file
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new logger that writes to the log file
	logger = log.New(logFile, "", log.Ldate|log.Ltime)
}

func main() {
	//Read the config file
	file, err := os.Open("config.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		logger.Fatal(err)
	}

	//Connect to the database
	dsn := fmt.Sprintf("%s:%s@/%s", config.DBuser, config.DBpass, config.DBname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err)
	}

	database = db
	if err := createTable(); err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	password := generatePassword(10, true, true)

	//Checks if password already exists and creates a new password if it already exists
	for checkPassword(db, password) {
		password = generatePassword(10, true, true)
	}

	addPassword(db, password)

	fmt.Println(password)
}

func createTable() error {
	createTableSQL := `
CREATE TABLE IF NOT EXISTS passwords (
  id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  password VARCHAR(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	var err error
	_, err = database.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
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
			logger.Fatal(err)
		}
		password += string(characters[index.Int64()])
	}
	logger.Println("Generated password:", password)
	return password
}

func checkPassword(db *sql.DB, password string) bool {
	//Checks if the password is in the database
	query := "SELECT COUNT(*) FROM passwords WHERE password = ?"
	var count int
	err := db.QueryRow(query, password).Scan(&count)
	if err != nil {
		logger.Fatal(err)
	}
	//Password already exists if count > 0
	return count > 0
}

func addPassword(db *sql.DB, password string) {
	//Adds the password to the database
	query := "INSERT INTO passwords(password) VALUES(?)"
	_, err := db.Exec(query, password)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Password added to the database")
}
