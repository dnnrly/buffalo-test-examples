package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: bcryp-check <password> <hash>\n")
		os.Exit(0)
	}

	start := time.Now()

	err := bcrypt.CompareHashAndPassword([]byte(os.Args[2]), []byte(os.Args[1]))
	end := time.Now()
	fmt.Printf("Time taken:%dms", end.Sub(start).Milliseconds())

	if err != nil {
		fmt.Printf("Error checking password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\tPassword correct\n")
}
