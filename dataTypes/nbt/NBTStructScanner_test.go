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

type compoundMapTest struct {
	CompoundMapString map[string]interface{} `nbt:"Level"`
}

func TestCompoundMapTest(t *testing.T) {
	inData, exception := ioutil.ReadFile("../../data/nbt-test/bigtest.nbt")
	assert.Nil(t, exception)

	compressed := bytes.NewReader(inData)
	zr, exception := gzip.NewReader(compressed)
	assert.Nil(t, exception)

	uncompressed, exception := ioutil.ReadAll(zr)
	assert.Nil(t, exception)

	compoundTestVal := compoundMapTest{}
	compound, _ := ReadNBT(uncompressed)
	NBTStructScan(&compoundTestVal, &compound)
	out, _ := json.Marshal(compoundTestVal)
	fmt.Println(string(out))
}
