package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: bcryp-password <password> <cost>\n")
		os.Exit(0)
	}

	cost, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Cost must be an integer: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), cost)
	end := time.Now()
	fmt.Printf("Time taken:%dms", end.Sub(start).Milliseconds())

	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\t%s\n", string(hash))
}
