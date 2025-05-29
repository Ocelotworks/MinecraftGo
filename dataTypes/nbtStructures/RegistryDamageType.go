package nbtStructures

type RegistryDamageTypeOuter struct {
	Inner RegistryDamageType `nbt:""`
}

type RegistryDamageType struct {
	Effects    string  `nbt:"effects"`
	Exhaustion float32 `nbt:"exhaustion"`
	MessageId  string  `nbt:"message_id"`
	Scaling    string  `nbt:"scaling"`
}
