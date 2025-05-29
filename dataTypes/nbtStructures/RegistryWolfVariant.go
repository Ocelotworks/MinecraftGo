package nbtStructures

type RegistryWolfVariantOuter struct {
	Inner RegistryWolfVariant `nbt:""`
}

type RegistryWolfVariant struct {
	Assets WolfAssets `nbt:"assets"`

	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"`
}

type WolfAssets struct {
	WildTexture  string `nbt:"wild"`
	TameTexture  string `nbt:"tame"`
	AngryTexture string `nbt:"angry"`
}
