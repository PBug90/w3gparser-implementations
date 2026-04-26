package w3ggo_test

import (
	"testing"

	w3g "w3ggotest"
)

func TestReplay200_HauntedGoldMine(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "goldmine test.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	p := result.Players[0]
	assertEqual(t, "ugol count", p.Buildings.SummaryMap["ugol"], 1)
	assertEqual(t, "buildings.order len", len(p.Buildings.Order), 1)
	assertEqual(t, "buildings.order[0].id", p.Buildings.Order[0].ID, "ugol")
	assertEqual(t, "buildings.order[0].ms", p.Buildings.Order[0].MS, 28435)
}

func TestReplay200_Version(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "goldmine test.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "2.00")
}

func TestReplay200_CustomMapUIComponents(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "TempReplay.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "2.00")
}

func TestReplay200_Retraining(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "retrainingissues.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "2.00")

	// Find player with Archmage (Hamg) hero
	var player *w3g.PlayerOutput
	for i := range result.Players {
		if findHero(result.Players[i].Heroes, "Hamg") != nil {
			player = &result.Players[i]
			break
		}
	}
	if player == nil {
		t.Fatal("player with Archmage hero not found")
	}

	hamg := findHero(player.Heroes, "Hamg")
	if hamg == nil {
		t.Fatal("Hamg hero not found")
	}
	assertEqual(t, "hamg.level", hamg.Level, 6)

	ao := hamg.AbilityOrder
	hasRetraining := false
	for _, e := range ao {
		if e.Type == "retraining" {
			hasRetraining = true
			break
		}
	}
	if !hasRetraining {
		t.Error("expected ability order to contain a retraining entry")
	}
}

func TestReplay200_202MeleeChat(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "2.0.2-Melee.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(result.Chat) < 2 {
		t.Errorf("expected at least 2 chat messages, got %d", len(result.Chat))
		return
	}

	var msg0 *w3g.ChatMessage
	for i := range result.Chat {
		if result.Chat[i].Message == "don't hurt me" {
			msg0 = &result.Chat[i]
			break
		}
	}
	if msg0 == nil {
		t.Fatal("first chat message not found")
	}
	assertEqual(t, "msg0.player_id", msg0.PlayerID, 1)

	var msg1 *w3g.ChatMessage
	for i := range result.Chat {
		if result.Chat[i].Message == "no more" {
			msg1 = &result.Chat[i]
			break
		}
	}
	if msg1 == nil {
		t.Fatal("second chat message not found")
	}
	assertEqual(t, "msg1.player_id", msg1.PlayerID, 2)
}

func TestReplay200_202FloTVSavedByWc3(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("200", "2.0.2-FloTVSavedByWc3.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(result.Players) < 1 {
		t.Error("expected at least 1 player")
	}
}
