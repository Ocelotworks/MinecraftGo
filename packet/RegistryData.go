package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type RegistryData struct {
	RegistryID  string            `proto:"string"`
	EntryLength int               `proto:"varInt"`
	Identifier  string            `proto:"string"`
	NBTData     *DimensionTypeNBT `proto:"prefixedOptionalNbt"`
	//Entries     RegistryEntry `proto:"nbt"`
}

type RegistryEntry struct {
	Identifier string        `proto:"string"`
	NBTData    DimensionType `proto:"prefixedOptionalNbt"`
}

func (apd *RegistryData) GetPacketId() int {
	return constants.CBRegistryDataConfiguration
}

type DimensionTypeNBT struct {
	DimensionType DimensionType `nbt:""`
}

type DimensionType struct {
	HasSkylight                 byte    `nbt:"has_skylight"`
	HasCeiling                  byte    `nbt:"has_ceiling"`
	UltraWarm                   byte    `nbt:"ultrawarm"`
	Natural                     byte    `nbt:"natural"`
	CoordinateScale             float64 `nbt:"coordinate_scale"`
	BedWorks                    byte    `nbt:"bed_works"`
	RespawnAnchorWorks          byte    `nbt:"respawn_anchor_works"`
	MinY                        int32   `nbt:"min_y"`
	Height                      int32   `nbt:"height"`
	LogicalHeight               int32   `nbt:"logical_height"`
	Infiniburn                  string  `nbt:"infiniburn"`
	Effects                     string  `nbt:"effects"`
	AmbientLight                float32 `nbt:"ambient_light"`
	PiglinSafe                  byte    `nbt:"piglin_safe"`
	HasRaids                    byte    `nbt:"has_raids"`
	MonsterSpawnLightLevel      int32   `nbt:"monster_spawn_light_level"`
	MonsterSpawnBlockLightLimit int32   `nbt:"monster_spawn_block_light_limit"`
}
