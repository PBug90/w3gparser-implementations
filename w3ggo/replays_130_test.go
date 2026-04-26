package w3ggo_test

import (
	"testing"

	w3g "w3ggo"
)

func TestReplay1302_Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("130", "standard_1302.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.30.2+")
	assertEqual(t, "matchup", result.Matchup, "NvU")
	assertEqual(t, "players len", len(result.Players), 2)
}

func TestReplay1303_Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("130", "standard_1303.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.30.2+")
	assertEqual(t, "players len", len(result.Players), 2)
}

func TestReplay1304_Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("130", "standard_1304.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.30.2+")
	assertEqual(t, "players len", len(result.Players), 2)
}

func TestReplay1304_2on2(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("130", "standard_1304.2on2.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.30.2+")
	assertEqual(t, "build_number", result.BuildNumber, 6061)
	assertEqual(t, "players len", len(result.Players), 4)
}

func TestReplay130_Standard(t *testing.T) {
	result, err := w3g.ParseFile(replayPath("130", "standard_130.w3g"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	assertEqual(t, "version", result.Version, "1.30")
	assertEqual(t, "matchup", result.Matchup, "NvU")
	assertEqual(t, "game_type", result.Type, "1on1")
	assertEqual(t, "players len", len(result.Players), 2)

	sheik := findPlayer(result.Players, "sheik")
	if sheik == nil {
		t.Fatal("sheik not found")
	}
	assertEqual(t, "sheik.race", sheik.Race, "U")
	assertEqual(t, "sheik.race_detected", sheik.RaceDetected, "U")

	other := findPlayer(result.Players, "123456789012345")
	if other == nil {
		t.Fatal("123456789012345 not found")
	}
	assertEqual(t, "other.race", other.Race, "N")
	assertEqual(t, "other.race_detected", other.RaceDetected, "N")

	// Heroes for sheik
	hero0 := findHero(sheik.Heroes, "Udea")
	if hero0 == nil {
		t.Fatal("Udea hero not found")
	}
	assertEqual(t, "hero0.level", hero0.Level, 6)

	hero1 := findHero(sheik.Heroes, "Ulic")
	if hero1 == nil {
		t.Fatal("Ulic hero not found")
	}
	assertEqual(t, "hero1.level", hero1.Level, 6)

	hero2 := findHero(sheik.Heroes, "Udre")
	if hero2 == nil {
		t.Fatal("Udre hero not found")
	}
	assertEqual(t, "hero2.level", hero2.Level, 3)

	assertEqual(t, "map.file", result.Map.File, "(4)TwistedMeadows.w3x")
	assertEqual(t, "map.checksum", result.Map.Checksum, "c3cae01d")
	assertEqual(t, "map.checksum_sha1", result.Map.ChecksumSha1, "23dc614cca6fd7ec232fbba4898d318a90b95bc6")
	assertEqual(t, "map.path", result.Map.Path, "Maps\\FrozenThrone\\(4)TwistedMeadows.w3x")
}

func TestReplay130_ResetElapsedMs(t *testing.T) {
	r1, err := w3g.ParseFile(replayPath("130", "standard_130.w3g"))
	if err != nil {
		t.Fatalf("parse 1 failed: %v", err)
	}
	r2, err := w3g.ParseFile(replayPath("130", "standard_130.w3g"))
	if err != nil {
		t.Fatalf("parse 2 failed: %v", err)
	}
	assertEqual(t, "duration", r1.Duration, r2.Duration)
}

func findHero(heroes []w3g.HeroInfo, id string) *w3g.HeroInfo {
	for i := range heroes {
		if heroes[i].ID == id {
			return &heroes[i]
		}
	}
	return nil
}
