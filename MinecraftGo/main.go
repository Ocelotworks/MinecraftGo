package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"./packet"
)

var keyPair *rsa.PrivateKey

func main() {

	listener, exception := net.Listen("tcp", "localhost:25565")

	if exception != nil {
		fmt.Println(exception)
		return
	}

	fmt.Println("Generating Keypair")
	key, exception := rsa.GenerateKey(rand.Reader, 1024)

	if exception != nil {
		fmt.Println(exception)
		return
	}

	keyPair = key

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
	packet.Init(connection, keyPair)
}
