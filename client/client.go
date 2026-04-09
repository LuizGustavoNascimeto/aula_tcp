package main

import (
	"bufio"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	originClient = "cliente"
	originServer = "servidor"
	terminal     = ">>>"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mu     sync.Mutex
}

var stdoutMu sync.Mutex

func main() {
	client := connect()
	defer client.conn.Close()
	appLog(originClient, "pronto para enviar comandos")
	go client.receiveResponses()
	for {
		text := readTerminal()
		client.sendMessage(text)
	}
}

func appLog(origin string, format string, args ...interface{}) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func connect() *Client {
	appLog(originClient, "conectando ao servidor em :8080")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		appLog(originClient, "erro ao conectar: %v", err)
		os.Exit(1)
	}
	appLog(originClient, "conexao estabelecida")
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *Client) sendMessage(message string) {
	command := handleCommand(message)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.writer.WriteString(command + "\n"); err != nil {
		appLog(originClient, "erro ao enviar comando: %v", err)
		return
	}
	if err := c.writer.Flush(); err != nil {
		appLog(originClient, "erro ao enviar comando: %v", err)
		return
	}
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)

	printPrompt()
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func handleCommand(msg string) string {
	/*Isso só existe para poder fazer a hash*/
	command := strings.Fields(msg)
	if len(command) == 0 {
		return msg
	}

	if command[0] == "CONNECT" {
		payload := strings.TrimSpace(strings.TrimPrefix(msg, "CONNECT"))
		parts := strings.SplitN(payload, ",", 2)

		if len(parts) == 2 {
			login := strings.TrimSpace(parts[0])
			password := strings.TrimSpace(parts[1])
			if login != "" && password != "" {
				msg = commandAuth(login, password)
			}
		} else if len(command) >= 3 {
			login := strings.TrimSuffix(command[1], ",")
			msg = commandAuth(login, command[2])
		}
	}
	if command[0] == "EXIT" {
		appLog(originClient, "encerrando conexao")
		os.Exit(0)
	}
	return msg
}

/*Isso tbm só existe para poder fazer a hash*/
func commandAuth(login string, password string) string {
	pass := sha512.Sum512([]byte(password))
	return fmt.Sprintf("CONNECT %s,%x", login, pass)
}

func summarizeCommand(message string) string {
	command := strings.Fields(message)
	if len(command) == 0 {
		return "comando vazio"
	}
	if command[0] == "CONNECT" && len(command) >= 3 {
		return fmt.Sprintf("CONNECT %s [hash]", command[1])
	}
	return message
}

func (c *Client) receiveResponses() {
	for {
		dir, err := c.readResponse()
		if err != nil {
			appLog(originClient, "erro ao ler resposta: %v", err)
			return
		}
		appLog(originServer, "%s", strings.TrimSpace(dir))
		printPrompt()
	}
}

func (c *Client) readResponse() (string, error) {
	return c.reader.ReadString('\n')
}

func printPrompt() {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()

	fmt.Print(terminal + " ")
}
