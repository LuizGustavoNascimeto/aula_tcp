package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func Run(addr string) error {
	return NewClient(addr).Run()
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("erro ao conectar: %w", err)
	}
	defer conn.Close()
	defer ClientLog("conexao encerrada")

	ClientLog("conectado em %s", c.addr)
	return c.handleConnection(conn)
}

func (c *Client) handleConnection(conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("erro ao ler entrada: %w", err)
			}
			return nil
		}

		if err := c.handleCommand(conn, scanner.Text()); err != nil {
			ClientLog("erro: %v", err)
		}
	}
}
