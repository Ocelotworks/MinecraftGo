package constants

import "testing"

// Packet ids are handed out per state and flow in alphabetical order of the
// packet's registry name, so a single packet added upstream renumbers every
// packet after it.  These pin the ids the join sequence depends on: if one of
// them fails after a version bump, the fix is to regenerate rather than to
// edit the expectation by hand.
//
//	go run ./cmd/mcdata -jar <client jar>
func TestJoinSequencePacketIDs(t *testing.T) {
	cases := []struct {
		name  string
		got   int
		want  int
		about string
	}{
		{"CBLoginFinished", CBLoginFinished, 0x02, "login/clientbound minecraft:login_finished"},
		{"CBRegistryDataConfiguration", CBRegistryDataConfiguration, 0x07, "configuration/clientbound minecraft:registry_data"},
		{"CBUpdateTagsConfiguration", CBUpdateTagsConfiguration, 0x0E, "configuration/clientbound minecraft:update_tags"},
		{"CBFinishConfiguration", CBFinishConfiguration, 0x03, "configuration/clientbound minecraft:finish_configuration"},
		{"CBLogin", CBLogin, 0x32, "play/clientbound minecraft:login — the join game packet"},
		{"CBSetChunkCacheCenter", CBSetChunkCacheCenter, 0x60, "play/clientbound minecraft:set_chunk_cache_center"},
		{"CBLevelChunkWithLight", CBLevelChunkWithLight, 0x2E, "play/clientbound minecraft:level_chunk_with_light"},
		{"CBLightUpdate", CBLightUpdate, 0x31, "play/clientbound minecraft:light_update"},
		{"CBSetDefaultSpawnPosition", CBSetDefaultSpawnPosition, 0x63, "play/clientbound minecraft:set_default_spawn_position"},
		{"CBPlayerPosition", CBPlayerPosition, 0x49, "play/clientbound minecraft:player_position"},
		{"CBSetHeldSlot", CBSetHeldSlot, 0x6B, "play/clientbound minecraft:set_held_slot"},
		{"SBClientInformation", SBClientInformation, 0x0E, "play/serverbound minecraft:client_information"},
		{"SBAcceptTeleportation", SBAcceptTeleportation, 0x00, "play/serverbound minecraft:accept_teleportation"},
		{"SBKeepAlive", SBKeepAlive, 0x1C, "play/serverbound minecraft:keep_alive"},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = 0x%02X, want 0x%02X (%s)", c.name, c.got, c.want, c.about)
		}
	}
}

// The states overlap: keep_alive, disconnect and custom_payload all exist in
// both configuration and play under different ids, so the constants have to
// stay distinct per state.
func TestStatesDoNotShareConstants(t *testing.T) {
	if CBKeepAliveConfiguration == CBKeepAlive {
		t.Error("configuration and play keep_alive share an id")
	}
	if SBCustomPayloadConfiguration == SBCustomPayload {
		t.Error("configuration and play custom_payload share an id")
	}
}
