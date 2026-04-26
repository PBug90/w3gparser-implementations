use w3grs::parse_file;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_129_netease_obs() {
    let result = parse_file(&replay_path("129", "netease_129_obs.nwg")).expect("parse failed");
    assert_eq!(result.version, "1.29");

    let rudan = result.players.iter().find(|p| p.name == "rudan").expect("rudan not found");
    assert_eq!(rudan.color, "#282828");

    assert_eq!(result.observers.len(), 1);
    assert_eq!(result.matchup, "NvN");
    assert_eq!(result.game_type, "1on1");
    assert_eq!(result.players.len(), 2);

    assert_eq!(result.map.checksum, "281f9d6a");
    assert_eq!(result.map.checksum_sha1, "c232d68286eb4604cc66db42d45e28017b78e3c4");
    assert_eq!(result.map.file, "(4)TurtleRock.w3x");
    assert_eq!(result.map.path, "Maps/1.29\\(4)TurtleRock.w3x");
}

#[test]
fn test_129_standard_obs() {
    let result = parse_file(&replay_path("129", "standard_129_obs.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.29");
    assert_eq!(result.players.len(), 2);
    assert_eq!(result.matchup, "OvO");
    assert_eq!(result.game_type, "1on1");
    assert_eq!(result.observers.len(), 4);
    assert!(result.chat.len() > 2);

    let sokol = result.players.iter().find(|p| p.name == "S.o.K.o.L").expect("S.o.K.o.L not found");
    assert_eq!(sokol.race_detected, "O");
    assert_eq!(sokol.id, 4);
    assert_eq!(sokol.teamid, 3);
    assert_eq!(sokol.color, "#00781e");
    assert_eq!(sokol.units.summary.get("opeo").copied().unwrap_or(0), 10);
    assert_eq!(sokol.units.summary.get("ogru").copied().unwrap_or(0), 5);
    assert_eq!(sokol.units.summary.get("orai").copied().unwrap_or(0), 6);
    assert_eq!(sokol.units.summary.get("ospm").copied().unwrap_or(0), 5);
    assert_eq!(sokol.units.summary.get("okod").copied().unwrap_or(0), 2);
    assert_eq!(sokol.actions.assign_group, 38);
    assert_eq!(sokol.actions.right_click, 1104);
    assert_eq!(sokol.actions.basic, 122);
    assert_eq!(sokol.actions.build_train, 111);
    assert_eq!(sokol.actions.ability, 59);
    assert_eq!(sokol.actions.item, 6);
    assert_eq!(sokol.actions.select, 538);
    assert_eq!(sokol.actions.remove_unit, 0);
    assert_eq!(sokol.actions.select_hotkey, 751);
    assert_eq!(sokol.actions.esc, 0);

    let stormhoof = result.players.iter().find(|p| p.name == "Stormhoof").expect("Stormhoof not found");
    assert_eq!(stormhoof.race_detected, "O");
    assert_eq!(stormhoof.color, "#9b0000");
    assert_eq!(stormhoof.id, 6);
    assert_eq!(stormhoof.teamid, 0);
    assert_eq!(stormhoof.units.summary.get("opeo").copied().unwrap_or(0), 11);
    assert_eq!(stormhoof.units.summary.get("ogru").copied().unwrap_or(0), 8);
    assert_eq!(stormhoof.units.summary.get("orai").copied().unwrap_or(0), 8);
    assert_eq!(stormhoof.units.summary.get("ospm").copied().unwrap_or(0), 4);
    assert_eq!(stormhoof.units.summary.get("okod").copied().unwrap_or(0), 3);
    assert_eq!(stormhoof.actions.assign_group, 111);
    assert_eq!(stormhoof.actions.right_click, 1595);
    assert_eq!(stormhoof.actions.basic, 201);
    assert_eq!(stormhoof.actions.build_train, 112);
    assert_eq!(stormhoof.actions.ability, 57);
    assert_eq!(stormhoof.actions.item, 5);
    assert_eq!(stormhoof.actions.select, 653);
    assert_eq!(stormhoof.actions.remove_unit, 0);
    assert_eq!(stormhoof.actions.select_hotkey, 1865);
    assert_eq!(stormhoof.actions.esc, 4);

    assert_eq!(result.map.checksum, "008ab7f1");
    assert_eq!(result.map.checksum_sha1, "79ba7579f28e5ccfd741a1ebfbff95a56813086e");
    assert_eq!(result.map.file, "w3arena__twistedmeadows__v3.w3x");
    assert_eq!(result.map.path, "Maps\\w3arena\\w3arena__twistedmeadows__v3.w3x");
}

#[test]
fn test_129_3on3_leaver_apm() {
    let result = parse_file(&replay_path("129", "standard_129_3on3_leaver.w3g")).expect("parse failed");

    let abmit = result.players.iter().find(|p| p.name == "abmitdirpic").expect("abmitdirpic not found");
    let first_left_minute = ((abmit.current_time_played as f64 / 1000.0 / 60.0).ceil()) as usize;
    let post_leave_sum: u32 = abmit.actions.timed[first_left_minute..].iter().sum();
    assert_eq!(post_leave_sum, 0);
    assert_eq!(abmit.apm, 98);
    assert_eq!(abmit.current_time_played, 4371069);
}
