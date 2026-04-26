package w3ggo_test

import (
	"strings"
	"testing"

	w3g "w3ggo"
)

func TestReplay132_Reforged1(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged1.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6091)
	assertEqual(t, "players len", len(result.Players), 2)
}

func TestReplay132_Reforged2(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged2.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6091)
	assertEqual(t, "players len", len(result.Players), 2)
}

func TestReplay132_Reforged2010(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged2010.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6102)
	assertEqual(t, "players len", len(result.Players), 6)
	if findPlayer(result.Players, "BEARAND#1604") == nil {
		t.Error("BEARAND#1604 not found")
	}
}

func TestReplay132_ReforgedRelease(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged_release.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6105)
	assertEqual(t, "players len", len(result.Players), 2)
	if findPlayer(result.Players, "anXieTy#2932") == nil {
		t.Error("anXieTy#2932 not found")
	}
	if findPlayer(result.Players, "IroNSoul#22724") == nil {
		t.Error("IroNSoul#22724 not found")
	}
}

func TestReplay132_ReforgedHunter2PrivateString(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged_hunter2_privatestring.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6105)
	assertEqual(t, "players len", len(result.Players), 2)
	if findPlayer(result.Players, "pischner#2950") == nil {
		t.Error("pischner#2950 not found")
	}
	if findPlayer(result.Players, "Wartoni#2638") == nil {
		t.Error("Wartoni#2638 not found")
	}
}

func TestReplay132_Netease(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "netease_132.nwg"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6105)
	assertEqual(t, "players len", len(result.Players), 2)
	if findPlayer(result.Players, "HurricaneBo") == nil {
		t.Error("HurricaneBo not found")
	}
	if findPlayer(result.Players, "SimplyHunteR") == nil {
		t.Error("SimplyHunteR not found")
	}
}

func TestReplay132_ReforgedTruncatedPlayernames(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged_truncated_playernames.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.32")
	assertEqual(t, "build_number", result.BuildNumber, 6105)
	assertEqual(t, "players len", len(result.Players), 2)
	if findPlayer(result.Players, "WaN#1734") == nil {
		t.Error("WaN#1734 not found")
	}
	found := false
	for _, p := range result.Players {
		if strings.Contains(p.Name, "1734") || strings.Contains(p.Name, "228941") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a player name containing '1734' or '228941'")
	}
}

func TestReplay132_RandomHeroRandomRaces(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "replay_randomhero_randomraces.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !result.Settings.RandomHero {
		t.Error("expected random_hero to be true")
	}
	if !result.Settings.RandomRaces {
		t.Error("expected random_races to be true")
	}
}

func TestReplay132_TeamsTogetherSettings(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "replay_teamstogether.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !result.Settings.FullSharedUnitControl {
		t.Error("expected full_shared_unit_control to be true")
	}
	if !result.Settings.TeamsTogether {
		t.Error("expected teams_together to be true")
	}
	if !result.Settings.FixedTeams {
		t.Error("expected fixed_teams to be true")
	}
	if result.Settings.RandomHero {
		t.Error("expected random_hero to be false")
	}
	if result.Settings.RandomRaces {
		t.Error("expected random_races to be false")
	}
}

func TestReplay132_FullObservers(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "replay_fullobs.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "settings.observer_mode", result.Settings.ObserverMode, w3g.ObserverModeFull)
}

func TestReplay132_Referees(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "replay_referee.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "settings.observer_mode", result.Settings.ObserverMode, w3g.ObserverModeReferees)
}

func TestReplay132_ObsOnDefeat(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "replay_obs_on_defeat.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "settings.observer_mode", result.Settings.ObserverMode, w3g.ObserverModeOnDefeat)
}

func TestReplay132_Hotkeys(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "reforged1.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	p0 := result.Players[0]
	p1 := result.Players[1]

	hk1p0 := p0.GroupHotkeys[1]
	assertEqual(t, "p0.hk1.assigned", hk1p0.Assigned, 1)
	assertEqual(t, "p0.hk1.used", hk1p0.Used, 29)

	hk2p0 := p0.GroupHotkeys[2]
	assertEqual(t, "p0.hk2.assigned", hk2p0.Assigned, 1)
	assertEqual(t, "p0.hk2.used", hk2p0.Used, 60)

	hk1p1 := p1.GroupHotkeys[1]
	assertEqual(t, "p1.hk1.assigned", hk1p1.Assigned, 21)
	assertEqual(t, "p1.hk1.used", hk1p1.Used, 106)

	hk2p1 := p1.GroupHotkeys[2]
	assertEqual(t, "p1.hk2.assigned", hk2p1.Assigned, 4)
	assertEqual(t, "p1.hk2.used", hk2p1.Used, 64)
}

func TestReplay132_KotGLevel6(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "706266088.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	var kotgPlayer *w3g.PlayerOutput
	for i := range result.Players {
		if findHero(result.Players[i].Heroes, "Ekee") != nil {
			kotgPlayer = &result.Players[i]
			break
		}
	}
	if kotgPlayer == nil {
		t.Fatal("player with KotG not found")
	}
	kotg := findHero(kotgPlayer.Heroes, "Ekee")
	assertEqual(t, "kotg.level", kotg.Level, 6)
}

func TestReplay132_Winner1640262494(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "1640262494.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 0)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "Happie")
}

func TestReplay132_Winner1448202825(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "1448202825.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "ThundeR#31281")
}

func TestReplay132_WinnerWanVsTrunks(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "wan_vs_trunks.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 0)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "WaN#1734")
}

func TestReplay132_WinnerBenjiii(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "benjiii_vs_Scars_Concealed_Hill.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "benjiii#1588")
}

func TestReplay132_WinnerESLCupChanger(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "esl_cup_vs_changer_1.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 0)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "TapioN#2351")
}

func TestReplay132_WinnerBuildingWinAnxiety(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "buildingwin_anxietyperspective.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "anXieTy#2932")
}

func TestReplay132_WinnerBuildingWinHelpstone(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("132", "buildingwin_helpstoneperspective.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "winning_team_id", result.WinningTeamID, 1)
	winner := findPlayerByTeam(result.Players, result.WinningTeamID)
	if winner == nil {
		t.Fatal("winner not found")
	}
	assertEqual(t, "winner.name", winner.Name, "anXieTy#2932")
}

func findPlayerByTeam(players []w3g.PlayerOutput, teamID int) *w3g.PlayerOutput {
	for i := range players {
		if players[i].TeamID == teamID {
			return &players[i]
		}
	}
	return nil
}
