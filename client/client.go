package main

import (
	"bufio"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	originClient = "cliente"
	originServer = "servidor"
)

var validCommands = map[string]struct{}{
	"CONNECT": {},
	"PWD":     {},
	"CHDIR":   {},
}

func appLog(origin string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[%s][%s] %s\n", timestamp, origin, msg)
}

func main() {
	conn := connect()
	defer conn.Close()
	reader := bufio.NewReader(conn)
	appLog(originClient, "pronto para enviar comandos")
	for {
		text := readTerminal()
		sendMessage(conn, reader, text)
	}
}

func connect() net.Conn {
	appLog(originClient, "conectando ao servidor em :8080")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		appLog(originClient, "erro ao conectar: %v", err)
		os.Exit(1)
	}
	appLog(originClient, "conexao estabelecida")
	return conn
}

func sendMessage(conn net.Conn, reader *bufio.Reader, message string) {
	c := handleCommand(message)
	_, err := conn.Write([]byte(c + "\n"))
	if err != nil {
		appLog(originClient, "erro ao enviar comando: %v", err)
		return
	}

	if hasResponse(message) {
		handleResponse(conn, reader)
		return
	}
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

/*Isso só existe para poder fazer a hash*/
func handleCommand(msg string) string {
	command := strings.Fields(msg)
	if len(command) >= 3 && isCommand(command[0]) && command[0] == "CONNECT" {
		msg = commandAuth(command[1], command[2])
	}
	return msg
}

/*Isso tbm só existe para poder fazer a hash*/
func commandAuth(login string, password string) string {
	pass := sha512.Sum512([]byte(password))
	return fmt.Sprintf("CONNECT %s %x", login, pass)
}

func hasResponse(message string) bool {
	command := strings.Fields(message)
	return len(command) > 0 && isCommand(command[0])
}

func isCommand(value string) bool {
	_, exists := validCommands[value]
	return exists
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

func handleResponse(conn net.Conn, reader *bufio.Reader) {
	dir, err := reader.ReadString('\n')
	if err != nil {
		appLog(originClient, "erro ao ler respota: %v", err)
		return
	}
	appLog(originServer, "%s", dir)
}
