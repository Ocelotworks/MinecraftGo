package dataTypes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The test world the server loads at start up (controller.WorldPath).
const testRegionFile = "../data/worlds/MCGO_FlatTest/region/r.0.0.mca"

func TestRegionLoad(t *testing.T) {
	inData, err := os.ReadFile(testRegionFile)
	require.NoError(t, err)

	region := ReadRegionFile(inData)

	assert.Len(t, region.ChunkOffsets, 1024)
	assert.Len(t, region.Chunks, 1024)

	loaded := 0
	for slot, chunk := range region.Chunks {
		if chunk == nil {
			// a chunk with no offset entry has never been generated
			assert.Zero(t, region.ChunkOffsets[slot].Size, "slot %d has data but did not load", slot)
			continue
		}
		loaded++

		assert.NotEmpty(t, chunk.Sections, "chunk %d,%d has no sections", chunk.XPos, chunk.ZPos)

		// a chunk's slot in the region is its position modulo 32, x fastest
		assert.Equal(t, int32(slot%32), chunk.XPos&31, "slot %d holds chunk x %d", slot, chunk.XPos)
		assert.Equal(t, int32(slot/32), chunk.ZPos&31, "slot %d holds chunk z %d", slot, chunk.ZPos)

		assert.Same(t, chunk, region.GetChunk(chunk.XPos, chunk.ZPos))
	}
	assert.NotZero(t, loaded, "no chunks loaded from %s", testRegionFile)
}
