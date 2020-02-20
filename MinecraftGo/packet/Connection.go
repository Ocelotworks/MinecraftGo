package packet

import (
	"fmt"
	"net"
	"reflect"

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

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":   dataTypes.ReadLong,
	"varInt": dataTypes.ReadVarInt,
}

var dataWriteMap = map[string]func(interface{}) []byte{
	"long":   dataTypes.WriteLong,
	"varInt": dataTypes.WriteVarInt,
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
			iLength, end := dataTypes.ReadVarInt(buf)
			length := iLength.(int)
			cursor := end

			fmt.Printf("Read length: %d, Reported Length: %d\n", readLength, length)

			iPacketType, end := dataTypes.ReadVarInt(buf[end:])
			packetType := iPacketType.(int)
			cursor += end
			fmt.Printf("Packet Type: %d\n", packetType)

			if packets[c.State] == nil {
				fmt.Println("Bad State ", c.State)
				continue
			}

			if packets[c.State][packetType] == nil {
				fmt.Println("Bad Packet Type ", packetType)
				continue
			}

			packet := packets[c.State][packetType]

			packetBuffer := buf[cursor : cursor+length]
			c.StructScan(&packet, packetBuffer)
			packet.Handle(packetBuffer, c)
		}
	}
}

func (c *Connection) StructScan(packet *Packet, buf []byte) {
	v := reflect.ValueOf(packet).Elem()
	t := reflect.TypeOf(packet).Elem()

	cursor := 0

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}
		if dataReadMap[tag] == nil {
			fmt.Println("Unknown tag type ", tag)
			continue
		}

		val, end := dataReadMap[tag](buf[cursor:])

		cursor += end

		v.FieldByName(field.Name).Set(reflect.ValueOf(val))
	}
}

func (c *Connection) SendPacket(packet *Packet) {
	v := reflect.ValueOf(packet).Elem()
	t := reflect.TypeOf(packet).Elem()

	payload := make([]byte, 0)

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}
		if dataReadMap[tag] == nil {
			fmt.Println("Unknown tag type ", tag)
			continue
		}

		segment := dataWriteMap[tag](v.FieldByName(field.Name).Elem())
		payload = append(payload, segment...)
	}

	payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	payload = append([]byte{byte((*packet).GetPacketId())}, payload...)

	_, exception := c.Conn.Write(payload)

	if exception != nil {
		fmt.Println("Exception Writing ", exception)
	}

}
