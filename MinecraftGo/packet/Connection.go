package packet

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"

	"../dataTypes"
)

type Connection struct {
	State State
	Conn  net.Conn
	Key   *rsa.PrivateKey
}

type State int

const (
	HANDSHAKING State = 0
	STATUS      State = 1
	LOGIN       State = 2
	PLAY        State = 3
)

var packets = map[State][]Packet{
	HANDSHAKING: {
		0x00: &Handshaking{},
	},
	STATUS: {
		0x00: &StatusRequest{},
		0x01: &StatusPing{},
	},
	LOGIN: {
		0x00: &LoginStart{},
		0x0B: &PluginMessage{IsServer: true},
		0x05: &ClientSettings{},
	},
	PLAY: {},
}

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":          dataTypes.ReadLong,
	"varInt":        dataTypes.ReadVarInt,
	"string":        dataTypes.ReadString,
	"raw":           dataTypes.ReadRaw,
	"unsignedShort": dataTypes.ReadUnsignedShort,
	"bool":          dataTypes.ReadBoolean,
	"unsignedByte":  dataTypes.ReadUnsignedByte,
	"int":           dataTypes.ReadInt,
	"float":         dataTypes.ReadFloat,
	"double":        dataTypes.ReadDouble,
}

var dataWriteMap = map[string]func(interface{}) []byte{
	"long":          dataTypes.WriteLong,
	"varInt":        dataTypes.WriteVarInt,
	"string":        dataTypes.WriteString,
	"raw":           dataTypes.WriteRaw,
	"unsignedShort": dataTypes.WriteUnsignedShort,
	"bool":          dataTypes.WriteBoolean,
	"unsignedByte":  dataTypes.WriteUnsignedByte,
	"int":           dataTypes.WriteInt,
	"intArray":      dataTypes.WriteIntArray,
	"float":         dataTypes.WriteFloat,
	"double":        dataTypes.WriteDouble,
}

func Init(conn net.Conn, key *rsa.PrivateKey) *Connection {
	fmt.Println("--New Connection!")
	newConnection := Connection{
		State: HANDSHAKING,
		Conn:  conn,
		Key:   key,
	}

	go newConnection.Handle()

	return &newConnection
}

func (c *Connection) Handle() {
	buf := make([]byte, 4096)
	for {
		// Read the incoming connection into the buffer.
		readLength, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			//_ = c.Conn.Close()
			return
		} else {
			iLength, end := dataTypes.ReadVarInt(buf)
			length := iLength.(int)
			cursor := end

			fmt.Printf("Read length: %d, Reported Length: %d\n", readLength, length)

			iPacketType, end := dataTypes.ReadVarInt(buf[end:])
			packetType := iPacketType.(int)
			cursor += end
			fmt.Printf("State %d, Packet Type: %d\n", c.State, packetType)

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

			fmt.Println(">>>INCOMING<<<")
			fmt.Println(hex.Dump(packetBuffer))

			c.StructScan(&packet, packetBuffer)
			packet.Handle(packetBuffer, c)
		}
	}
}

func (c *Connection) StructScan(packet *Packet, buf []byte) {
	v := reflect.ValueOf(*packet).Elem()
	t := reflect.TypeOf(*packet).Elem()

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

		if len(buf) < cursor {
			fmt.Println("Cursor overrun")
			continue
		}
		val, end := dataReadMap[tag](buf[cursor:])

		cursor += end
		//fmt.Printf("Reading tag %s into field %s value %v cursor %d\n", tag, field.Name, val, cursor)

		v.FieldByName(field.Name).Set(reflect.ValueOf(val))
	}
}

func (c *Connection) SendPacket(packet *Packet) {
	fmt.Println("Send packet time")
	v := reflect.ValueOf(*packet).Elem()
	t := reflect.TypeOf(*packet).Elem()

	payload := make([]byte, 0)

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}
		if dataWriteMap[tag] == nil {
			fmt.Println("Unknown tag type ", tag)
			continue
		}
		val := v.FieldByName(field.Name).Interface()
		if val == nil {
			fmt.Println("nil value!!!", field.Name)
			continue
		}

		segment := dataWriteMap[tag](v.FieldByName(field.Name).Interface())

		fmt.Printf("Field %s type %s coded as value %v\n", field.Name, tag, segment)

		payload = append(payload, segment...)
	}

	payload = append([]byte{byte((*packet).GetPacketId())}, payload...)
	actualLength := len(payload)
	payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	sendingLength, _ := dataTypes.ReadVarInt(payload)
	fmt.Println("Sending Length ", sendingLength)
	fmt.Println("Payload length ", actualLength)
	fmt.Println("Writing to connection...")
	_, exception := c.Conn.Write(payload)
	fmt.Println(">>>OUTGOING<<<")
	fmt.Println(hex.Dump(payload))

	if exception != nil {
		fmt.Println("Exception Writing ", exception)
	}

}
