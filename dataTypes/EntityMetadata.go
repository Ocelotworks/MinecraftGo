package dataTypes

import (
	"../entity"
	"encoding/json"
)

type Metadata interface {
	Write() []byte
}

type EntityMetadata struct {
	Effect            *byte
	Air               *int
	CustomName        *entity.ChatMessageComponent
	CustomNameVisible *bool
	Silent            *bool
	NoGravity         *bool
	Pose              *int
}

//This is going to be horrific
func (em *EntityMetadata) Write() []byte {
	output := make([]byte, 0)
	if em.Effect != nil {
		output = append(output, WriteEntityMetadataProperty(0, EMValueByte, WriteUnsignedByte(*em.Effect))...)
	}
	if em.Air != nil {
		output = append(output, WriteEntityMetadataProperty(1, EMValueVarInt, WriteVarInt(*em.Air))...)
	}
	if em.CustomName != nil {
		output, exception := json.Marshal(em.CustomName)
		if exception != nil {
			output = append(output, WriteEntityMetadataProperty(2, EMValueOptChat, WriteBoolean(false))...)
		} else {
			output = append(output, WriteEntityMetadataProperty(2, EMValueOptChat, WriteBoolean(true))...)
			output = append(output, WriteString(output)...)
		}
	}
	if em.CustomNameVisible != nil {
		output = append(output, WriteEntityMetadataProperty(3, EMValueBool, WriteBoolean(*em.CustomNameVisible))...)
	}

	if em.Silent != nil {
		output = append(output, WriteEntityMetadataProperty(4, EMValueBool, WriteBoolean(*em.Silent))...)
	}

	if em.NoGravity != nil {
		output = append(output, WriteEntityMetadataProperty(5, EMValueBool, WriteBoolean(*em.NoGravity))...)
	}

	if em.Pose != nil {
		output = append(output, WriteEntityMetadataProperty(6, EMValuePose, WriteVarInt(*em.Pose))...)
	}

	return output
}

const (
	PoseStanding   = 0
	PoseFallFlying = 1
	PoseSleeping   = 2
	PoseSwimming   = 3
	PoseSpinAttack = 4
	PoseSneaking   = 5
	PoseDying      = 6
)

const (
	EMValueByte         = 0
	EMValueVarInt       = 1
	EMValueFloat        = 2
	EMValueString       = 3
	EMValueChat         = 4
	EMValueOptChat      = 5
	EMValueSlot         = 6
	EMValueBool         = 7
	EMValueRotation     = 8
	EmValuePosition     = 9
	EMValueOptPosition  = 10
	EMValueDirection    = 11
	EMValueOptUUID      = 12
	EMValueOptBlockID   = 13
	EMValueNBT          = 14
	EMValueParticle     = 15
	EMValueVillagerData = 16
	EMValueOptVarInt    = 17
	EMValuePose         = 18
)

func WriteEntityMetadata(input interface{}) []byte {
	metadata := input.(Metadata)
	output := metadata.Write()
	return append(output, 0xff)
}

func WriteEntityMetadataProperty(index byte, valueType int, value []byte) []byte {
	output := WriteUnsignedByte(index)
	output = append(output, WriteVarInt(valueType)...)
	return append(output, value...)
}
