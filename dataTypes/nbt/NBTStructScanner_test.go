package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type helloWorldStruct struct {
	HelloWorld struct {
		Name string `nbt:"name"`
	} `nbt:"hello world"`
}

func TestNBTStructScanner(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/nbt-test/hello_world.nbt")
	assert.Nil(t, exception)

	compound, _ := ReadNBT(inData)

	helloStruct := helloWorldStruct{}

	NBTStructScan(&helloStruct, &compound)

	assert.Equal(t, "Bananrama", helloStruct.HelloWorld.Name)
}

func TestNBTMarshal(t *testing.T) {

	helloStruct := helloWorldStruct{struct {
		Name string `nbt:"name"`
	}(struct{ Name string }{Name: "Bananrama"})}

	output := NBTMarshal(&helloStruct)

	inData, exception := ioutil.ReadFile("../../data/nbt-test/hello_world.nbt")
	assert.Nil(t, exception)

	assert.Equal(t, inData, output[:len(inData)])

}

type listTest struct {
	CompoundList []struct {
		Key   string
		Value string
	}

	StringArrayTest []string

	IntArrayTest  []int32
	LongArrayTest []int64
}

func TestNBTMarshalList(t *testing.T) {
	listTestVal := listTest{
		CompoundList: []struct {
			Key   string
			Value string
		}{{
			Key:   "Hello1",
			Value: "Hello2",
		}, {
			Key:   "Hello3",
			Value: "Hello4",
		}},
		StringArrayTest: []string{"Hello", "world", "reflected", "lists!"},
		IntArrayTest:    []int32{1, 2, 3, 4},
		LongArrayTest:   []int64{4, 3, 2, 1},
	}

	output := NBTMarshal(&listTestVal)

	rereadNbt, _ := ReadNBT(output)

	listTestReread := listTest{}

	NBTStructScan(&listTestReread, &rereadNbt)

	assert.Equal(t, listTestVal, listTestReread)
}

func TestLoadCodecNBT(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/codec.nbt")
	assert.Nil(t, exception)

	compressed := bytes.NewReader(inData)
	zr, exception := gzip.NewReader(compressed)
	assert.Nil(t, exception)

	uncompressed, exception := ioutil.ReadAll(zr)
	assert.Nil(t, exception)

	compound, _ := ReadNBT(uncompressed)

	codec := dataTypes.CodecOuterCompound{}

	NBTStructScan(&codec, &compound)

	data, _ := json.Marshal(codec)
	fmt.Println(string(data))
}
