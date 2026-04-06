package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn := connect()
	for {
		text := readTerminal()
		fmt.Println("Texto digitado:", text)
		sendMessage(conn, text)
	}
}

func connect() net.Conn {
	fmt.Println("Conectando...")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Conectado!")
	return conn
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message))

	resp, _ := bufio.NewReader(conn).ReadString('\n')
	if resp != "ACK: "+message {
		fmt.Println("Resposta inesperada do servidor:", resp)
	}
}

func readTerminal() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	return text
}
