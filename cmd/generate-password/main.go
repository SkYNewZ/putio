package main

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	log.SetFlags(0)
}

func main() {
	// Read password
	plain, err := getPasswordFromInput()
	if err != nil {
		log.Fatalln(err)
	}

	// Hash it
	hash, err := hashPassword(plain)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("\nHash: %s", hash)

}

func getPasswordFromInput() ([]byte, error) {
	// Prompt the user to enter a password
	fmt.Print("Enter a password: ")

	pwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}

	return pwd, nil
}

func hashPassword(password []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return string(bytes), err
}
