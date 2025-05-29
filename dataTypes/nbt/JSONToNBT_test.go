package nbt

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/helpers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJSONToNBT(t *testing.T) {

	t.Run("Simple JSON", func(t *testing.T) {
		compoundValue, err := JSONToNBT([]byte(`{
  "effects": "burning",
  "exhaustion": 0.1,
  "message_id": "inFire",
  "scaling": "when_caused_by_living_non_player"
}`))

		assert.NoError(t, err)

		type DamageTypeNBT struct {
			Effects    string  `nbt:"effects"`
			Exhaustion float64 `nbt:"exhaustion"`
			MessageId  string  `nbt:"message_id"`
			Scaling    string  `nbt:"scaling"`
		}

		damageType := DamageTypeNBT{}
		NBTStructScan(&damageType, &compoundValue)

		assert.Equal(t, "burning", damageType.Effects)
		assert.Equal(t, 0.1, damageType.Exhaustion)
		assert.Equal(t, "inFire", damageType.MessageId)
		assert.Equal(t, "when_caused_by_living_non_player", damageType.Scaling)
	})

	t.Run("Nested compounds", func(t *testing.T) {
		compoundValue, err := JSONToNBT([]byte(`{
  "assets": {
    "angry": "minecraft:entity/wolf/wolf_angry",
    "tame": "minecraft:entity/wolf/wolf_tame",
    "wild": "minecraft:entity/wolf/wolf"
  },
  "spawn_conditions": [
    {
      "priority": 0
    }
  ]
}`))
		type WolfVariant struct {
			Assets struct {
				Angry string `nbt:"angry"`
				Tame  string `nbt:"tame"`
				Wild  string `nbt:"wild"`
			} `nbt:"assets"`
			SpawnConditions []struct {
				Priority float64 `nbt:"priority"`
			} `nbt:"spawn_conditions"`
		}

		wolfVariant := WolfVariant{}

		NBTStructScan(&wolfVariant, &compoundValue)
		fmt.Println(wolfVariant)

		assert.Equal(t, "minecraft:entity/wolf/wolf_angry", wolfVariant.Assets.Angry)
		assert.Equal(t, "minecraft:entity/wolf/wolf_tame", wolfVariant.Assets.Tame)
		assert.Equal(t, "minecraft:entity/wolf/wolf", wolfVariant.Assets.Wild)
		assert.Equal(t, float64(0), wolfVariant.SpawnConditions[0].Priority)

		assert.NoError(t, err)
	})

	t.Run("Flowering Azalea Leaves", func(t *testing.T) {
		f, _ := os.ReadFile("data/minecraft/loot_table/blocks/flowering_azalea_leaves.json")
		JSONToNBT(f)
	})

	t.Run("Parse all registries", func(t *testing.T) {
		registries := helpers.LoadAllRegistries()

		for registryName, registryEntries := range registries {
			for entryFile, entryData := range registryEntries {
				t.Run(registryName+"/"+entryFile, func(t *testing.T) {
					JSONToNBT([]byte(entryData))
				})
			}
		}
	})

}
