package packet

import (
	"crypto/rsa"
	"fmt"
	"net"
	"reflect"
	"time"

	"../dataTypes"
	"../entity"
)

type Connection struct {
	State       State
	Conn        net.Conn
	Key         *rsa.PrivateKey
	Minecraft   *Minecraft
	Player      *entity.Player
	KeepAliveID int64
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
	},
	PLAY: {
		0x00: &TeleportConfirm{},
		0x05: &ClientSettings{},
		0x0B: &PluginMessage{IsServer: true},
		0x0F: &KeepAlive{},
		0x11: &PlayerPosition{},
		0x12: &PlayerPositionAndRotation{},
		0x13: &PlayerRotation{},
		0x14: &PlayerMovement{},
		0x1A: &PlayerDigging{},
		0x23: &HeldItemChange{IsServer: true},
		0x2A: &Animation{},
	},
}

var dataReadMap = map[string]func(buf []byte) (interface{}, int){
	"long":          dataTypes.ReadLong,
	"varInt":        dataTypes.ReadVarInt,
	"string":        dataTypes.ReadString,
	"raw":           dataTypes.ReadRaw,
	"short":         dataTypes.ReadShort,
	"unsignedShort": dataTypes.ReadUnsignedShort,
	"bool":          dataTypes.ReadBoolean,
	"unsignedByte":  dataTypes.ReadUnsignedByte,
	"int":           dataTypes.ReadInt,
	"float":         dataTypes.ReadFloat,
	"double":        dataTypes.ReadDouble,
	"uuid":          dataTypes.ReadUUID,
}

var dataWriteMap = map[string]func(interface{}) []byte{
	"long":          dataTypes.WriteLong,
	"varInt":        dataTypes.WriteVarInt,
	"string":        dataTypes.WriteString,
	"raw":           dataTypes.WriteRaw,
	"short":         dataTypes.WriteShort,
	"unsignedShort": dataTypes.WriteUnsignedShort,
	"bool":          dataTypes.WriteBoolean,
	"unsignedByte":  dataTypes.WriteUnsignedByte,
	"int":           dataTypes.WriteInt,
	"intArray":      dataTypes.WriteIntArray,
	"float":         dataTypes.WriteFloat,
	"double":        dataTypes.WriteDouble,
	"uuid":          dataTypes.WriteUUID,
}

func Init(conn net.Conn, key *rsa.PrivateKey, minecraft *Minecraft) *Connection {
	fmt.Println("--New Connection!")
	newConnection := Connection{
		State:       HANDSHAKING,
		Conn:        conn,
		Key:         key,
		Minecraft:   minecraft,
		KeepAliveID: 0,
	}

	go newConnection.Handle()

	return &newConnection
}

func (c *Connection) sendKeepAlive() {
	for {
		//fmt.Println("Waiting")
		<-time.After(15 * time.Second)
		fmt.Println("Sending keepalive")
		keepAlive := Packet(&KeepAlive{
			ID: c.KeepAliveID,
		})

		c.KeepAliveID++

		exception := c.SendPacket(&keepAlive)
		if exception != nil {
			fmt.Println("client has probably gone away")
			break
		}
	}
}

func (c *Connection) Handle() {
	buf := make([]byte, 4096)
	for {
		// Read the incoming connection into the buffer.
		_, err := c.Conn.Read(buf)
		if err != nil {
			//fmt.Println("Error reading:", err.Error())
			_ = c.Conn.Close()
			if c.Player != nil {
				c.Minecraft.ConnectedPlayers--
			}
			for i, conn := range c.Minecraft.Connections {
				if conn == c {
					c.Minecraft.Connections[i] = c.Minecraft.Connections[len(c.Minecraft.Connections)-1]
					c.Minecraft.Connections = c.Minecraft.Connections[:len(c.Minecraft.Connections)-1]
					break
				}
			}
			return
		} else {
			iLength, end := dataTypes.ReadVarInt(buf)
			length := iLength.(int)
			cursor := end

			//fmt.Printf("Read length: %d, Reported Length: %d\n", readLength, length)

			iPacketType, end := dataTypes.ReadVarInt(buf[end:])
			packetType := iPacketType.(int)
			cursor += end
			fmt.Printf("State %d, Packet Type: %d\n", c.State, packetType)

			if packets[c.State] == nil {
				fmt.Println("!!! Bad State ", c.State)
				continue
			}

			if len(packets[c.State]) < packetType || packets[c.State][packetType] == nil {
				fmt.Printf("!!! Bad Packet Type 0x%X\n", packetType)
				continue
			}

			packet := packets[c.State][packetType]

			packetBuffer := buf[cursor : cursor+length]

			//fmt.Println(">>>INCOMING<<<")
			//fmt.Println(hex.Dump(packetBuffer))

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
			//fmt.Println("!!! Unknown tag type ", tag)
			continue
		}

		if len(buf) < cursor {
			//fmt.Println("Cursor overrun")
			continue
		}
		val, end := dataReadMap[tag](buf[cursor:])

		cursor += end
		//fmt.Printf("Reading tag %s into field %s value %v cursor %d\n", tag, field.Name, val, cursor)

		v.FieldByName(field.Name).Set(reflect.ValueOf(val))
	}
}

func UnmarshalData(input interface{}) []byte {
	v := reflect.ValueOf(input).Elem()
	t := reflect.TypeOf(input).Elem()

	payload := make([]byte, 0)

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("proto")
		if !exists {
			continue
		}

		val := v.FieldByName(field.Name).Interface()
		if val == nil {
			fmt.Println("nil value!!!", field.Name)
			continue
		}

		var segment []byte

		if tag == "playerArray" {
			playerData := val.([]Player)
			segment = dataTypes.WriteVarInt(len(playerData))
			fmt.Println("Player Data Array Length ", len(playerData))
			for _, player := range playerData {
				segment = append(segment, UnmarshalData(&player)...)
			}
		} else if tag == "playerPropertiesArray" {
			playerProperties := val.([]PlayerProperty)
			fmt.Println("Player Property Array Length ", len(playerProperties))
			segment = dataTypes.WriteVarInt(len(playerProperties))
			if len(playerProperties) > 0 {
				for _, playerProperty := range playerProperties {
					segment = append(segment, UnmarshalData(&playerProperty)...)
				}
			}
		} else {
			if dataWriteMap[tag] == nil {
				fmt.Println("!!!! Unknown tag type ", tag)
				continue
			}

			segment = dataWriteMap[tag](v.FieldByName(field.Name).Interface())
		}
		//if len(segment) < 100 {
		//	fmt.Printf("Field %s type %s coded as value %v (Between  %d - %d)\n", field.Name, tag, val, len(payload), len(payload)+len(segment))
		//}else{
		//	fmt.Printf("Field %s type %s coded as value [Big value] (Between  %d - %d)\n", field.Name, tag, len(payload), len(payload)+len(segment))
		//}

		payload = append(payload, segment...)
	}

	return payload
}

func (c *Connection) SendPacket(packet *Packet) error {
	fmt.Println("Send packet time")
	payload := UnmarshalData(*packet)

	payload = append([]byte{byte((*packet).GetPacketId())}, payload...)
	_ = len(payload)
	payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	_, _ = dataTypes.ReadVarInt(payload)
	//fmt.Println("Sending Length ", sendingLength)
	//fmt.Println("Payload length ", actualLength)
	//fmt.Println("Writing to connection...")
	_, exception := c.Conn.Write(payload)
	//if len(payload) < 1024 {
	//	fmt.Println(">>>OUTGOING<<<")
	//	fmt.Println(hex.Dump(payload))
	//}

	if exception != nil {
		//fmt.Println("Exception Writing ", exception)
		return exception
	}
	return nil
}
