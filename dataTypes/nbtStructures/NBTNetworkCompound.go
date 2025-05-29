package nbtStructures

type NBTNetworkCompound[T any] struct {
	Inner T `nbt:""`
}
