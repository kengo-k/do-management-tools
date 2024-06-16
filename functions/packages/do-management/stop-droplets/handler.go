package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"embed"

	"github.com/digitalocean/godo"
)

type DropletConfig struct {
	Project      string   `json:"project"`
	DropletNames []string `json:"droplet-names"`
}

type Config struct {
	DropletsList []DropletConfig `json:"droplets-list"`
}

//go:embed config.json
var configFile embed.FS

func Main(args map[string]interface{}) map[string]interface{} {

	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("Error: DIGITALOCEAN_ACCESS_TOKEN environment variable is not set.")
		os.Exit(1)
	}

	file, err := configFile.Open("config.json")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error: reading config file: %v", err),
		}
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error: parsing config file: %v", err),
		}
	}

	fmt.Printf("Config: %v\n", config)

	client := godo.NewFromToken(token)
	ctx := context.Background()

	projects, _, err := client.Projects.List(ctx, nil)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Error: list projects: %v", err),
		}
	}
	for _, p := range projects {
		fmt.Printf("project name: %v\n", p.Name)
		for _, dropletConfig := range config.DropletsList {
			if p.Name == dropletConfig.Project {
				resources, _, err := client.Projects.ListResources(ctx, p.ID, nil)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("Error: list resources for project `%v`: %v", p.Name, err),
					}
				}
				for _, resource := range resources {
					fmt.Printf("Resource URN: %v\n", resource.URN)
					dropletIdStr, err := getDropletID(resource.URN)
					if err != nil {
						fmt.Printf("Resource `%v` is not a droplet, skipped\n", resource.URN)
					}
					dropletId, err := strconv.Atoi(dropletIdStr)
					if err != nil {
						return map[string]interface{}{
							"error": fmt.Sprintf("Error: invalid dropletId `%v`: %v", dropletIdStr, err),
						}
					}
					_, _, err = client.DropletActions.PowerOff(ctx, dropletId)
					if err != nil {
						return map[string]interface{}{
							"error": fmt.Sprintf("Error: shutdown droplet `%v`: %v", dropletId, err),
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"message": "success",
	}
}

func getDropletID(urn string) (string, error) {
	prefix := "do:droplet:"
	if !strings.HasPrefix(urn, prefix) {
		return "", fmt.Errorf("not a Droplet URN: %s", urn)
	}
	dropletID := strings.TrimPrefix(urn, prefix)

	return dropletID, nil
}
