package main

import (
	"client_byte/client"
	"fmt"
)

func main() {
	if err := client.Run("127.0.0.1:8080"); err != nil {
		fmt.Printf("erro: %v\n", err)
	}
}
