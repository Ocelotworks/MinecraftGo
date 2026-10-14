package controller

import (
	"net"
	"testing"
	"time"

	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/helpers/joinprobe"
	"github.com/Ocelotworks/MinecraftGo/helpers/testutil"
)

// Walks a real connection through handshake, login and configuration against the
// running server code and decodes the first chunk packet it receives, so the
// bytes checked are the ones that actually go on the wire.
//
// This checks the layout the test's own reader expects.  cmd/mcverify checks the
// same packets against the game's decoders, which is the check that matters when
// the two disagree.
func TestChunkPacketOverRealConnection(t *testing.T) {
	testutil.ChdirToRepoRoot(t) // data/ is relative to the repo root

	mc := CreateMinecraft()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		Init(conn, nil, mc)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	packets, err := joinprobe.Join(conn, "Tester", 30*time.Second, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// each phase has to have run its course
	seen := map[string]map[int]bool{}
	for _, packet := range packets {
		if seen[packet.Phase] == nil {
			seen[packet.Phase] = map[int]bool{}
		}
		seen[packet.Phase][packet.ID] = true
	}
	for phase, id := range map[string]int{
		joinprobe.PhaseLogin:         constants.CBLoginFinished,
		joinprobe.PhaseConfiguration: constants.CBFinishConfiguration,
		joinprobe.PhasePlay:          constants.CBLogin,
	} {
		if !seen[phase][id] {
			t.Errorf("no packet 0x%02X in the %s phase", id, phase)
		}
	}

	chunks := 0
	for _, packet := range packets {
		if packet.Phase == joinprobe.PhasePlay && packet.ID == constants.CBLevelChunkWithLight {
			chunks++
			decodeChunkPacket(t, packet.Body)
		}
	}
	if chunks == 0 {
		t.Error("no chunk packet arrived")
	}
	t.Logf("%d packets received, %d of them chunks", len(packets), chunks)
}
