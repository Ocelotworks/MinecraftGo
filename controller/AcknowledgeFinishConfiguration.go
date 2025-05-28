package controller

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type AcknowledgeFinishConfiguration struct {
	CurrentPacket *packetType.AcknowledgeFinishConfiguration
}

func (lpr *AcknowledgeFinishConfiguration) GetPacketStruct() packetType.Packet {
	return &packetType.AcknowledgeFinishConfiguration{}
}

func (lpr *AcknowledgeFinishConfiguration) Init(currentPacket packetType.Packet) {
	lpr.CurrentPacket = currentPacket.(*packetType.AcknowledgeFinishConfiguration)
}

func (lpr *AcknowledgeFinishConfiguration) Handle(packet []byte, connection *Connection) {
	connection.State = PLAY

	//os.WriteFile("ours2.nbt", compound.Write(), 0644)
	//os.WriteFile("theirs2.nbt",, 0644)

	joinGame := packetType.Packet(&packetType.JoinGame{
		EntityID:            connection.Player.EntityID,
		IsHardcore:          false,
		DimensionNames:      []string{"minecraft:overworld"}, // TODO: this but properly
		MaxPlayers:          connection.Minecraft.MaxPlayers,
		ViewDistance:        32,
		SimulationDistance:  32,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		DoLimitedCrafting:   false,
		DimensionType:       0,
		DimensionName:       "minecraft:overworld",
		HashedSeed:          71495747907944700,
		Gamemode:            1,
		PreviousGamemode:    1,
		IsDebug:             false,
		IsFlat:              false,
		HasDeathLocation:    false,
		PortalCooldown:      0,
		SeaLevel:            64,
		EnforcesSecureChat:  false,
	})

	connection.SendPacket(&joinGame)

	pluginMessage := packetType.Packet(&packetType.PluginMessage{
		IsServer:   false,
		Identifier: "minecraft:brand",
		ByteArray:  dataTypes.WriteString("BigPMC"),
	})

	connection.SendPacket(&pluginMessage)

	difficulty := packetType.Packet(&packetType.ServerDifficulty{
		Difficulty:       0,
		DifficultyLocked: false,
	})

	connection.SendPacket(&difficulty)
}
