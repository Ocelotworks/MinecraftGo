package controller

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Ocelotworks/MinecraftGo/constants"

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

	// writeLock serialises SendPacket.  Keepalives, chat and entity updates
	// are all sent from their own goroutines, and two of them writing at once
	// interleaves their bytes, which leaves the client reading a length
	// prefix from the middle of another packet.
	writeLock sync.Mutex
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
		constants.SBHelloLogin:             &LoginStart{},
		constants.SBKeyLogin:               &EncryptionResponse{},
		constants.SBCustomQueryAnswerLogin: &LoginPluginResponse{},
		constants.SBLoginAcknowledged:      &LoginAcknowledged{},
	},
	TRANSFER: {},
	CONFIGURATION: {
		constants.SBClientInformationConfiguration: &ClientInformation{},
		constants.SBCustomPayloadConfiguration:     &PluginMessage{},
		constants.SBSelectKnownPacksConfiguration:  &KnownPacks{},
		constants.SBFinishConfiguration:            &AcknowledgeFinishConfiguration{},
	},
	PLAY: {
		constants.SBClientTickEnd:        &ClientTickEnd{},
		constants.SBAcceptTeleportation:  &TeleportConfirm{},
		constants.SBChat:                 &IncomingChatMessage{},
		constants.SBClientInformation:    &ClientSettings{},
		constants.SBCustomPayload:        &PluginMessage{},
		constants.SBKeepAlive:            &KeepAlive{},
		constants.SBMovePlayerPos:        &PlayerPosition{},
		constants.SBMovePlayerPosRot:     &PlayerPositionAndRotation{},
		constants.SBMovePlayerRot:        &PlayerRotation{},
		constants.SBMovePlayerStatusOnly: &PlayerMovement{},
		constants.SBPlayerAbilities:      &PlayerAbilities{},
		//constants.SBPlayerAction:  &PlayerAction{},
		//constants.SBPlayerCommand: &PlayerCommand{},
		constants.SBSetCarriedItem: &HeldItemChange{},
		constants.SBInteract:       &Interact{},
		constants.SBClientCommand:  &ClientStatus{},
		constants.SBUseItemOn:      &PlayerBlockPlacement{},
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

// errUnframed means the stream is no longer sitting on a packet boundary: a
// length prefix has to be a VarInt, and nothing past five bytes is one.
var errUnframed = errors.New("packet length is not a VarInt")

// readPacketLength reads the VarInt every packet is prefixed with.  It reports
// the length, how many bytes the prefix took, and whether the buffer holds the
// whole prefix yet.
func readPacketLength(buf []byte) (length int, prefix int, complete bool, exception error) {
	for i := 0; i < len(buf); i++ {
		length |= int(buf[i]&0b01111111) << (7 * i)
		if buf[i]&0b10000000 == 0 {
			return length, i + 1, true, nil
		}
		if i == 4 {
			return 0, 0, false, errUnframed
		}
	}

	return 0, 0, false, nil
}

func (c *Connection) Handle() {
	defer c.HandleError()

	readBuffer := make([]byte, 4096)

	// One read is not one packet.  The client coalesces several packets into
	// a single write — in play it sends a handful every tick — and a big one
	// arrives split across several reads.  So bytes are collected here and
	// taken back out a whole packet at a time, as the length prefix says.
	var pending []byte

	for {
		count, err := c.Conn.Read(readBuffer)

		if err != nil {
			//fmt.Println("Error reading:", err.Error())
			c.Drop()
			return
		}

		received := make([]byte, count)

		if c.EnableEncryption {
			// CFB8 is a stream cipher, so this has to run over every byte
			// exactly once and in order.  Decrypting a read at a time rather
			// than a packet at a time is what keeps it in step.
			c.DecryptionCipher.XORKeyStream(received, readBuffer[:count])
		} else {
			copy(received, readBuffer[:count])
		}

		pending = append(pending, received...)

		for {
			length, prefix, complete, exception := readPacketLength(pending)

			if exception != nil {
				fmt.Println("!!! Lost the packet boundary:", exception)
				c.Drop()
				return
			}

			if !complete || len(pending) < prefix+length {
				// The rest of this packet has not arrived yet.
				break
			}

			c.HandlePacket(pending[prefix : prefix+length])
			pending = pending[prefix+length:]
		}

		if len(pending) == 0 {
			pending = nil
		}
	}
}

// HandlePacket decodes one whole packet — everything after its length prefix —
// and hands it to the controller registered for its type in the current state.
func (c *Connection) HandlePacket(packetData []byte) {
	packetStructScanner := structScanner.PacketStructScanner{}

	if c.EnableCompression {
		dataLength, cursor := dataTypes.ReadVarInt(packetData)
		packetData = packetData[cursor:]

		if dataLength.(int) > 0 {
			//fmt.Println("Data length is above compression threshold")
			r, exception := zlib.NewReader(bytes.NewReader(packetData))

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

			packetData = uncompressedPacket.Bytes()
		}
	}

	iPacketType, cursor := dataTypes.ReadVarInt(packetData)
	currentPacketType := iPacketType.(int)

	if controllers[c.State] == nil {
		fmt.Println("!!! Bad State ", c.State)
		return
	}

	if currentPacketType < 0 || currentPacketType >= len(controllers[c.State]) || controllers[c.State][currentPacketType] == nil {
		fmt.Printf("!!! Serverbound packet of type 0x%X has no handler in state %d\n", currentPacketType, c.State)
		return
	}

	packetController := controllers[c.State][currentPacketType]

	packet := packetController.GetPacketStruct()

	packetBuffer := packetData[cursor:]

	if c.State < PLAY {
		fmt.Println("Current state ", c.State)
		fmt.Println("Incoming packet type ", currentPacketType)
		fmt.Println(">>>INCOMING<<<")
		fmt.Println(hex.Dump(packetBuffer))
	}

	packetStructScanner.StructScan(&packet, packetBuffer)

	packetController.Init(packet, c.Minecraft)
	packetController.Handle(packetBuffer, c)
}

// Drop closes the connection and takes it out of the server's list.
func (c *Connection) Drop() {
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
}

func (c *Connection) HandleError() {
	if r := recover(); r != nil {
		fmt.Println("Connection recovery:", r)
		debug.PrintStack()
	}
	_ = c.Conn.Close()
}

// maxDumpedPacketBytes caps the outgoing hex dump.  A chunk packet carries a
// full set of sky light arrays, so dumping one in full is fifty thousand bytes
// of console, and a join sends a few hundred of them.
const maxDumpedPacketBytes = 512

func (c *Connection) SendPacket(packet *packetType.Packet) error {
	// Everything below writes to the same socket, and the callers are not all
	// on the same goroutine: keepalives, chat and entity updates each have
	// their own.  Two of them writing at once splices one packet into the
	// middle of another and the client stops being able to find a packet
	// boundary at all.
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	var payload []byte

	packetIDInt := (*packet).GetPacketId()
	packetIDBytes := dataTypes.WriteVarInt(packetIDInt)
	packetStructScanner := structScanner.PacketStructScanner{}
	fmt.Printf("Sending packet 0x%X\n", packetIDInt)

	if c.EnableCompression {
		uncompressedPayload := append(packetIDBytes, packetStructScanner.UnmarshalData(*packet)...)
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
		payload = append(packetIDBytes, packetStructScanner.UnmarshalData(*packet)...)
		payload = append(dataTypes.WriteVarInt(len(payload)), payload...)
	}

	fmt.Println(">>>OUTGOING<<<")
	if len(payload) <= maxDumpedPacketBytes {
		fmt.Println(hex.Dump(payload))
	} else {
		fmt.Printf("%s[... %d more bytes]\n\n", hex.Dump(payload[:maxDumpedPacketBytes]), len(payload)-maxDumpedPacketBytes)
	}

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
