package helpers

import (
	"math"
)

// BlocksPerEntry returns the number of blocks in a (64-bit) entry
func BlocksPerEntry(bitsPerBlock int64) int {
	return int(math.Ceil(float64(64 / bitsPerBlock)))
}

// BlockPosToSectionIndex returns the entry number and bit offset within that entry
func BlockPosToSectionIndex(bitsPerBlock int64, sectionY int, x int, y int, z int) (int, int) {
	blocksPerEntry := BlocksPerEntry(bitsPerBlock)
	// The Y needs to be converted to be relative to the section Y value
	relativeY := int(math.Abs(float64(sectionY*16 - y)))
	// The absolute index of the block, in the order of X by Z by (relative) Y
	blockIndex := (relativeY%16)*16*16 + z*16 + x

	// The section data index that contains this block
	entryNumber := int(math.Ceil(float64(blockIndex / blocksPerEntry)))
	// The bit offset within the entry that corresponds to this block. This starts from the MSB and works backwards
	bitOffset := 64 - ((blockIndex % blocksPerEntry) * int(bitsPerBlock)) - int(bitsPerBlock)

	return entryNumber, bitOffset
}
