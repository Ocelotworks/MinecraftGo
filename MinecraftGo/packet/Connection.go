package packet

import (
	"../dataTypes"
	"fmt"
	"net"
)

type Connection struct {
	State State
	Conn  net.Conn
}

type State int

const (
	HANDSHAKING State = 0
	STATUS      State = 1
	LOGIN       State = 2
	PLAY        State = 3
)

func Init(conn net.Conn) *Connection {
	newConnection := Connection{
		State: HANDSHAKING,
		Conn:  conn,
	}

	go newConnection.Handle()

	return &newConnection
}

func (c *Connection) Handle() {
	buf := make([]byte, 1024)
	for {
		// Read the incoming connection into the buffer.
		readLength, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		} else {
			length, end := dataTypes.ReadVarInt(buf)
			cursor := end

			fmt.Printf("Read length: %d, Reported Length: %d\n", readLength, length)

			packetType, end := dataTypes.ReadVarInt(buf[end:])
			cursor += end
			fmt.Printf("Packet Type: %d\n", packetType)

			var packet Packet
			switch packetType {
			case 0:
				packet = &Handshaking{}
				break
			}

			packet.Handle(buf[cursor:cursor+length], c)
		}
	}
}
