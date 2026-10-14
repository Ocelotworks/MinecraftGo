package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func LoadAllRegistries() map[string]map[string]string {
	output := make(map[string]map[string]string)

	filepath.Walk("data/minecraft", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".json") || strings.HasPrefix(path, "data/minecraft/tags") {
			return nil
		}

		cleanPath, _ := strings.CutPrefix(strings.ReplaceAll(path, "\\", "/"), "data/minecraft/")

		dir, file := filepath.Split(cleanPath)

		// ignore meta files and non-json
		if strings.HasPrefix(file, "_") {
			return nil
		}

		registryName, _ := strings.CutSuffix(dir, "/")

		registryName = "minecraft:" + registryName

		_, exists := output[registryName]

		if !exists {
			output[registryName] = make(map[string]string)
		}

		jsonData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		entryName, _ := strings.CutSuffix(file, ".json")
		output[registryName]["minecraft:"+entryName] = string(jsonData)

		return nil
	})

	return output
}

type registryTag struct {
	Values []string `json:"values"`
}

func LoadAllTags() map[string]map[string][]string {
	output := make(map[string]map[string][]string)

	err := filepath.Walk("data/minecraft/tags", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		cleanPath, _ := strings.CutPrefix(strings.ReplaceAll(path, "\\", "/"), "data/minecraft/tags/")

		// First path component is the registry name; the rest (minus .json) is the tag
		// name.  filepath.Split can't be used here because tags like banner_pattern's
		// "pattern_item/flower" have nested paths: the registry is only the first segment.
		slashIdx := strings.Index(cleanPath, "/")
		if slashIdx == -1 {
			return nil
		}
		registryName := "minecraft:" + cleanPath[:slashIdx]
		tagName, _ := strings.CutSuffix(cleanPath[slashIdx+1:], ".json")

		// ignore meta files (prefixed with _)
		if strings.HasPrefix(filepath.Base(tagName), "_") {
			return nil
		}

		jsonData, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tagData := registryTag{}
		err = json.Unmarshal(jsonData, &tagData)
		if err != nil {
			return err
		}
		_, ok := output[registryName]
		if !ok {
			output[registryName] = make(map[string][]string)
		}
		output[registryName][tagName] = tagData.Values
		return nil
	})

	if err != nil {
		panic(err)
	}
	//fmt.Println(output)
	return output
}
