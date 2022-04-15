package dataTypes

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRegionLoad(t *testing.T) {
	blockFile, exception := ioutil.ReadFile("../data/blocks.json")
	assert.NoError(t, exception)

	blockData := make(map[string]entity.BlockData)
	exception = json.Unmarshal(blockFile, &blockData)
	assert.NoError(t, exception)

	inData, exception := ioutil.ReadFile("../data/worlds/MCGO_FlatTest/region/r.0.0.mca")
	assert.NoError(t, exception)

	region := ReadRegionFile(inData, blockData)

	fmt.Println(region)

}
