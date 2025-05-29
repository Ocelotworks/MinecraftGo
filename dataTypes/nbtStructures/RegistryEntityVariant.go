package nbtStructures

type RegistryEntityVariantOuter struct {
	Inner RegistryEntityVariant `nbt:""`
}

type RegistryEntityVariant struct {
	AssetId         string           `nbt:"asset_id"`
	Model           string           `nbt:"model"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"`
}

type SpawnCondition struct {
	Priority int32 `nbt:"priority"`
	//Condition *SpawnConditionData `nbt:"condition"` // not required
}

type SpawnConditionData struct {
	Type       string `nbt:"type"`
	Structures string `nbt:"structures"`
}
