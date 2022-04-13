package dataTypes

// Hacky fixes
type CodecOuterCompound struct {
	Inner DimensionCodec `nbt:""`
}

type DimensionOuterCompound struct {
	Inner DimensionType `nbt:""`
}

type DimensionCodec struct {
	DimensionType DimensionTypeRegistry `nbt:"minecraft:dimension_type"`
	BiomeRegistry BiomeRegistry         `nbt:"minecraft:worldgen/biome"`
}

type DimensionTypeRegistry struct {
	Type  string                       `nbt:"type"`
	Value []DimensionTypeRegistryEntry `nbt:"value"`
}

type DimensionTypeRegistryEntry struct {
	Name    string        `nbt:"name"`
	Id      int32         `nbt:"id"`
	Element DimensionType `nbt:"element"`
}

type DimensionType struct {
	PiglinSafe   byte    `nbt:"piglin_safe"`
	Natural      byte    `nbt:"natural"`
	AmbientLight float32 `nbt:"ambient_light"`
	// FixedTime // TODO: optional fields
	Infiniburn         string  `nbt:"infiniburn"`
	RespawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
	HasSkylight        byte    `nbt:"has_skylight"`
	BedWorks           byte    `nbt:"bed_works"`
	Effects            string  `nbt:"effects"`
	HasRaids           byte    `nbt:"has_raids"`
	MinY               int32   `nbt:"min_y"`
	Height             int32   `nbt:"height"`
	LogicalHeight      int32   `nbt:"logical_height"`
	CoordinateScale    float64 `nbt:"coordinate_scale"`
	Ultrawarm          byte    `nbt:"ultrawarm"`
	HasCeiling         byte    `nbt:"has_ceiling"`
}

type BiomeRegistry struct {
	Type  string               `nbt:"type"`
	Value []BiomeRegistryEntry `nbt:"value"`
}

type BiomeRegistryEntry struct {
	Name    string `nbt:"name"`
	Id      int32  `nbt:"id"`
	Element Biome  `nbt:"element"`
}

type Biome struct {
	Precipitation string  `nbt:"precipitation"`
	Depth         float32 `nbt:"depth"`
	Temperature   float32 `nbt:"temperature"`
	Scale         float32 `nbt:"scale"`
	Downfall      float32 `nbt:"downfall"`
	Category      string  `nbt:"category"`
	//TemperatureModifier
	Effects BiomeEffect `nbt:"effects"`
	// Particle
}

type BiomeEffect struct {
	SkyColor      int32 `nbt:"sky_color"`
	WaterFogColor int32 `nbt:"water_fog_color"`
	FogColor      int32 `nbt:"fog_color"`
	WaterColor    int32 `nbt:"water_color"`
	//FoliageColor int32 `nbt:"foliage_color"`
	//GrassColor
	//GrassColorModifier
	//Music
	//AmbientSound
	//AdditionsSound
	//MoodSound
}
