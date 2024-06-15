package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
)

func Main(args map[string]interface{}) map[string]interface{} {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("Error: DIGITALOCEAN_ACCESS_TOKEN environment variable is not set.")
		os.Exit(1)
	}

	client := godo.NewFromToken(token)
	ctx := context.TODO()

	droplets, _, err := client.Droplets.List(ctx, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, droplet := range droplets {
		fmt.Printf("Droplet ID: %d, Name: %s, Status: %s\n", droplet.ID, droplet.Name, droplet.Status)
	}
	msg := make(map[string]interface{})
	msg["body"] = "Hello World!"
	return msg
}
