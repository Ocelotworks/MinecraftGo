package packet

import (
	"fmt"
	"net"

	"../dataTypes"
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

var packets = map[State][]Packet{
	HANDSHAKING: {0x00: &Handshaking{}, 0xFE: nil /*Legacy type*/},
	STATUS:      {0x00: &StatusRequest{}, 0x01: &StatusPing{}},
}

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

			if packets[c.State] == nil {
				fmt.Println("Bad State ", c.State)
			} else if packets[c.State][packetType] == nil {
				fmt.Println("Bad Packet Type ", packetType)
			} else {
				packets[c.State][packetType].Handle(buf[cursor:cursor+length], c)
			}
		}
	}
}
