use w3grs::parse_file;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_131_action_0x7a() {
    let result = parse_file(&replay_path("131", "action0x7a.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.31");
    assert_eq!(result.players.len(), 1);
    assert_eq!(result.winning_team_id, -1);
}

#[test]
fn test_131_tome_of_retraining() {
    let result = parse_file(&replay_path("131", "standard_tomeofretraining_1.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.31");
    assert_eq!(result.build_number, 6072);
    assert_eq!(result.players.len(), 2);

    // Find a player with Hamg hero
    let player_with_hamg = result.players.iter()
        .find(|p| p.heroes.iter().any(|h| h.id == "Hamg"))
        .expect("No player with Hamg hero");

    let hamg = player_with_hamg.heroes.iter().find(|h| h.id == "Hamg").unwrap();
    assert_eq!(hamg.level, 4);
    assert_eq!(hamg.abilities.get("AHab").copied().unwrap_or(0), 2);
    assert_eq!(hamg.abilities.get("AHbz").copied().unwrap_or(0), 2);
    assert_eq!(hamg.retraining_history.len(), 1);
    assert_eq!(hamg.retraining_history[0].time, 1136022);
    assert_eq!(hamg.retraining_history[0].abilities.get("AHab").copied().unwrap_or(0), 2);
    assert_eq!(hamg.retraining_history[0].abilities.get("AHwe").copied().unwrap_or(0), 2);

    // Check ability order
    let ao = &hamg.ability_order;
    assert_eq!(ao.len(), 9);

    use w3grs::player::AbilityOrderEntry;
    match &ao[0] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 124366); assert_eq!(value, "AHwe"); } _ => panic!() }
    match &ao[1] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 234428); assert_eq!(value, "AHab"); } _ => panic!() }
    match &ao[2] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 293007); assert_eq!(value, "AHwe"); } _ => panic!() }
    match &ao[3] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 1060007); assert_eq!(value, "AHab"); } _ => panic!() }
    match &ao[4] { AbilityOrderEntry::Retraining { time } => { assert_eq!(*time, 1136022); } _ => panic!() }
    match &ao[5] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 1140944); assert_eq!(value, "AHbz"); } _ => panic!() }
    match &ao[6] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 1141147); assert_eq!(value, "AHbz"); } _ => panic!() }
    match &ao[7] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 1141460); assert_eq!(value, "AHab"); } _ => panic!() }
    match &ao[8] { AbilityOrderEntry::Ability { time, value } => { assert_eq!(*time, 1141569); assert_eq!(value, "AHab"); } _ => panic!() }
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "[OCG]shocker");
}

#[test]
fn test_131_roc_map_name() {
    let result = parse_file(&replay_path("131", "roc-losttemple-mapname.w3g")).expect("parse failed");
    assert_eq!(result.version, "1.31");
    assert_eq!(result.build_number, 6072);
    assert_eq!(result.map.file, "(4)LostTemple [Unforged 0.5 RoC].w3x");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "syNtec");
}
