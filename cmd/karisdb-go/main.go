package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"karisdb/internal/app"
)

func main() {
	command := "help"
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: karisdb-go [help|version|status]")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		command = os.Args[1]
	}

	router, err := app.NewRouter()
	if err != nil {
		fmt.Fprintln(os.Stderr, "karisdb-go: core could not initialize")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, appErr := router.Execute(ctx, command, nil)
	if appErr != nil {
		fmt.Fprintln(os.Stderr, appErr.Message)
		os.Exit(1)
	}

	switch command {
	case "help":
		fmt.Println("Commands: help, version, status")
	case "version":
		fmt.Printf("KarisDB %s\n", app.Version)
	default:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			fmt.Fprintln(os.Stderr, "karisdb-go: could not format response")
			os.Exit(2)
		}
	}
}
