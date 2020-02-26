package controller

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Connection struct {
	State             State
	Conn              net.Conn
	Key               *rsa.PrivateKey
	Minecraft         *Minecraft
	Player            *entity.Player
	KeepAliveID       int64
	EnableCompression bool
	EnableEncryption  bool

	VerifyToken      []byte
	SharedSecret     []byte
	Cipher           cipher.Stream
	DecryptionCipher cipher.Stream
}

type State int

const (
	HANDSHAKING State = 0
	STATUS      State = 1
	LOGIN       State = 2
	PLAY        State = 3
)

var controllers = map[State][]Packet{
	HANDSHAKING: {
		0x00: &Handshaking{},
	},
	STATUS: {
		0x00: &StatusRequest{},
		0x01: &StatusPing{},
	},
	LOGIN: {
		0x00: &packetType.LoginStart{},
		0x01: &packetType.EncryptionResponse{},
	},
	PLAY: {
		0x00: &packetType.TeleportConfirm{},
		0x03: &packetType.IncomingChatMessage{},
		0x05: &packetType.ClientSettings{},
		0x0B: &packetType.PluginMessage{IsServer: true},
		0x0F: &packetType.KeepAlive{},
		0x11: &packetType.PlayerPosition{},
		0x12: &packetType.PlayerPositionAndRotation{},
		0x13: &packetType.PlayerRotation{},
		0x14: &packetType.PlayerMovement{},
		0x1A: &packetType.PlayerDigging{},
		0x1B: &packetType.EntityAction{},
		0x23: &packetType.HeldItemChange{IsServer: true},
		0x2A: &packetType.Animation{},
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
	//"intArray":		 dataTypes.ReadIntArray,
	"float":           dataTypes.ReadFloat,
	"double":          dataTypes.ReadDouble,
	"uuid":            dataTypes.ReadUUID,
	"varIntByteArray": dataTypes.ReadVarIntByteArray,
}

var dataWriteMap = map[string]func(interface{}) []byte{
	"long":           dataTypes.WriteLong,
	"varInt":         dataTypes.WriteVarInt,
	"string":         dataTypes.WriteString,
	"raw":            dataTypes.WriteRaw,
	"short":          dataTypes.WriteShort,
	"unsignedShort":  dataTypes.WriteUnsignedShort,
	"bool":           dataTypes.WriteBoolean,
	"unsignedByte":   dataTypes.WriteUnsignedByte,
	"int":            dataTypes.WriteInt,
	"intArray":       dataTypes.WriteIntArray,
	"float":          dataTypes.WriteFloat,
	"double":         dataTypes.WriteDouble,
	"uuid":           dataTypes.WriteUUID,
	"entityMetadata": dataTypes.WriteEntityMetadata,
	"varIntArray":    dataTypes.WriteVarIntArray,
}

func Init(conn net.Conn, key *rsa.PrivateKey, minecraft *Minecraft) *Connection {
	fmt.Println("--New Connection!")
	newConnection := Connection{
		State:             HANDSHAKING,
		Conn:              conn,
		Key:               key,
		Minecraft:         minecraft,
		KeepAliveID:       0,
		EnableCompression: false,
		EnableEncryption:  false,
	}

	go newConnection.Handle()

	return &newConnection
}

func (c *Connection) sendKeepAlive() {
	for {
		//fmt.Println("Waiting")
		<-time.After(15 * time.Second)
		fmt.Println("Sending keepalive")
		keepAlive := packetType.Packet(&packetType.KeepAlive{
			ID: c.KeepAliveID,
		})

		c.KeepAliveID++

		exception := c.SendPacket(&keepAlive)
		if exception != nil {
			fmt.Println("client has probably gone away")

			go c.Minecraft.PlayerLeave(c)
			break
		}
	}
}

func (c *Connection) Handle() {
	buf := make([]byte, 4096)
	for {
		// Read the incoming connection into the buffer.
		count, err := c.Conn.Read(buf)

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
			decryptedBuf := make([]byte, count)

			if c.EnableEncryption {
				c.DecryptionCipher.XORKeyStream(decryptedBuf, buf[:count])
			} else {
				copy(decryptedBuf, buf[:count])
			}

			_, cursor := dataTypes.ReadVarInt(decryptedBuf)
			//length := iLength.(int)

			//fmt.Println("Packet length ", length)

			if c.EnableCompression {
				dataLength, end := dataTypes.ReadVarInt(decryptedBuf[cursor:])
				//fmt.Println("Data length ", dataLength)
				cursor += end
				if dataLength.(int) > 0 {
					//fmt.Println("Data length is above compression threshold")
					compressedPacket := bytes.NewReader(decryptedBuf[cursor:])
					r, exception := zlib.NewReader(compressedPacket)

					if exception != nil {
						fmt.Println("Exception decompressing packet ", exception)
						return
					}

					var uncompressedPacket bytes.Buffer
					_, exception = io.Copy(&uncompressedPacket, r)
					r.Close()
					if exception != nil {
						fmt.Println("Exception reading decompressed packet ", exception)
						return
					}

					decryptedBuf = uncompressedPacket.Bytes()
					cursor = 0
				}
			}

			//fmt.Printf("Read length: %d, Reported Length: %d\n", readLength, length)

			iPacketType, end := dataTypes.ReadVarInt(decryptedBuf[cursor:])
			packetType := iPacketType.(int)
			cursor += end

			if packets[c.State] == nil {
				fmt.Println("!!! Bad State ", c.State)
				continue
			}

			if packetType < 0 || len(packets[c.State]) < packetType || packets[c.State][packetType] == nil {
				fmt.Printf("!!! Bad Packet Type 0x%X\n", packetType)
				continue
			}

			packet := packets[c.State][packetType]

			//if cursor+length > len(decryptedBuf) {
			//	fmt.Println("Buffer overflow", cursor+length, len(decryptedBuf))
			//	continue
			//}

			packetBuffer := decryptedBuf[cursor:]

			//fmt.Println(">>>INCOMING<<<")
			//fmt.Println(hex.Dump(packetBuffer))

			c.StructScan(&packet, packetBuffer)
			//packet.Handle(packetBuffer, c)
		}
	}
}

func (c *Connection) StructScan(packet *packetType.Packet, buf []byte) {
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
			fmt.Println("!!! Unknown tag type ", tag)
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
			playerData := val.([]entity.Player)
			segment = dataTypes.WriteVarInt(len(playerData))
			//fmt.Println("Player Data Array Length ", len(playerData))
			for _, player := range playerData {
				segment = append(segment, UnmarshalData(&player)...)
			}
			//fmt.Println(hex.Dump(segment))
		} else if tag == "playerPropertiesArray" {
			playerProperties := val.([]entity.PlayerProperty)
			//fmt.Println("Player Property Array Length ", len(playerProperties))
			segment = dataTypes.WriteVarInt(len(playerProperties))
			if len(playerProperties) > 0 {
				for _, playerProperty := range playerProperties {
					if playerProperty.Signature != "" {
						playerProperty.Signed = true
					}
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

func (c *Connection) SendPacket(packet *packetType.Packet) error {
	var payload []byte
	packetID := byte((*packet).GetPacketId())
	fmt.Printf("Sending packet 0x%X\n", packetID)

	if c.EnableCompression {
		uncompressedPayload := append([]byte{packetID}, UnmarshalData(*packet)...)
		dataLength := len(uncompressedPayload)
		if dataLength >= c.Minecraft.CompressionThreshold {
			payload = dataTypes.WriteVarInt(dataLength) //Length of uncompressed payload - needs packet length before it

			var compressedPayload bytes.Buffer
			w := zlib.NewWriter(&compressedPayload)
			w.Write(uncompressedPayload)
			w.Close()

			payload = append(payload, compressedPayload.Bytes()...)
			payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
		} else {
			//Packet is below compression threshold so dataLength is 0
			uncompressedPayload = append(dataTypes.WriteVarInt(0), uncompressedPayload...)
			uncompressedPayload = append(dataTypes.WriteVarInt(len(uncompressedPayload)), uncompressedPayload...)
			payload = uncompressedPayload
		}
	} else {
		payload = UnmarshalData(*packet)
		payload = append([]byte{packetID}, payload...)
		payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	}

	//if len(payload) < 1024 {
	//	fmt.Println(">>>OUTGOING<<<")
	//	fmt.Println(hex.Dump(payload))
	//}

	if packetID == 0x34 {
		fmt.Println(payload)
	}

	fmt.Println("Writing payload")
	var exception error
	if c.EnableEncryption {
		//fmt.Println("Packet is encrypted")
		encryptedPayload := make([]byte, len(payload))
		c.Cipher.XORKeyStream(encryptedPayload, payload)
		_, exception = c.Conn.Write(encryptedPayload)
		//fmt.Println(hex.Dump(encryptedPayload))
	} else {
		_, exception = c.Conn.Write(payload)
	}

	//_, _ = dataTypes.ReadVarInt(payload)
	//fmt.Println("Sending Length ", sendingLength)
	//fmt.Println("Payload length ", actualLength)
	//fmt.Println("Writing to connection...")

	if exception != nil {
		//fmt.Println("Exception Writing ", exception)
		return exception
	}
	return nil
}
