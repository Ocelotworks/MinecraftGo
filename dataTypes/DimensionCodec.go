package dataTypes

type DimensionCodec struct {
	DimensionType DimensionType `nbt:"minecraft:dimension_type"`
	BiomeRegistry BiomeRegistry `nbt:"minecraft:worldgen/biome"`
}

type DimensionType struct {
}

type BiomeRegistry struct {
}
