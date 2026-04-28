package w3ggo_test

import (
	"path/filepath"
	"runtime"
	"testing"

	w3g "github.com/PBug90/w3gparser-implementations/w3ggo"
)

func replayPath(version, name string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "testdata", "replays", version, name)
}

func TestReplay126_2on2Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("126", "999.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.26")
	assertEqual(t, "players len", len(result.Players), 4)

	p0 := result.Players[0]
	assertEqual(t, "p0.ID", p0.ID, 2)
	assertEqual(t, "p0.TeamID", p0.TeamID, 0)

	p1 := result.Players[1]
	assertEqual(t, "p1.ID", p1.ID, 4)
	assertEqual(t, "p1.TeamID", p1.TeamID, 0)

	p2 := result.Players[2]
	assertEqual(t, "p2.ID", p2.ID, 3)
	assertEqual(t, "p2.TeamID", p2.TeamID, 1)

	p3 := result.Players[3]
	assertEqual(t, "p3.ID", p3.ID, 5)
	assertEqual(t, "p3.TeamID", p3.TeamID, 1)

	assertEqual(t, "matchup", result.Matchup, "HUvHU")
	assertEqual(t, "game_type", result.Type, "2on2")
	assertEqual(t, "map.checksum", result.Map.Checksum, "b4230d1e")
	assertEqual(t, "map.checksum_sha1", result.Map.ChecksumSha1, "1f75e2a24fd995a6d7b123bb44d8afae7b5c6222")
	assertEqual(t, "map.file", result.Map.File, "w3arena__maelstrom__v2.w3x")
	assertEqual(t, "map.path", result.Map.Path, "Maps\\w3arena\\w3arena__maelstrom__v2.w3x")
}

func TestReplay126_Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("126", "standard_126.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.26")
	assertEqual(t, "observers len", len(result.Observers), 8)
	assertEqual(t, "players len", len(result.Players), 2)
	assertEqual(t, "matchup", result.Matchup, "HvU")
	assertEqual(t, "game_type", result.Type, "1on1")

	happy := findPlayer(result.Players, "Happy_")
	if happy == nil {
		t.Fatal("Happy_ not found")
	}
	assertEqual(t, "happy.race_detected", happy.RaceDetected, "U")
	assertEqual(t, "happy.color", happy.Color, "#0042ff")

	u2 := findPlayer(result.Players, "u2.sok")
	if u2 == nil {
		t.Fatal("u2.sok not found")
	}
	assertEqual(t, "u2.race_detected", u2.RaceDetected, "H")
	assertEqual(t, "u2.color", u2.Color, "#ff0303")

	assertEqual(t, "map.checksum", result.Map.Checksum, "51a1c63b")
	assertEqual(t, "map.checksum_sha1", result.Map.ChecksumSha1, "0b4f05ca7dcc23b9501422b4fa26a86c7d2a0ee0")
	assertEqual(t, "map.file", result.Map.File, "w3arena__amazonia__v3.w3x")
	assertEqual(t, "map.path", result.Map.Path, "Maps\\w3arena\\w3arena__amazonia__v3.w3x")
}

// helpers

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func findPlayer(players []w3g.PlayerOutput, name string) *w3g.PlayerOutput {
	for i := range players {
		if players[i].Name == name {
			return &players[i]
		}
	}
	return nil
}
