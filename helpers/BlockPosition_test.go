package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlocksPerEntry(t *testing.T) {
	// Values between 4 and 8 are expected
	assert.Equal(t, 16, BlocksPerEntry(4))
	assert.Equal(t, 12, BlocksPerEntry(5))
	assert.Equal(t, 10, BlocksPerEntry(6))
	assert.Equal(t, 9, BlocksPerEntry(7))
	assert.Equal(t, 8, BlocksPerEntry(8))
}

func TestBlockPosToSectionIndex(t *testing.T) {

	t.Run("X Coord Calculation", func(t *testing.T) {
		// (0,0,0) should be the highest 4 bits of the first entry, which will start at 60 and move left 4 places
		doubleEqual(t, 4, 0, 0, 0, 0, 0, 60)

		// The first X coordinate is the second value in total, which is another 4 bits down from the original
		doubleEqual(t, 4, 0, 1, 0, 0, 0, 56)

		// The 16th X coordinate is the end of the very first entry, so the bit offset starts at 0
		doubleEqual(t, 4, 0, 15, 0, 0, 0, 0)

		// The 17th X coordinate is the start of the second entry, so starts at the highest 4 bits
		doubleEqual(t, 4, 0, 16, 0, 0, 1, 60)

	})

	t.Run("Y Coord Calculation", func(t *testing.T) {
		// SectionY starting at -16 means that the Y value is 0 relative to the start of the section, so this is entry 0 bit 60 again
		doubleEqual(t, 4, -16, 0, -16, 0, 0, 60)

		// With 4 bits per block, each entry is one line of blocks, so 1 Y value upwards is the 16th entry bit 60
		doubleEqual(t, 4, 0, 0, 1, 0, 16, 60)
		doubleEqual(t, 4, 0, 1, 1, 0, 16, 56)
		doubleEqual(t, 4, 0, 15, 1, 0, 16, 0)
		doubleEqual(t, 4, 0, 0, 2, 0, 32, 60)
		doubleEqual(t, 4, 0, 0, 3, 0, 48, 60)

		doubleEqual(t, 4, -16, 0, -15, 0, 16, 60)
		doubleEqual(t, 4, -16, 0, -14, 0, 32, 60)
		doubleEqual(t, 4, -16, 0, 0, 0, 256, 60)
	})

}

func doubleEqual(t *testing.T, bitsPerBlock int64, sectionY int, x int, y int, z int, expectedEntryNumber int, expectedBitOffset int) {
	t.Run(fmt.Sprintf("bpb=%d;section=%d;(%d,%d,%d)", bitsPerBlock, sectionY, x, y, z), func(t *testing.T) {
		entryNumber, bitOffset := BlockPosToSectionIndex(bitsPerBlock, sectionY, x, y, z)
		assert.Equal(t, expectedEntryNumber, entryNumber)
		assert.Equal(t, expectedBitOffset, bitOffset)
	})
}
