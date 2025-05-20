package controller

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"io/ioutil"
)

type DataStore struct {
	BlockData map[string]entity.BlockData
	Map       [][]*dataTypes.RegionMetadata
	Codec     *dataTypes.CodecOuterCompound
}

func NewDataStore() *DataStore {
	ds := &DataStore{}

	ds.LoadBlockData()
	ds.LoadStartingArea()
	ds.LoadCodec()

	return ds
}

func (ds *DataStore) LoadBlockData() {
	fmt.Println("Loading BlockData...")
	blockFile, exception := ioutil.ReadFile("data/blocks.json")
	exception = json.Unmarshal(blockFile, &ds.BlockData)

	if exception != nil {
		fmt.Println("Error reading block data", exception)
	}
}

func (ds *DataStore) LoadStartingArea() {
	fmt.Println("Loading starting region...")
	inData, exception := ioutil.ReadFile("C:\\Users\\unacc\\AppData\\Roaming\\.minecraft\\saves\\flat\\region\\r.0.0.mca")

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return
	}

	region := dataTypes.ReadRegionFile(inData)

	ds.Map = [][]*dataTypes.RegionMetadata{{&region}}
}

func (ds *DataStore) LoadCodec() {
	fmt.Println("Loading codec...")
	inData, exception := ioutil.ReadFile("data/codec.nbt")

	compressed := bytes.NewReader(inData)
	zr, exception := gzip.NewReader(compressed)

	if exception != nil {
		fmt.Println(exception)
	}

	uncompressed, exception := ioutil.ReadAll(zr)

	if exception != nil {
		fmt.Println(exception)
	}

	compound, _ := nbt.ReadNBT(uncompressed)
	ds.Codec = &dataTypes.CodecOuterCompound{}

	nbt.NBTStructScan(ds.Codec, &compound)
}
