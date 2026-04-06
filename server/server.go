package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	openConn()

}

func openConn() {
	fmt.Println("Porta aberta :8080")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Conexão encerrada:", err)
			return
		}
		fmt.Print("Recebido do cliente ", conn.RemoteAddr(), ": ", msg)
		// responde de volta
		conn.Write([]byte("ACK: " + msg))
	}
}
