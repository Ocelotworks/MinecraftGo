package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLoadHelloWorld(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/nbt-test/hello_world.nbt")
	assert.Nil(t, exception)

	compound := ReadNBT(inData)
	out, _ := json.Marshal(compound)
	fmt.Println(string(out))

	assert.IsType(t, &Compound{}, compound.Data["hello world"])

	innerCompound := compound.Data["hello world"].(*Compound)
	assert.IsType(t, &String{}, innerCompound.Data["name"])
	assert.Equal(t, "Bananrama", innerCompound.Data["name"].GetValue())

	// our compound ends with an extra 0x00, because it expects to close the outer compound, but that is never done in a file
	assert.Equal(t, inData, compound.Write()[:len(inData)])
}

func TestLoadBigNbt(t *testing.T) {
	buf, exception := ioutil.ReadFile("../../data/nbt-test/bigtest.nbt")
	assert.Nil(t, exception)

	compressed := bytes.NewReader(buf)
	zr, exception := gzip.NewReader(compressed)
	assert.Nil(t, exception)

	uncompressed, exception := ioutil.ReadAll(zr)
	assert.Nil(t, exception)

	compound := ReadNBT(uncompressed)
	out, _ := json.Marshal(compound)
	fmt.Println(string(out))

	assert.IsType(t, &Compound{}, compound.Data["Level"])

	level := compound.Data["Level"].(*Compound)
	assert.Equal(t, int32(2147483647), level.Data["intTest"].GetValue())
	assert.Equal(t, byte(127), level.Data["byteTest"].GetValue())
	assert.Equal(t, "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!", level.Data["stringTest"].GetValue())
	assert.Equal(t, 0.49312871321823148, level.Data["doubleTest"].GetValue())
	assert.Equal(t, float32(0.49823147058486938), level.Data["floatTest"].GetValue())
	assert.Equal(t, int64(9223372036854775807), level.Data["longTest"].GetValue())
	assert.Equal(t, int16(32767), level.Data["shortTest"].GetValue())

	// Can't compare the contents as map order is not
	assert.Len(t, compound.Write()[:len(uncompressed)], len(uncompressed))
}

func TestLoadMapNBT(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/worlds/MCGO_FlatTest/level.dat")
	assert.Nil(t, exception)

	compressed := bytes.NewReader(inData)
	zr, exception := gzip.NewReader(compressed)
	assert.Nil(t, exception)

	uncompressed, exception := ioutil.ReadAll(zr)
	assert.Nil(t, exception)

	compound := ReadNBT(uncompressed)
	out, _ := json.Marshal(compound)
	fmt.Println(string(out))
}

func TestLoadMap2NBT(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/worlds/world/level.dat")
	assert.Nil(t, exception)

	compressed := bytes.NewReader(inData)
	zr, exception := gzip.NewReader(compressed)
	assert.Nil(t, exception)

	uncompressed, exception := ioutil.ReadAll(zr)
	assert.Nil(t, exception)

	compound := ReadNBT(uncompressed)
	out, _ := json.Marshal(compound)
	fmt.Println(string(out))
}
