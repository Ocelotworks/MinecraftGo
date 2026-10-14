package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

// UpdateTagsConfiguration sends tag bindings for synchronized registries
// during the configuration phase. Data is a manually-assembled binary
// payload: [VarInt registryCount] { [String registry] [VarInt tagCount]
// { [String tag] [VarInt entryCount] { [VarInt entryID] ... } } }
// Entry IDs are 0-based indices from the RegistryData packet for that registry.
type UpdateTagsConfiguration struct {
	Data []byte `proto:"raw"`
}

func (t *UpdateTagsConfiguration) GetPacketId() int {
	return constants.CBUpdateTagsConfiguration
}
