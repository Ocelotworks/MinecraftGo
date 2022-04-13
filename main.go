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

	//inData, exception := ioutil.ReadFile("data/worlds/world/region/r.0.0.mca") //ioutil.ReadFile("C:\\Users\\Peter\\AppData\\Roaming\\.minecraft\\saves\\MCGO Flat Test 2\\region\\r.0.0.mca")
	//
	//if exception != nil {
	//	fmt.Println("Reading file")
	//	fmt.Println(exception)
	//	return
	//}
	//
	//region := dataTypes.ReadRegionFile(inData, nil)
	//
	////fmt.Println(level)
	//
	////palette := level["Sections"].(map[string]interface{})["List-0"].(map[string]interface{})["Compound_1"].(map[string]interface{})["Palette"].(map[string]interface{})
	//
	////byte((len(castBlockStates)*64)/4096),
	//
	//for _, chunk := range region.Chunks {
	//	if chunk == nil {
	//		continue
	//	}
	//	//fmt.Println("Chunk",chunk)
	//	chunkRaw := dataTypes.WriteChunk(chunk.Sections)
	//	fmt.Println("Chunk Raw", chunkRaw)
	//}
	//
	//return

	listener, exception := net.Listen("tcp", ":25566")

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

	fmt.Println("Listening on port 25566")

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
