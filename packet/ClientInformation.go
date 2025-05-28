package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ClientInformation struct {
	Locale              string `proto:"string"`
	ViewDistance        byte   `proto:"byte"`
	ChatMode            int    `proto:"varInt"`
	ChatColors          bool   `proto:"bool"`
	DisplayedSkinParts  byte   `proto:"unsignedByte"`
	MainHand            int    `proto:"varInt"`
	EnableTextFiltering bool   `proto:"bool"`
	AllowServerListings bool   `proto:"bool"`
	ParticleStatus      int    `proto:"varInt"`
}

func (ci *ClientInformation) GetPacketId() int {
	return constants.SBClientInformation
}
