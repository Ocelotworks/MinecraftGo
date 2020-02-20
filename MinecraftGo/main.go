package main

import (
	"./dataTypes"
	"./packet"
	"fmt"
	"net"
)

func main() {

	test, end := dataTypes.ReadVarInt([]byte{0xc2, 0x04})
	fmt.Printf("Test 0 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0x01}, 0)
	//fmt.Printf("Test 1 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0x02}, 0)
	//fmt.Printf("Test 2 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0x7f}, 0)
	//fmt.Printf("Test 127 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0x80, 0x01}, 0)
	//fmt.Printf("Test 128 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0xff, 0x01}, 0)
	//fmt.Printf("Test 255 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x07}, 0)
	//fmt.Printf("Test 2147483647 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x0f}, 0)
	//fmt.Printf("Test -1 = %d (end %d)\n", test, end)
	//
	//test, end = dataTypes.ReadVarInt([]byte{0x80, 0x80, 0x80, 0x80, 0x08}, 0)
	//fmt.Printf("Test -2147483648 = %d (end %d)\n", test, end)

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
