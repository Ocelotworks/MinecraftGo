package structScanner

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

func encodedKnownPacks(packs ...[3]string) []byte {
	b := dataTypes.WriteVarInt(len(packs))
	for _, p := range packs {
		for _, s := range p {
			b = append(b, dataTypes.WriteString(s)...)
		}
	}
	return b
}

func TestPrefixedArrayRoundTrip(t *testing.T) {
	tests := map[string][][3]string{
		"empty":    {},
		"single":   {{"minecraft", "core", "1.21.5"}},
		"multiple": {{"minecraft", "core", "1.21.5"}, {"example", "extra", "0.1"}, {"a", "b", "c"}},
	}

	for name, packs := range tests {
		t.Run(name, func(t *testing.T) {
			expected := make([]packetType.KnownPack, 0, len(packs))
			for _, p := range packs {
				expected = append(expected, packetType.KnownPack{Namespace: p[0], ID: p[1], Version: p[2]})
			}

			pss := PacketStructScanner{}
			wire := encodedKnownPacks(packs...)

			written := pss.UnmarshalData(&packetType.ServerboundKnownPacks{KnownPacks: expected})
			if !bytes.Equal(written, wire) {
				t.Fatalf("write mismatch\n got %x\nwant %x", written, wire)
			}

			var packet packetType.Packet = &packetType.ServerboundKnownPacks{}
			pss.StructScan(&packet, wire)
			got := packet.(*packetType.ServerboundKnownPacks).KnownPacks
			if len(got) != len(expected) || (len(expected) > 0 && !reflect.DeepEqual(got, expected)) {
				t.Fatalf("read mismatch\n got %v\nwant %v", got, expected)
			}
		})
	}
}

type arrayThenTrailer struct {
	Packs   []packetType.KnownPack `proto:"prefixedArray"`
	Trailer int                    `proto:"varInt"`
}

func (a *arrayThenTrailer) GetPacketId() int { return 0 }

// Fields after the array must be read from the correct offset.
func TestPrefixedArrayAdvancesCursor(t *testing.T) {
	pss := PacketStructScanner{}
	wire := append(encodedKnownPacks([3]string{"minecraft", "core", "1"}, [3]string{"x", "y", "z"}), dataTypes.WriteVarInt(300)...)

	var packet packetType.Packet = &arrayThenTrailer{}
	pss.StructScan(&packet, wire)

	got := packet.(*arrayThenTrailer)
	if len(got.Packs) != 2 || got.Trailer != 300 {
		t.Fatalf("unexpected result %+v", got)
	}
}

func TestPrefixedArrayInvalidLength(t *testing.T) {
	pss := PacketStructScanner{}

	// Claims a huge element count with no data behind it.
	var packet packetType.Packet = &packetType.ServerboundKnownPacks{}
	pss.StructScan(&packet, dataTypes.WriteVarInt(1_000_000))

	if got := packet.(*packetType.ServerboundKnownPacks).KnownPacks; len(got) != 0 {
		t.Fatalf("expected no elements, got %v", got)
	}
}
