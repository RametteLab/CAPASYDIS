package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files (HTML, CSS, JS) from the "current" directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// Start the server on port 8080
	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
