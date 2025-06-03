package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"io/ioutil"
)

type DataStore struct {
	BlockData map[string]entity.BlockData
	Map       [][]*dataTypes.RegionMetadata
}

func NewDataStore() *DataStore {
	ds := &DataStore{}

	ds.LoadBlockData()
	//ds.LoadStartingArea()

	return ds
}

func (ds *DataStore) LoadBlockData() {
	fmt.Println("Loading BlockData...")
	blockFile, exception := ioutil.ReadFile("data/blocks.json")

	blockData := make([]entity.BlockData, 0)

	exception = json.Unmarshal(blockFile, &blockData)

	if exception != nil {
		fmt.Println("Error reading block data", exception)
	}

	ds.BlockData = make(map[string]entity.BlockData)
	for _, block := range blockData {
		ds.BlockData["minecraft:"+block.Name] = block
	}

}

func (ds *DataStore) loadRegionFile(x int, y int) *dataTypes.RegionMetadata {
	inData, exception := ioutil.ReadFile(fmt.Sprintf("C:\\Users\\unacc\\IdeaProjects\\MinecraftGo\\data\\worlds\\MCGO_FlatTest\\region\\r.%d.%d.mca", x, y))

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return nil
	}

	region := dataTypes.ReadRegionFile(inData)
	return &region
}

func (ds *DataStore) LoadStartingArea() {
	fmt.Println("Loading starting region...")

	ds.Map = [][]*dataTypes.RegionMetadata{
		{ds.loadRegionFile(-1, -1)},
	}
}
