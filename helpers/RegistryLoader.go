package helpers

import (
	"os"
	"path/filepath"
	"strings"
)

func LoadAllRegistries() map[string]map[string]string {

	output := make(map[string]map[string]string)

	filepath.Walk("data/minecraft", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
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

		output[registryName][entryName] = string(jsonData)

		return nil
	})
	return output
}
