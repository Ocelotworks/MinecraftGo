package nbtStructures

type RegistryPaintingVariantOuter struct {
	Inner RegistryPaintingVariant `nbt:""`
}

type RegistryPaintingVariant struct {
	AssetId string `nbt:"asset_id"`
	Height  int32  `nbt:"height"`
	Width   int32  `nbt:"width"`
	Title   string `nbt:"title"`
	Author  string `nbt:"author"`
}
