package main

import (
	"html/template"
	"log"
	"net/http"
)

// Data struct to pass data to the template
type Data struct {
	HiddenValue string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Define your data
		data := Data{
			HiddenValue: "some value",
		}

		// Parse your HTML template
		tmpl, err := template.ParseFiles("template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}

		// Execute the template and pass data to it
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
