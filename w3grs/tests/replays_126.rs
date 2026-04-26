use w3grs::parse_file;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_126_2on2_standard() {
    let result = parse_file(&replay_path("126", "999.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.26");
    assert_eq!(result.players.len(), 4);

    let p0 = &result.players[0];
    assert_eq!(p0.id, 2);
    assert_eq!(p0.teamid, 0);

    let p1 = &result.players[1];
    assert_eq!(p1.id, 4);
    assert_eq!(p1.teamid, 0);

    let p2 = &result.players[2];
    assert_eq!(p2.id, 3);
    assert_eq!(p2.teamid, 1);

    let p3 = &result.players[3];
    assert_eq!(p3.id, 5);
    assert_eq!(p3.teamid, 1);

    assert_eq!(result.matchup, "HUvHU");
    assert_eq!(result.game_type, "2on2");

    assert_eq!(result.map.checksum, "b4230d1e");
    assert_eq!(result.map.checksum_sha1, "1f75e2a24fd995a6d7b123bb44d8afae7b5c6222");
    assert_eq!(result.map.file, "w3arena__maelstrom__v2.w3x");
    assert_eq!(result.map.path, "Maps\\w3arena\\w3arena__maelstrom__v2.w3x");
}

#[test]
fn test_126_standard() {
    let result = parse_file(&replay_path("126", "standard_126.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.26");
    assert_eq!(result.observers.len(), 8);
    assert_eq!(result.players.len(), 2);
    assert_eq!(result.matchup, "HvU");
    assert_eq!(result.game_type, "1on1");

    // players are sorted by teamid then id
    let happy = result.players.iter().find(|p| p.name == "Happy_").expect("Happy_ not found");
    assert_eq!(happy.race_detected, "U");
    assert_eq!(happy.color, "#0042ff");

    let u2 = result.players.iter().find(|p| p.name == "u2.sok").expect("u2.sok not found");
    assert_eq!(u2.race_detected, "H");
    assert_eq!(u2.color, "#ff0303");

    assert_eq!(result.map.checksum, "51a1c63b");
    assert_eq!(result.map.checksum_sha1, "0b4f05ca7dcc23b9501422b4fa26a86c7d2a0ee0");
    assert_eq!(result.map.file, "w3arena__amazonia__v3.w3x");
    assert_eq!(result.map.path, "Maps\\w3arena\\w3arena__amazonia__v3.w3x");
}
