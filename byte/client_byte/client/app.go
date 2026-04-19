package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Run(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("erro ao conectar: %w", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("conectado")
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("erro ao ler entrada: %w", err)
			}
			return nil
		}

		if err := processCommand(conn, scanner.Text()); err != nil {
			fmt.Printf("erro: %v\n", err)
		}
	}
}
