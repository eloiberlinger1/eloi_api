package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {

	dataBytes, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Impossible to reqd data.json: %v", err)
	}

	var staticData StaticData
	err = json.Unmarshal(dataBytes, &staticData)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	cfg := config{
		addr: ":8080",
	}

	app := application{
		config: cfg,
		data:   staticData,
	}

	h := app.mount()

	err = app.run(h)
	if err != nil {
		log.Println("Server has failed to start")
		os.Exit(1)
	}
}
