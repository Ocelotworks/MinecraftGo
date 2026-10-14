package dataTypes

import "encoding/binary"

// A BitSet is ByteBufCodecs.BIT_SET: a VarInt byte count followed by
// java.util.BitSet.toByteArray(), which is the bits packed least significant
// first with every trailing zero byte dropped.  An empty set is a single 0.
//
// This used to be the FriendlyByteBuf.writeBitSet form, a VarInt count of
// big-endian longs, and a client reading that form takes the 0x00 high bytes of
// the first long as a series of empty sets.
//
// The sets are passed around as the longs java.util.BitSet.valueOf takes, bit n
// being bit n%64 of long n/64.

func ReadBitSet(buf []byte) (interface{}, int) {
	size, cursor := ReadVarInt(buf)
	byteCount := size.(int)

	output := make([]int64, (byteCount+7)/8)
	for i := 0; i < byteCount; i++ {
		output[i/8] |= int64(buf[cursor+i]) << uint(8*(i%8))
	}

	return output, cursor + byteCount
}

func WriteBitSet(input interface{}) []byte {
	longs := input.([]int64)

	bytes := make([]byte, 0, len(longs)*8)
	for _, long := range longs {
		bytes = binary.LittleEndian.AppendUint64(bytes, uint64(long))
	}
	for len(bytes) > 0 && bytes[len(bytes)-1] == 0 {
		bytes = bytes[:len(bytes)-1]
	}

	return append(WriteVarInt(len(bytes)), bytes...)
}
