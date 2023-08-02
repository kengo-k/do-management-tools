package main

import (
	"encoding/json"
	lib "helloworld/lib/example"
	"log"
)

func main() {
	jsonStr := `{"Name":"John"}`

	var person lib.Person
	err := json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	lib.Hello(person)
}
