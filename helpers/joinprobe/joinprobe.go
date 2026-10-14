// Package joinprobe is the client half of joining a server, with nothing in it
// but the protocol: it walks a connection through handshake, login and
// configuration, then listens, and hands back every packet the server sent.
//
// It exists so the server's output can be checked from outside the server.
// Tests use it to look at what really goes on the wire, and cmd/mcverify feeds
// the result to the game's own packet decoders.
package joinprobe

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

// Phases a packet can arrive in, named as the game names the protocols.
const (
	PhaseLogin         = "login"
	PhaseConfiguration = "configuration"
	PhasePlay          = "play"
)

// Packet is one clientbound packet, as the server sent it.
type Packet struct {
	Phase string
	ID    int
	// Body is everything after the packet id.
	Body []byte
}

// Frame is the packet as the client's codec reads it: the id, then the body.
func (p Packet) Frame() []byte {
	return append(dataTypes.WriteVarInt(p.ID), p.Body...)
}

// Join connects as username and records what the server sends.
//
// It gives the server joinTimeout to get through login and configuration; once
// play has started it keeps listening until the server has been quiet for
// quiet, which is how it knows the server has finished its join sequence.
// Compression and encryption are not negotiated, so the server must have both
// off, as it does by default.
func Join(conn net.Conn, username string, joinTimeout time.Duration, quiet time.Duration) ([]Packet, error) {
	conn.SetDeadline(time.Now().Add(joinTimeout))

	err := writeFrame(conn, 0x00, // handshake
		dataTypes.WriteVarInt(constants.ProtocolVersion),
		dataTypes.WriteString("localhost"),
		dataTypes.WriteUnsignedShort(uint16(25565)),
		dataTypes.WriteVarInt(2)) // next state: login
	if err != nil {
		return nil, err
	}
	if err := writeFrame(conn, constants.SBHelloLogin, dataTypes.WriteString(username), make([]byte, 16)); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	phase := PhaseLogin
	var packets []Packet

	for {
		id, body, err := readFrame(reader)
		if err != nil {
			if phase == PhasePlay && (errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, io.EOF)) {
				return packets, nil // the server has said all it is going to
			}
			return packets, fmt.Errorf("reading a packet in the %s phase: %w", phase, err)
		}
		packets = append(packets, Packet{Phase: phase, ID: id, Body: body})

		switch {
		case phase == PhaseLogin && id == constants.CBLoginFinished:
			phase = PhaseConfiguration
			err = writeFrame(conn, constants.SBLoginAcknowledged)

		case phase == PhaseConfiguration && id == constants.CBSelectKnownPacksConfiguration:
			// claim the pack the real client has, which is what the server offered
			err = writeFrame(conn, constants.SBSelectKnownPacksConfiguration,
				dataTypes.WriteVarInt(1),
				dataTypes.WriteString("minecraft"), dataTypes.WriteString("core"), dataTypes.WriteString(constants.GameVersion))

		case phase == PhaseConfiguration && id == constants.CBFinishConfiguration:
			phase = PhasePlay
			err = writeFrame(conn, constants.SBFinishConfiguration)
		}
		if err != nil {
			return packets, err
		}

		if phase == PhasePlay {
			conn.SetReadDeadline(time.Now().Add(quiet))
		}
	}
}

// readFrame reads one length-prefixed packet and splits off its id.
func readFrame(r *bufio.Reader) (id int, body []byte, err error) {
	length, shift := 0, 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		length |= int(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	frame := make([]byte, length)
	if _, err := io.ReadFull(r, frame); err != nil {
		return 0, nil, err
	}

	idValue, cursor := dataTypes.ReadVarInt(frame)
	return idValue.(int), frame[cursor:], nil
}

func writeFrame(conn net.Conn, id int, body ...[]byte) error {
	payload := dataTypes.WriteVarInt(id)
	for _, part := range body {
		payload = append(payload, part...)
	}
	_, err := conn.Write(append(dataTypes.WriteVarInt(len(payload)), payload...))
	return err
}
