package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func main() {
	generatePassword(10, true, true)

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
