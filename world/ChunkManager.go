package world

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

// chunksPerRegionAxis is the width of a region file in chunks; a chunk's region
// is found by shifting its coordinate right by 5, which floors for negatives
// where dividing by 32 would not.
const chunksPerRegionAxis = 32

type regionPos struct {
	X, Z int
}

type ChunkManager struct {
	// WorldPath is the directory holding the world, i.e. the parent of region/.
	WorldPath string
	// Height is the vertical extent of the dimension these chunks belong to.
	Height WorldHeight

	// Chunks is indexed by chunk X then chunk Z.
	Chunks map[int32]map[int32]*Chunk

	// loadedRegions records which region files have been read so a second
	// request for the same region is not a second parse.
	loadedRegions map[regionPos]bool
}

func NewChunkManager(worldPath string, height WorldHeight) *ChunkManager {
	return &ChunkManager{
		WorldPath:     worldPath,
		Height:        height,
		Chunks:        make(map[int32]map[int32]*Chunk),
		loadedRegions: make(map[regionPos]bool),
	}
}

func (cm *ChunkManager) loadRegionFile(x int, z int) (*dataTypes.RegionMetadata, error) {
	path := filepath.Join(cm.WorldPath, "region", fmt.Sprintf("r.%d.%d.mca", x, z))

	inData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading region file %s: %w", path, err)
	}

	region := dataTypes.ReadRegionFile(inData)
	return &region, nil
}

// LoadRegion reads one region file and adds its chunks to whatever is already
// loaded.  Regions that have already been read, and regions with no file on
// disc, are skipped rather than treated as an error — an absent region file
// just means that part of the world has never been generated.
func (cm *ChunkManager) LoadRegion(x int, z int) error {
	if cm.loadedRegions[regionPos{x, z}] {
		return nil
	}

	regionMetadata, err := cm.loadRegionFile(x, z)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cm.loadedRegions[regionPos{x, z}] = true
			return nil
		}
		return err
	}

	// Mark the region loaded before parsing so a chunk that fails to parse
	// does not cause the whole file to be re-read on the next request.
	cm.loadedRegions[regionPos{x, z}] = true

	for _, regionChunk := range regionMetadata.Chunks {
		if regionChunk == nil {
			continue
		}

		chunk, err := ParseChunkFromRegionChunk(regionChunk, cm.Height)
		if err != nil {
			// One malformed chunk should not take the rest of the region
			// with it; the client renders a hole instead of disconnecting.
			fmt.Println("Skipping chunk:", err)
			continue
		}

		cm.store(chunk)
	}

	return nil
}

func (cm *ChunkManager) store(chunk *Chunk) {
	column, ok := cm.Chunks[chunk.X]
	if !ok {
		column = make(map[int32]*Chunk)
		cm.Chunks[chunk.X] = column
	}
	column[chunk.Z] = chunk
}

// GetChunk returns the chunk at a chunk coordinate, reading its region file
// first if that has not happened yet.  It returns nil when the chunk does not
// exist in the world.
func (cm *ChunkManager) GetChunk(x int32, z int32) *Chunk {
	if chunk, ok := cm.Chunks[x][z]; ok {
		return chunk
	}

	regionX, regionZ := int(x)>>5, int(z)>>5
	if cm.loadedRegions[regionPos{regionX, regionZ}] {
		return nil
	}

	if err := cm.LoadRegion(regionX, regionZ); err != nil {
		fmt.Println("Loading region:", err)
		return nil
	}

	return cm.Chunks[x][z]
}

// ChunksAround returns the chunks within viewDistance of a chunk coordinate,
// nearest first so the client fills in the area around the player before the
// edges of the view.
func (cm *ChunkManager) ChunksAround(centreX int32, centreZ int32, viewDistance int) []*Chunk {
	chunks := make([]*Chunk, 0)

	for ring := 0; ring <= viewDistance; ring++ {
		for x := centreX - int32(ring); x <= centreX+int32(ring); x++ {
			for z := centreZ - int32(ring); z <= centreZ+int32(ring); z++ {
				// Only the outermost ring is new on each pass.
				if abs32(x-centreX) != int32(ring) && abs32(z-centreZ) != int32(ring) {
					continue
				}
				if chunk := cm.GetChunk(x, z); chunk != nil {
					chunks = append(chunks, chunk)
				}
			}
		}
	}

	return chunks
}

func abs32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}
