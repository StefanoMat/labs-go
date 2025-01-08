package main

import (
	"fmt"
	"net/http"

	"github.com/stefanomat/weather/weather"
)

func main() {
	http.HandleFunc("/weather", weather.ZipCodeHandler)
	fmt.Println("Server is listening on port 8080...")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		fmt.Errorf("Error starting server: %v", error)
	}
}
