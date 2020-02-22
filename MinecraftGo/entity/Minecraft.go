package entity

type Minecraft struct {
	//Connections      []*packet.Connection
	ServerName       ChatMessageComponent
	MaxPlayers       int
	EnableEncryption bool
}
