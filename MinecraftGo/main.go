package main

import (
	"./packet"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
)

var keyPair *rsa.PrivateKey

var minecraft *packet.Minecraft

func main() {

	//inData, exception := ioutil.ReadFile("nbt-test/bigtest.nbt")
	//
	//if exception != nil {
	//	fmt.Println("Reading file")
	//	fmt.Println(exception)
	//	return
	//}
	//
	//read, _ := dataTypes.ReadNBT(inData)
	//write := dataTypes.NBTWriteCompound(read)
	//
	//readAgain, _ := dataTypes.ReadNBT(write)
	//fmt.Println("Original: ",read)
	//fmt.Println("New:", readAgain)
	//
	//return

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

	minecraft = packet.CreateMinecraft()

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
