package controller

import (
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type LoginAcknowledged struct {
	CurrentPacket *packetType.LoginAcknowledged
}

func (lpr *LoginAcknowledged) GetPacketStruct() packetType.Packet {
	return &packetType.LoginAcknowledged{}
}

func (lpr *LoginAcknowledged) Init(currentPacket packetType.Packet) {
	lpr.CurrentPacket = currentPacket.(*packetType.LoginAcknowledged)
}

func (lpr *LoginAcknowledged) Handle(packet []byte, connection *Connection) {

	connection.State = CONFIGURATION
	knownPacks := packetType.Packet(&packetType.KnownPacks{
		ArrayLength: 1,
		Namespace:   "minecraft",
		ID:          "core",
		Version:     "1.21.5",
	})

	connection.SendPacket(&knownPacks)

	registryData := packetType.Packet(&packetType.RegistryData{
		RegistryID:  "minecraft:dimension_type",
		EntryLength: 1,
		Identifier:  "minecraft:overworld",
		NBTData: &packetType.DimensionTypeNBT{
			DimensionType: packetType.DimensionType{
				HasSkylight:                 1,
				HasCeiling:                  0,
				UltraWarm:                   0,
				Natural:                     1,
				CoordinateScale:             1,
				BedWorks:                    1,
				RespawnAnchorWorks:          0,
				MinY:                        0,
				Height:                      256,
				LogicalHeight:               250,
				Infiniburn:                  "#",
				Effects:                     "minecraft:overworld",
				AmbientLight:                0,
				PiglinSafe:                  0,
				HasRaids:                    1,
				MonsterSpawnLightLevel:      0,
				MonsterSpawnBlockLightLimit: 0,
			},
		},
	})

	// TODO:
	// 		minecraft:root/minecraft:cat_variant: Registry must be non-empty: minecraft:cat_variant
	//		minecraft:root/minecraft:chicken_variant: Registry must be non-empty: minecraft:chicken_variant
	//		minecraft:root/minecraft:cow_variant: Registry must be non-empty: minecraft:cow_variant
	//		minecraft:root/minecraft:frog_variant: Registry must be non-empty: minecraft:frog_variant
	//		minecraft:root/minecraft:painting_variant: Registry must be non-empty: minecraft:painting_variant
	//		minecraft:root/minecraft:pig_variant: Registry must be non-empty: minecraft:pig_variant
	//		minecraft:root/minecraft:wolf_sound_variant: Registry must be non-empty: minecraft:wolf_sound_variant
	//		minecraft:root/minecraft:wolf_variant: Registry must be non-empty: minecraft:wolf_variant
	// https://mcasset.cloud/1.21.5-pre2/data/minecraft

	connection.SendPacket(&registryData)

	finishConfiguration := packetType.Packet(&packetType.FinishConfiguration{})
	connection.SendPacket(&finishConfiguration)

}
