package entity

import "../packet"

type Minecraft struct {
	Connections      []*packet.Connection
	ServerName       ChatMessageComponent
	MaxPlayers       int
	EnableEncryption bool
}
