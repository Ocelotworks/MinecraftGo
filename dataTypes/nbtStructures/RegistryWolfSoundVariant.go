package nbtStructures

type RegistryWolfSoundVariantOuter struct {
	Inner RegistryWolfSoundVariant `nbt:""`
}

type RegistryWolfSoundVariant struct {
	AmbientSound string `nbt:"ambient_sound"`
	DeathSound   string `nbt:"death_sound"`
	GrowlSound   string `nbt:"growl_sound"`
	HurtSound    string `nbt:"hurt_sound"`
	PantSound    string `nbt:"pant_sound"`
	WhineSound   string `nbt:"whine_sound"`
}
