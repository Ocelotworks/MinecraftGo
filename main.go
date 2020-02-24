package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/Ocelotworks/MinecraftGo/packet"
)

var keyPair *rsa.PrivateKey

var minecraft *packet.Minecraft

func main() {

	//inData, exception := ioutil.ReadFile("world/region/r.0.0.mca")
	//
	//if exception != nil {
	//	fmt.Println("Reading file")
	//	fmt.Println(exception)
	//	return
	//}
	//
	//region := dataTypes.ReadRegionFile(inData)
	//
	//chunk := region.Chunks[0]
	//
	//output, _ := json.Marshal(dataTypes.NBTAsMap(chunk.Data))
	//fmt.Println("Chunky chunk")
	//fmt.Println(string(output))
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
