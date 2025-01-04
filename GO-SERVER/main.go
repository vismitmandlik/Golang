package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func formHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// Handle form submission here
		// You can read form data and process it as needed
		fmt.Fprintf(w, "Form submitted!")

		return
	}

	// Serve the HTML form for GET requests
	http.ServeFile(w, r, "./static/form.html")
}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)

	http.HandleFunc("/form", formHandler)

	fmt.Printf("Server started at 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {

		log.Fatalln(err)
	}
}
