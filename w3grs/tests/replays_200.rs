use w3grs::parse_file;

fn replay_path(version: &str, name: &str) -> String {
    format!("{}/tests/replays/{}/{}", env!("CARGO_MANIFEST_DIR"), version, name)
}

#[test]
fn test_200_haunted_gold_mine() {
    let result = parse_file(&replay_path("200", "goldmine test.w3g")).expect("parse failed");
    let p = &result.players[0];
    assert_eq!(p.buildings.summary.get("ugol").copied().unwrap_or(0), 1);
    assert_eq!(p.buildings.order.len(), 1);
    assert_eq!(p.buildings.order[0].id, "ugol");
    assert_eq!(p.buildings.order[0].ms, 28435);
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "anXieTy#2932");
}

#[test]
fn test_200_version() {
    let result = parse_file(&replay_path("200", "goldmine test.w3g")).expect("parse failed");
    assert_eq!(result.version, "2.00");
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "anXieTy#2932");
}

#[test]
fn test_200_custom_map_ui_components() {
    let result = parse_file(&replay_path("200", "TempReplay.w3g")).expect("parse failed");
    assert_eq!(result.version, "2.00");
    assert_eq!(result.winning_team_id, -1);
}

#[test]
fn test_200_retraining() {
    let result = parse_file(&replay_path("200", "retrainingissues.w3g")).expect("parse failed");
    assert_eq!(result.version, "2.00");

    // Find player with Paladin (Hpal) hero
    let player = result.players.iter()
        .find(|p| p.heroes.iter().any(|h| h.id == "Hamg"))
        .expect("player with Archmage hero not found");

    let hamg = player.heroes.iter().find(|h| h.id == "Hamg").unwrap();
    assert_eq!(hamg.level, 6);

    let ao = &hamg.ability_order;
    // Check the ability order matches the test expectation
    use w3grs::player::AbilityOrderEntry;
    // Find the retraining entry
    assert!(ao.iter().any(|e| matches!(e, AbilityOrderEntry::Retraining { .. })));
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "Cechi#21571");
}

#[test]
fn test_200_202_melee_chat() {
    let result = parse_file(&replay_path("200", "2.0.2-Melee.w3g")).expect("parse failed");
    // chatlog[0].playerId == 1, message == "don't hurt me"
    // chatlog[1].playerId == 2, message == "no more"
    let chat = &result.chat;
    assert!(chat.len() >= 2);
    let msg0 = chat.iter().find(|c| c.message == "don't hurt me").expect("first chat not found");
    assert_eq!(msg0.player_id, 1);
    let msg1 = chat.iter().find(|c| c.message == "no more").expect("second chat not found");
    assert_eq!(msg1.player_id, 2);
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "BogdanW4");
}

#[test]
fn test_200_202_flo_tv_saved_by_wc3() {
    let result = parse_file(&replay_path("200", "2.0.2-FloTVSavedByWc3.w3g")).expect("parse failed");
    assert!(result.players.len() >= 1);
    assert_eq!(result.winning_team_id, 0);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "DeMusliM#21529");
}

#[test]
fn test_200_winner_leon_tidehunters() {
    let result = parse_file(&replay_path("200", "148249993_Edo_Leon_Tidehunters 12.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "LeonXIV#21701");
}

#[test]
fn test_200_winner_leon_springtime() {
    let result = parse_file(&replay_path("200", "791786117_Edo_Leon_Springtime 13.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "LeonXIV#21701");
}

#[test]
fn test_200_winner_changer_hammerfall() {
    let result = parse_file(&replay_path("200", "455872485_PhoeNix_Changer_Hammerfall.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "Changer#21493");
}

#[test]
fn test_200_winner_changer_concealed_hill() {
    let result = parse_file(&replay_path("200", "2661392198_PhoeNix_Changer_Concealed Hill.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "Changer#21493");
}

#[test]
fn test_200_winner_dorrani_springtime() {
    let result = parse_file(&replay_path("200", "3320738873_Changer_PhoeNix_Springtime 13.w3g")).expect("parse failed");
    assert_eq!(result.winning_team_id, 1);
    let winner = result.players.iter().find(|p| p.teamid == result.winning_team_id as u8).unwrap();
    assert_eq!(winner.name, "DorraNi#2523");
}
