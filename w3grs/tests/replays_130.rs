use w3grs::parse_file;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_1302_standard() {
    let result = parse_file(&replay_path("130", "standard_1302.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.30.2+");
    assert_eq!(result.matchup, "NvU");
    assert_eq!(result.players.len(), 2);
}

#[test]
fn test_1303_standard() {
    let result = parse_file(&replay_path("130", "standard_1303.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.30.2+");
    assert_eq!(result.players.len(), 2);
}

#[test]
fn test_1304_standard() {
    let result = parse_file(&replay_path("130", "standard_1304.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.30.2+");
    assert_eq!(result.players.len(), 2);
}

#[test]
fn test_1304_2on2() {
    let result = parse_file(&replay_path("130", "standard_1304.2on2.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.30.2+");
    assert_eq!(result.build_number, 6061);
    assert_eq!(result.players.len(), 4);
}

#[test]
fn test_130_standard() {
    let result = parse_file(&replay_path("130", "standard_130.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.30");
    assert_eq!(result.matchup, "NvU");
    assert_eq!(result.game_type, "1on1");
    assert_eq!(result.players.len(), 2);

    let sheik = result.players.iter().find(|p| p.name == "sheik").expect("sheik not found");
    assert_eq!(sheik.race, "U");
    assert_eq!(sheik.race_detected, "U");

    let other = result.players.iter().find(|p| p.name == "123456789012345").expect("123456789012345 not found");
    assert_eq!(other.race, "N");
    assert_eq!(other.race_detected, "N");

    // Heroes for sheik
    let hero0 = sheik.heroes.iter().find(|h| h.id == "Udea").expect("Udea hero not found");
    assert_eq!(hero0.level, 6);
    let hero1 = sheik.heroes.iter().find(|h| h.id == "Ulic").expect("Ulic hero not found");
    assert_eq!(hero1.level, 6);
    let hero2 = sheik.heroes.iter().find(|h| h.id == "Udre").expect("Udre hero not found");
    assert_eq!(hero2.level, 3);

    assert_eq!(result.map.file, "(4)TwistedMeadows.w3x");
    assert_eq!(result.map.checksum, "c3cae01d");
    assert_eq!(result.map.checksum_sha1, "23dc614cca6fd7ec232fbba4898d318a90b95bc6");
    assert_eq!(result.map.path, "Maps\\FrozenThrone\\(4)TwistedMeadows.w3x");
}

#[test]
fn test_130_reset_elapsed_ms() {
    // Parse twice and ensure ms elapsed is the same
    let r1 = parse_file(&replay_path("130", "standard_130.w3g")).expect("parse failed");
    let r2 = parse_file(&replay_path("130", "standard_130.w3g")).expect("parse failed");
    assert_eq!(r1.duration, r2.duration);
}
