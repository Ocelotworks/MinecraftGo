package main

import (
	"fmt"
	"net"

	"./packet"
)

func main() {

	listener, exception := net.Listen("tcp", "localhost:25565")

	if exception != nil {
		fmt.Println(exception)
		return
	}

	fmt.Println("Listening on port 25565")

	for {
		connection, exception := listener.Accept()

		if exception != nil {
			fmt.Println(exception)
		} else {
			go handleConnection(connection)
		}
	}
}

func handleConnection(connection net.Conn) {
	packet.Init(connection)
}
