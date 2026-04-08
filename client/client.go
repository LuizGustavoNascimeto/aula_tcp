package main

import (
	"bufio"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	originClient = "cliente"
	originServer = "servidor"
)

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
	appLog(originClient, "enviando: %s", summarizeCommand(c))
	_, err := conn.Write([]byte(c + "\n"))
	if err != nil {
		appLog(originClient, "erro ao enviar comando: %v", err)
		return
	}

	if isPWDCommand(message) {
		handlePWDResponse(conn, reader)
		return
	}

	resp, err := reader.ReadString('\n')
	if err != nil {
		appLog(originClient, "erro ao ler resposta: %v", err)
		return
	}

	appLog(originServer, strings.TrimSpace(resp))
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Comando: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func handleCommand(msg string) string {
	command := strings.Fields(msg)
	if len(command) >= 3 && command[0] == "CONNECT" {
		msg = commandAuth(command[1], command[2])
	}
	return msg
}

func commandAuth(login string, password string) string {
	pass := sha512.Sum512([]byte(password))
	return fmt.Sprintf("CONNECT %s %x", login, pass)
}

func isPWDCommand(message string) bool {
	command := strings.Fields(message)
	return len(command) > 0 && command[0] == "PWD"
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

func handlePWDResponse(conn net.Conn, reader *bufio.Reader) {
	ack, err := reader.ReadString('\n')
	if err != nil {
		appLog(originClient, "erro ao ler ACK do PWD: %v", err)
		return
	}
	appLog(originServer, strings.TrimSpace(ack))

	if strings.TrimSpace(ack) != "ACK: PWD" {
		return
	}

	dir, err := readPWDDir(conn, reader)
	if err != nil {
		appLog(originClient, "erro ao ler diretorio atual: %v", err)
		return
	}

	appLog(originServer, "Diretorio atual: %s", dir)
	_, _ = conn.Write([]byte("ACK: PWD_RECEIVED\n"))
	appLog(originClient, "confirmacao de PWD enviada")
}

func readPWDDir(conn net.Conn, reader *bufio.Reader) (string, error) {
	_ = conn.SetReadDeadline(time.Now().Add(700 * time.Millisecond))
	defer conn.SetReadDeadline(time.Time{})

	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return "", errors.New("timeout esperando diretorio do PWD")
		}
		if errors.Is(err, io.EOF) {
			return "", errors.New("conexao encerrada antes de receber diretorio")
		}
		return "", err
	}

	return strings.TrimSpace(string(buf[:n])), nil
}
