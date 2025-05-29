package controller

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"io"
	"net"
	"runtime/debug"
	"time"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"github.com/Ocelotworks/MinecraftGo/helpers/structScanner"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Connection struct {
	State             State
	Conn              net.Conn
	Key               *rsa.PrivateKey
	Minecraft         *Minecraft
	Player            *entity.Player
	Ping              int
	EnableCompression bool
	EnableEncryption  bool
	Joined            bool

	VerifyToken      []byte
	SharedSecret     []byte
	Cipher           cipher.Stream
	DecryptionCipher cipher.Stream
}

type State int

const (
	HANDSHAKING   State = 0
	STATUS        State = 1
	LOGIN         State = 2
	TRANSFER      State = 3
	CONFIGURATION State = 4
	PLAY          State = 5
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
		constants.SBLoginStart:              &LoginStart{},
		constants.SBLoginEncryptionResponse: &EncryptionResponse{},
		constants.SBLoginPluginResponse:     &LoginPluginResponse{},
		constants.SBLoginAcknowledged:       &LoginAcknowledged{},
	},
	CONFIGURATION: {
		constants.SBClientInformationConfiguration: &ClientInformation{},
		constants.SBPluginMessageConfiguration:     &PluginMessage{},
		constants.SBKnownPacksConfiguration:        &KnownPacks{},
		constants.SBAcknowledgeFinishConfiguration: &AcknowledgeFinishConfiguration{},
	},
	PLAY: {
		constants.SBConfirmTeleportation:      &TeleportConfirm{},
		constants.SBChatMessage:               &IncomingChatMessage{},
		constants.SBClientInformation:         &ClientSettings{},
		constants.SBPluginMessage:             &PluginMessage{},
		constants.SBKeepAlive:                 &KeepAlive{},
		constants.SBPlayerPosition:            &PlayerPosition{},
		constants.SBPlayerPositionAndRotation: &PlayerPositionAndRotation{},
		constants.SBPlayerRotation:            &PlayerRotation{},
		constants.SBPlayerMovement:            &PlayerMovement{},
		constants.SBPlayerAbilities:           &PlayerAbilities{},
		//constants.SBPlayerDigging:             &PlayerAction{},
		//constants.SBEntityAction:              &PlayerCommand{},
		constants.SBSetHeldItem:  &HeldItemChange{},
		constants.SBSwingArm:     &Animation{},
		constants.SBClientStatus: &ClientStatus{},
		constants.SBUseItemOn:    &PlayerBlockPlacement{},
	},
}

func Init(conn net.Conn, key *rsa.PrivateKey, minecraft *Minecraft) *Connection {
	fmt.Printf("--New Connection from %s!\n", conn.RemoteAddr().String())

	newConnection := Connection{
		State:             HANDSHAKING,
		Conn:              conn,
		Key:               key,
		Minecraft:         minecraft,
		EnableCompression: false,
		EnableEncryption:  false,
		Joined:            false,
		Ping:              -1,
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
			ID: time.Now().Unix(),
		})

		exception := c.SendPacket(&keepAlive)
		if exception != nil {
			fmt.Println("client has probably gone away")

			go c.Minecraft.PlayerLeave(c)
			break
		}
	}
}

func (c *Connection) Handle() {
	defer c.HandleError()
	packetStructScanner := structScanner.PacketStructScanner{}
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
			currentPacketType := iPacketType.(int)
			cursor += end

			if controllers[c.State] == nil {
				fmt.Println("!!! Bad State ", c.State)
				continue
			}

			if currentPacketType < 0 || len(controllers[c.State]) < currentPacketType || currentPacketType >= len(controllers[c.State]) || controllers[c.State][currentPacketType] == nil {
				fmt.Printf("!!! Serverbound packet of type 0x%X has no handler in state %d\n", currentPacketType, c.State)
				continue
			}

			fmt.Println("Current state ", c.State)
			fmt.Println("Incoming packet type ", currentPacketType)

			packetController := controllers[c.State][currentPacketType]

			packet := packetController.GetPacketStruct()

			//if cursor+length > len(decryptedBuf) {
			//	fmt.Println("Buffer overflow", cursor+length, len(decryptedBuf))
			//	continue
			//}

			packetBuffer := decryptedBuf[cursor:]

			fmt.Println(">>>INCOMING<<<")
			fmt.Println(hex.Dump(packetBuffer))

			packetStructScanner.StructScan(&packet, packetBuffer)

			packetController.Init(packet, c.Minecraft)
			packetController.Handle(decryptedBuf, c)
		}
	}
}

func (c *Connection) HandleError() {
	if r := recover(); r != nil {
		fmt.Println("Connection recovery:", r)
		debug.PrintStack()
	}
	_ = c.Conn.Close()
}

func (c *Connection) SendPacket(packet *packetType.Packet) error {
	var payload []byte
	packetID := byte((*packet).GetPacketId())
	packetStructScanner := structScanner.PacketStructScanner{}
	fmt.Printf("Sending packet 0x%X\n", packetID)

	if c.EnableCompression {
		uncompressedPayload := append([]byte{packetID}, packetStructScanner.UnmarshalData(*packet)...)
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
		payload = packetStructScanner.UnmarshalData(*packet)
		payload = append([]byte{packetID}, payload...)
		payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	}

	//if len(payload) < 1024 {
	fmt.Println(">>>OUTGOING<<<")
	fmt.Println(hex.Dump(payload))
	///}

	// fmt.Println("Writing payload")
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
