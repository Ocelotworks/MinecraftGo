package controller

import (
	"encoding/binary"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/helpers/structScanner"
	"github.com/Ocelotworks/MinecraftGo/helpers/testutil"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/Ocelotworks/MinecraftGo/world"
)

// wireReader reads the primitives the client's stream codecs are built from.
type wireReader struct {
	t   *testing.T
	buf []byte
}

func (r *wireReader) take(n int, what string) []byte {
	r.t.Helper()
	if n < 0 || n > len(r.buf) {
		r.t.Fatalf("reading %s: need %d bytes, %d left", what, n, len(r.buf))
	}
	out := r.buf[:n]
	r.buf = r.buf[n:]
	return out
}

func (r *wireReader) varInt(what string) int {
	r.t.Helper()
	value, shift := 0, 0
	for {
		b := r.take(1, what)[0]
		value |= int(b&0x7F) << shift
		if b&0x80 == 0 {
			return value
		}
		shift += 7
	}
}

// byteArray is ByteBufCodecs.byteArray(max).
func (r *wireReader) byteArray(max int, what string) []byte {
	r.t.Helper()
	n := r.varInt(what + " length")
	if n > max {
		r.t.Fatalf("%s: byte array of %d is bigger than the allowed %d", what, n, max)
	}
	return r.take(n, what)
}

// bitSet is ByteBufCodecs.BIT_SET: java.util.BitSet.toByteArray() as a byte array
// with a VarInt length.  The test works on the raw bytes, since only the number
// of bits set matters to it.
func (r *wireReader) bitSet(what string) []byte {
	r.t.Helper()
	return r.byteArray(1<<20, what)
}

func popCount(bytes []byte) int {
	n := 0
	for _, b := range bytes {
		for ; b != 0; b &= b - 1 {
			n++
		}
	}
	return n
}

// The chunk packet has to be readable by ClientboundLevelChunkWithLightPacket's
// STREAM_CODEC: x, z, then the chunk data (heightmaps, section buffer, block
// entities), then the light data (four BitSets and two lists of 2048-byte arrays).
func TestChunkPacketMatchesClientCodec(t *testing.T) {
	sectionData := []byte{1, 2, 3, 4, 5}
	light := newChunkLight(world.Overworld)

	packet := packetType.Packet(newChunkPacket(-3, 7, sectionData, light))
	scanner := structScanner.PacketStructScanner{}
	r := &wireReader{t: t, buf: scanner.UnmarshalData(packet)}

	if x := int32(binary.BigEndian.Uint32(r.take(4, "x"))); x != -3 {
		t.Errorf("x = %d", x)
	}
	if z := int32(binary.BigEndian.Uint32(r.take(4, "z"))); z != 7 {
		t.Errorf("z = %d", z)
	}

	// heightmaps: map of (type, long[])
	for i, n := 0, r.varInt("heightmap count"); i < n; i++ {
		r.varInt("heightmap type")
		r.take(8*r.varInt("heightmap length"), "heightmap longs")
	}

	// chunk buffer
	if buffer := r.byteArray(2097152, "chunk buffer"); string(buffer) != string(sectionData) {
		t.Errorf("chunk buffer = %v, want %v", buffer, sectionData)
	}

	// block entities
	if n := r.varInt("block entity count"); n != 0 {
		t.Fatalf("test only handles an empty block entity list, got %d", n)
	}

	// light data
	skyMask := r.bitSet("sky mask")
	blockMask := r.bitSet("block mask")
	emptySkyMask := r.bitSet("empty sky mask")
	emptyBlockMask := r.bitSet("empty block mask")

	skyCount := r.varInt("sky array count")
	if want := popCount(skyMask); skyCount != want {
		t.Errorf("%d sky arrays for a mask with %d bits set", skyCount, want)
	}
	for i := 0; i < skyCount; i++ {
		r.byteArray(2048, "sky array")
	}

	blockCount := r.varInt("block array count")
	if want := popCount(blockMask); blockCount != want {
		t.Errorf("%d block arrays for a mask with %d bits set", blockCount, want)
	}
	for i := 0; i < blockCount; i++ {
		r.byteArray(2048, "block array")
	}

	if popCount(emptySkyMask) != 0 || popCount(emptyBlockMask) != lightSectionCount(world.Overworld) {
		t.Errorf("empty masks have %d and %d bits set", popCount(emptySkyMask), popCount(emptyBlockMask))
	}
	if len(r.buf) != 0 {
		t.Errorf("%d bytes left over after the packet", len(r.buf))
	}

}

// decodeChunkPacket reads one chunk packet the way the client does and returns the
// section buffer, failing the test if any length or count is inconsistent.
func decodeChunkPacket(t *testing.T, wire []byte) []byte {
	t.Helper()
	r := &wireReader{t: t, buf: wire}

	r.take(8, "x and z")
	for i, n := 0, r.varInt("heightmap count"); i < n; i++ {
		r.varInt("heightmap type")
		r.take(8*r.varInt("heightmap length"), "heightmap longs")
	}
	buffer := r.byteArray(2097152, "chunk buffer")
	if n := r.varInt("block entity count"); n != 0 {
		t.Fatalf("unexpected %d block entities", n)
	}

	skyMask := r.bitSet("sky mask")
	blockMask := r.bitSet("block mask")
	r.bitSet("empty sky mask")
	r.bitSet("empty block mask")

	for _, mask := range [][]byte{skyMask, blockMask} {
		count := r.varInt("light array count")
		if count != popCount(mask) {
			t.Fatalf("%d light arrays for a mask with %d bits set", count, popCount(mask))
		}
		for i := 0; i < count; i++ {
			r.byteArray(2048, "light array")
		}
	}
	if len(r.buf) != 0 {
		t.Fatalf("%d bytes left over after the packet", len(r.buf))
	}
	return buffer
}

// Every chunk the server would send around the spawn must decode, not just the
// synthetic one above.
func TestRealChunkPacketsMatchClientCodec(t *testing.T) {
	testutil.ChdirToRepoRoot(t) // data/ is relative to the repo root

	mc := CreateMinecraft()
	light := newChunkLight(mc.ChunkManager.Height)
	scanner := structScanner.PacketStructScanner{}

	chunks := mc.ChunkManager.ChunksAround(0, 0, maxSendViewDistance)
	if len(chunks) == 0 {
		t.Fatal("no chunks around spawn")
	}
	for _, chunk := range chunks {
		sections := chunk.WriteNetworkData(mc.WorldRegistry)
		packet := packetType.Packet(newChunkPacket(int(chunk.X), int(chunk.Z), sections, light))

		buffer := decodeChunkPacket(t, scanner.UnmarshalData(packet))
		if len(buffer) != len(sections) {
			t.Errorf("chunk %d,%d: buffer is %d bytes, wrote %d", chunk.X, chunk.Z, len(buffer), len(sections))
		}
	}
	t.Logf("%d chunks decoded", len(chunks))
}
