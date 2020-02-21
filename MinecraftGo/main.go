package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"./entity"
	"./packet"
)

var keyPair *rsa.PrivateKey

var minecraft entity.Minecraft

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

	purple := entity.Purple
	minecraft = entity.Minecraft{
		Connections: make([]*packet.Connection, 0),
		ServerName: entity.ChatMessageComponent{
			Text:   "Petecraft",
			Colour: &purple,
		},
		MaxPlayers:       0,
		EnableEncryption: false,
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
	minecraft.Connections = append(minecraft.Connections, packet.Init(connection, keyPair, minecraft))
}
