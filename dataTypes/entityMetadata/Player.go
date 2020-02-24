package entityMetadata

type PlayerMetadata struct {
	LivingEntityMetadata
	AdditionalHearts    float32
	Score               int
	DisplayedSkinParts  byte
	MainHand            byte
	LeftShoulderEntity  []interface{} //NBT
	RightShoulderEntity []interface{} //NBT
}

func (em *PlayerMetadata) Write() []byte {
	output := em.LivingEntityMetadata.Write()
	//TODO
	return output
}
