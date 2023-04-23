package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Reading and waiting")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Signalling ready...")
		w.WriteHeader(http.StatusNoContent)
	}))
}
