use w3grs::parse_file;
use w3grs::types::ObserverMode;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_132_reforged1() {
    let result = parse_file(&replay_path("132", "reforged1.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6091);
    assert_eq!(result.players.len(), 2);
}

#[test]
fn test_132_reforged2() {
    let result = parse_file(&replay_path("132", "reforged2.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6091);
    assert_eq!(result.players.len(), 2);
}

#[test]
fn test_132_reforged2010() {
    let result = parse_file(&replay_path("132", "reforged2010.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6102);
    assert_eq!(result.players.len(), 6);
    assert_eq!(result.players.iter().find(|p| p.name == "BEARAND#1604").is_some(), true);
}

#[test]
fn test_132_reforged_release() {
    let result = parse_file(&replay_path("132", "reforged_release.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6105);
    assert_eq!(result.players.len(), 2);
    assert!(result.players.iter().any(|p| p.name == "anXieTy#2932"));
    assert!(result.players.iter().any(|p| p.name == "IroNSoul#22724"));
}

#[test]
fn test_132_reforged_hunter2_privatestring() {
    let result = parse_file(&replay_path("132", "reforged_hunter2_privatestring.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6105);
    assert_eq!(result.players.len(), 2);
    assert!(result.players.iter().any(|p| p.name == "pischner#2950"));
    assert!(result.players.iter().any(|p| p.name == "Wartoni#2638"));
}

#[test]
fn test_132_netease() {
    let result = parse_file(&replay_path("132", "netease_132.nwg")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6105);
    assert_eq!(result.players.len(), 2);
    assert!(result.players.iter().any(|p| p.name == "HurricaneBo"));
    assert!(result.players.iter().any(|p| p.name == "SimplyHunteR"));
}

#[test]
fn test_132_reforged_truncated_playernames() {
    let result = parse_file(&replay_path("132", "reforged_truncated_playernames.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.32");
    assert_eq!(result.build_number, 6105);
    assert_eq!(result.players.len(), 2);
    assert!(result.players.iter().any(|p| p.name == "WaN#1734"));
    assert!(result.players.iter().any(|p| p.name.contains("1734") || p.name.contains("228941")));
}

#[test]
fn test_132_random_hero_random_races() {
    let result = parse_file(&replay_path("132", "replay_randomhero_randomraces.w3g")).expect("parse failed");
    assert!(result.settings.random_hero);
    assert!(result.settings.random_races);
}

#[test]
fn test_132_teams_together_settings() {
    let result = parse_file(&replay_path("132", "replay_teamstogether.w3g")).expect("parse failed");
    assert!(result.settings.full_shared_unit_control);
    assert!(result.settings.teams_together);
    assert!(result.settings.fixed_teams);
    assert!(!result.settings.random_hero);
    assert!(!result.settings.random_races);
}

#[test]
fn test_132_full_observers() {
    let result = parse_file(&replay_path("132", "replay_fullobs.w3g")).expect("parse failed");
    assert_eq!(result.settings.observer_mode, ObserverMode::Full);
}

#[test]
fn test_132_referees() {
    let result = parse_file(&replay_path("132", "replay_referee.w3g")).expect("parse failed");
    assert_eq!(result.settings.observer_mode, ObserverMode::Referees);
}

#[test]
fn test_132_obs_on_defeat() {
    let result = parse_file(&replay_path("132", "replay_obs_on_defeat.w3g")).expect("parse failed");
    assert_eq!(result.settings.observer_mode, ObserverMode::OnDefeat);
}

#[test]
fn test_132_hotkeys() {
    let result = parse_file(&replay_path("132", "reforged1.w3g")).expect("parse failed");
    // players[0]: groupHotkeys[1] = {assigned:1, used:29}, groupHotkeys[2] = {assigned:1, used:60}
    // players[1]: groupHotkeys[1] = {assigned:21, used:106}, groupHotkeys[2] = {assigned:4, used:64}
    let p0 = &result.players[0];
    let p1 = &result.players[1];

    let hk1_p0 = p0.group_hotkeys.get(&1).expect("hotkey 1 not found");
    assert_eq!(hk1_p0.assigned, 1);
    assert_eq!(hk1_p0.used, 29);

    let hk2_p0 = p0.group_hotkeys.get(&2).expect("hotkey 2 not found");
    assert_eq!(hk2_p0.assigned, 1);
    assert_eq!(hk2_p0.used, 60);

    let hk1_p1 = p1.group_hotkeys.get(&1).expect("hotkey 1 not found");
    assert_eq!(hk1_p1.assigned, 21);
    assert_eq!(hk1_p1.used, 106);

    let hk2_p1 = p1.group_hotkeys.get(&2).expect("hotkey 2 not found");
    assert_eq!(hk2_p1.assigned, 4);
    assert_eq!(hk2_p1.used, 64);
}

#[test]
fn test_132_kotg_level_6() {
    let result = parse_file(&replay_path("132", "706266088.w3g")).expect("parse failed");
    let kotg_player = result.players.iter()
        .find(|p| p.heroes.iter().any(|h| h.id == "Ekee"))
        .expect("player with KotG not found");
    let kotg = kotg_player.heroes.iter().find(|h| h.id == "Ekee").unwrap();
    assert_eq!(kotg.level, 6);
}

#[test]
fn test_132_winner_1640262494() {
    let result = parse_file(&replay_path("132", "1640262494.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "Happie");
}

#[test]
fn test_132_winner_1448202825() {
    let result = parse_file(&replay_path("132", "1448202825.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "ThundeR#31281");
}

#[test]
fn test_132_winner_wan_vs_trunks() {
    let result = parse_file(&replay_path("132", "wan_vs_trunks.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "WaN#1734");
}

#[test]
fn test_132_winner_benjiii() {
    let result = parse_file(&replay_path("132", "benjiii_vs_Scars_Concealed_Hill.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "benjiii#1588");
}

#[test]
fn test_132_winner_esl_cup_changer() {
    let result = parse_file(&replay_path("132", "esl_cup_vs_changer_1.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "TapioN#2351");
}

#[test]
fn test_132_winner_buildingwin_anxiety() {
    let result = parse_file(&replay_path("132", "buildingwin_anxietyperspective.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "anXieTy#2932");
}

#[test]
fn test_132_winner_buildingwin_helpstone() {
    let result = parse_file(&replay_path("132", "buildingwin_helpstoneperspective.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "anXieTy#2932");
}
