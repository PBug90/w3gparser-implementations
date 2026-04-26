package w3ggo_test

import (
	"math"
	"testing"

	w3g "w3ggo"
)

func TestReplay129_NeteaseObs(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("129", "netease_129_obs.nwg"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.29")

	rudan := findPlayer(result.Players, "rudan")
	if rudan == nil {
		t.Fatal("rudan not found")
	}
	assertEqual(t, "rudan.color", rudan.Color, "#282828")

	assertEqual(t, "observers len", len(result.Observers), 1)
	assertEqual(t, "matchup", result.Matchup, "NvN")
	assertEqual(t, "game_type", result.Type, "1on1")
	assertEqual(t, "players len", len(result.Players), 2)

	assertEqual(t, "map.checksum", result.Map.Checksum, "281f9d6a")
	assertEqual(t, "map.checksum_sha1", result.Map.ChecksumSha1, "c232d68286eb4604cc66db42d45e28017b78e3c4")
	assertEqual(t, "map.file", result.Map.File, "(4)TurtleRock.w3x")
	assertEqual(t, "map.path", result.Map.Path, "Maps/1.29\\(4)TurtleRock.w3x")
}

func TestReplay129_StandardObs(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("129", "standard_129_obs.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.29")
	assertEqual(t, "players len", len(result.Players), 2)
	assertEqual(t, "matchup", result.Matchup, "OvO")
	assertEqual(t, "game_type", result.Type, "1on1")
	assertEqual(t, "observers len", len(result.Observers), 4)
	if len(result.Chat) <= 2 {
		t.Errorf("expected more than 2 chat messages, got %d", len(result.Chat))
	}

	sokol := findPlayer(result.Players, "S.o.K.o.L")
	if sokol == nil {
		t.Fatal("S.o.K.o.L not found")
	}
	assertEqual(t, "sokol.race_detected", sokol.RaceDetected, "O")
	assertEqual(t, "sokol.id", sokol.ID, 4)
	assertEqual(t, "sokol.teamid", sokol.TeamID, 3)
	assertEqual(t, "sokol.color", sokol.Color, "#00781e")
	assertEqual(t, "sokol.units.opeo", sokol.Units.SummaryMap["opeo"], 10)
	assertEqual(t, "sokol.units.ogru", sokol.Units.SummaryMap["ogru"], 5)
	assertEqual(t, "sokol.units.orai", sokol.Units.SummaryMap["orai"], 6)
	assertEqual(t, "sokol.units.ospm", sokol.Units.SummaryMap["ospm"], 5)
	assertEqual(t, "sokol.units.okod", sokol.Units.SummaryMap["okod"], 2)
	assertEqual(t, "sokol.actions.assign_group", sokol.Actions.AssignGroup, 38)
	assertEqual(t, "sokol.actions.right_click", sokol.Actions.RightClick, 1104)
	assertEqual(t, "sokol.actions.basic", sokol.Actions.Basic, 122)
	assertEqual(t, "sokol.actions.build_train", sokol.Actions.BuildTrain, 111)
	assertEqual(t, "sokol.actions.ability", sokol.Actions.Ability, 59)
	assertEqual(t, "sokol.actions.item", sokol.Actions.Item, 6)
	assertEqual(t, "sokol.actions.select", sokol.Actions.Select, 538)
	assertEqual(t, "sokol.actions.remove_unit", sokol.Actions.RemoveUnit, 0)
	assertEqual(t, "sokol.actions.select_hotkey", sokol.Actions.SelectHotkey, 751)
	assertEqual(t, "sokol.actions.esc", sokol.Actions.ESC, 0)

	stormhoof := findPlayer(result.Players, "Stormhoof")
	if stormhoof == nil {
		t.Fatal("Stormhoof not found")
	}
	assertEqual(t, "stormhoof.race_detected", stormhoof.RaceDetected, "O")
	assertEqual(t, "stormhoof.color", stormhoof.Color, "#9b0000")
	assertEqual(t, "stormhoof.id", stormhoof.ID, 6)
	assertEqual(t, "stormhoof.teamid", stormhoof.TeamID, 0)
	assertEqual(t, "stormhoof.units.opeo", stormhoof.Units.SummaryMap["opeo"], 11)
	assertEqual(t, "stormhoof.units.ogru", stormhoof.Units.SummaryMap["ogru"], 8)
	assertEqual(t, "stormhoof.units.orai", stormhoof.Units.SummaryMap["orai"], 8)
	assertEqual(t, "stormhoof.units.ospm", stormhoof.Units.SummaryMap["ospm"], 4)
	assertEqual(t, "stormhoof.units.okod", stormhoof.Units.SummaryMap["okod"], 3)
	assertEqual(t, "stormhoof.actions.assign_group", stormhoof.Actions.AssignGroup, 111)
	assertEqual(t, "stormhoof.actions.right_click", stormhoof.Actions.RightClick, 1595)
	assertEqual(t, "stormhoof.actions.basic", stormhoof.Actions.Basic, 201)
	assertEqual(t, "stormhoof.actions.build_train", stormhoof.Actions.BuildTrain, 112)
	assertEqual(t, "stormhoof.actions.ability", stormhoof.Actions.Ability, 57)
	assertEqual(t, "stormhoof.actions.item", stormhoof.Actions.Item, 5)
	assertEqual(t, "stormhoof.actions.select", stormhoof.Actions.Select, 653)
	assertEqual(t, "stormhoof.actions.remove_unit", stormhoof.Actions.RemoveUnit, 0)
	assertEqual(t, "stormhoof.actions.select_hotkey", stormhoof.Actions.SelectHotkey, 1865)
	assertEqual(t, "stormhoof.actions.esc", stormhoof.Actions.ESC, 4)

	assertEqual(t, "map.checksum", result.Map.Checksum, "008ab7f1")
	assertEqual(t, "map.checksum_sha1", result.Map.ChecksumSha1, "79ba7579f28e5ccfd741a1ebfbff95a56813086e")
	assertEqual(t, "map.file", result.Map.File, "w3arena__twistedmeadows__v3.w3x")
	assertEqual(t, "map.path", result.Map.Path, "Maps\\w3arena\\w3arena__twistedmeadows__v3.w3x")
}

func TestReplay129_3on3LeaverAPM(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("129", "standard_129_3on3_leaver.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	abmit := findPlayer(result.Players, "abmitdirpic")
	if abmit == nil {
		t.Fatal("abmitdirpic not found")
	}
	firstLeftMinute := int(math.Ceil(float64(abmit.CurrentTimePlayed) / 1000.0 / 60.0))
	var postLeaveSum int
	for _, v := range abmit.Actions.Timed[firstLeftMinute:] {
		postLeaveSum += v
	}
	assertEqual(t, "post_leave_sum", postLeaveSum, 0)
	assertEqual(t, "abmit.apm", abmit.APM, 98)
	assertEqual(t, "abmit.current_time_played", abmit.CurrentTimePlayed, 4371069)
}
