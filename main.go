package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Pelico - Video Game Collection Manager")
	fmt.Println("Please use 'go run cmd/server/main.go' to start the server")
	fmt.Println("Or run the server binary directly: './pelico'")
	os.Exit(1)
}