package w3ggo_test

import (
	"testing"

	w3g "github.com/PBug90/w3gparser-implementations/w3ggo"
)

func TestReplay131_Action0x7a(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("131", "action0x7a.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.31")
	assertEqual(t, "players len", len(result.Players), 1)
	assertEqual(t, "winning_team_id", result.WinningTeamID, -1)
}

func TestReplay131_TomeOfRetraining(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("131", "standard_tomeofretraining_1.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.31")
	assertEqual(t, "build_number", result.BuildNumber, 6072)
	assertEqual(t, "players len", len(result.Players), 2)

	// Find a player with Hamg hero
	var playerWithHamg *w3g.PlayerOutput
	for i := range result.Players {
		if findHero(result.Players[i].Heroes, "Hamg") != nil {
			playerWithHamg = &result.Players[i]
			break
		}
	}
	if playerWithHamg == nil {
		t.Fatal("No player with Hamg hero")
	}

	hamg := findHero(playerWithHamg.Heroes, "Hamg")
	if hamg == nil {
		t.Fatal("Hamg hero not found")
	}
	assertEqual(t, "hamg.level", hamg.Level, 4)
	assertEqual(t, "hamg.abilities.AHab", hamg.Abilities["AHab"], 2)
	assertEqual(t, "hamg.abilities.AHbz", hamg.Abilities["AHbz"], 2)
	assertEqual(t, "hamg.retraining_history len", len(hamg.RetrainingHistory), 1)
	assertEqual(t, "hamg.retraining_history[0].time", hamg.RetrainingHistory[0].Time, 1136022)
	assertEqual(t, "hamg.retraining_history[0].abilities.AHab", hamg.RetrainingHistory[0].Abilities["AHab"], 2)
	assertEqual(t, "hamg.retraining_history[0].abilities.AHwe", hamg.RetrainingHistory[0].Abilities["AHwe"], 2)

	// Check ability order
	ao := hamg.AbilityOrder
	assertEqual(t, "ao len", len(ao), 9)

	assertAbility := func(i, wantTime int, wantValue string) {
		t.Helper()
		if i >= len(ao) {
			t.Errorf("ao[%d]: out of range", i)
			return
		}
		if ao[i].Type != "ability" {
			t.Errorf("ao[%d].Type: got %s, want ability", i, ao[i].Type)
		}
		if ao[i].Time != wantTime {
			t.Errorf("ao[%d].Time: got %d, want %d", i, ao[i].Time, wantTime)
		}
		if ao[i].Value != wantValue {
			t.Errorf("ao[%d].Value: got %s, want %s", i, ao[i].Value, wantValue)
		}
	}
	assertRetraining := func(i, wantTime int) {
		t.Helper()
		if i >= len(ao) {
			t.Errorf("ao[%d]: out of range", i)
			return
		}
		if ao[i].Type != "retraining" {
			t.Errorf("ao[%d].Type: got %s, want retraining", i, ao[i].Type)
		}
		if ao[i].Time != wantTime {
			t.Errorf("ao[%d].Time: got %d, want %d", i, ao[i].Time, wantTime)
		}
	}

	assertAbility(0, 124366, "AHwe")
	assertAbility(1, 234428, "AHab")
	assertAbility(2, 293007, "AHwe")
	assertAbility(3, 1060007, "AHab")
	assertRetraining(4, 1136022)
	assertAbility(5, 1140944, "AHbz")
	assertAbility(6, 1141147, "AHbz")
	assertAbility(7, 1141460, "AHab")
	assertAbility(8, 1141569, "AHab")
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "[OCG]shocker")
}

func TestReplay131_RocMapName(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("131", "roc-losttemple-mapname.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.31")
	assertEqual(t, "build_number", result.BuildNumber, 6072)
	assertEqual(t, "map.file", result.Map.File, "(4)LostTemple [Unforged 0.5 RoC].w3x")
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "syNtec")
}
