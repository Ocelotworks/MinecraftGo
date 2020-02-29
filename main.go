package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/Ocelotworks/MinecraftGo/controller"
)

var keyPair *rsa.PrivateKey

var minecraft *controller.Minecraft

func main() {

	listener, exception := net.Listen("tcp", ":25565")

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

	minecraft = controller.CreateMinecraft()

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
	minecraft.Connections = append(minecraft.Connections, controller.Init(connection, keyPair, minecraft))
}
