use std::collections::{HashMap, HashSet};
use std::time::Instant;
use sha2::{Sha256, Digest};

use crate::convert::{game_version, map_filename};
use crate::parser::action::{Action, format_object_id, ObjectId};
use crate::parser::game_data::{parse_game_data, GameDataBlock};
use crate::parser::metadata::{parse as parse_metadata, SlotRecord};
use crate::parser::raw::{parse as parse_raw, decompress_blocks};
use crate::player::Player;
use crate::types::*;
use crate::ReplayHandler;

pub fn parse_file(path: &str) -> Option<ParserOutput> {
    let data = std::fs::read(path).ok()?;
    parse_bytes(&data)
}

pub fn parse_file_with_handler<H: ReplayHandler>(path: &str, handler: &mut H) -> Option<ParserOutput> {
    let data = std::fs::read(path).ok()?;
    parse_bytes_with_handler(&data, handler)
}

struct NoopHandler;
impl ReplayHandler for NoopHandler {}

pub fn parse_bytes(input: &[u8]) -> Option<ParserOutput> {
    parse_bytes_with_handler(input, &mut NoopHandler)
}

pub fn parse_bytes_with_handler<H: ReplayHandler>(input: &[u8], handler: &mut H) -> Option<ParserOutput> {
    let start = Instant::now();

    let (header, subheader, blocks) = parse_raw(input)?;
    let uncompressed = decompress_blocks(&blocks);
    let metadata = parse_metadata(&uncompressed)?;

    let game_data_blocks = parse_game_data(&metadata.game_data, metadata.is_post_202);

    // Build players from slots
    let mut temp_players: HashMap<u8, crate::parser::metadata::PlayerRecord> = HashMap::new();
    for pr in &metadata.player_records {
        temp_players.insert(pr.player_id, pr.clone());
    }

    // Apply reforged player names
    for extra in &metadata.reforged_player_metadata {
        if let Some(rec) = temp_players.get_mut(&(extra.player_id as u8)) {
            rec.player_name = extra.name.clone();
        }
    }

    let mut players: HashMap<u8, Player> = HashMap::new();
    let mut teams: HashMap<u8, Vec<u8>> = HashMap::new();
    let mut slot_to_player_id: HashMap<usize, u8> = HashMap::new();

    for (slot_index, slot) in metadata.slot_records.iter().enumerate() {
        if slot.slot_status > 1 {
            slot_to_player_id.insert(slot_index, slot.player_id);
            teams.entry(slot.team_id).or_default().push(slot.player_id);

            let name = temp_players.get(&slot.player_id)
                .map(|r| r.player_name.clone())
                .unwrap_or_else(|| "Computer".to_string());

            let race = race_flag_to_string(slot.race_flag);
            players.insert(slot.player_id, Player::new(
                slot.player_id, name, slot.team_id, slot.color, race,
            ));
        }
    }

    let known_player_ids: HashSet<u8> = players.keys().copied().collect();

    // Emit on_basic_replay_information before game data processing
    let basic_players: Vec<BasicPlayerInfo> = metadata.slot_records.iter()
        .filter(|s| s.slot_status > 1)
        .map(|s| BasicPlayerInfo {
            player_id: s.player_id,
            name: temp_players.get(&s.player_id)
                .map(|r| r.player_name.clone())
                .unwrap_or_else(|| "Computer".to_string()),
            team_id: s.team_id,
            color: s.color,
            race: race_flag_to_string(s.race_flag),
        })
        .collect();
    let basic_info = BasicReplayInfo {
        build_number: subheader.build_no,
        version: game_version(subheader.version),
        game_name: metadata.game_name.clone(),
        random_seed: metadata.random_seed,
        start_spots: metadata.start_spot_count,
        map: MapInfo {
            path: metadata.map.map_name.clone(),
            file: map_filename(&metadata.map.map_name),
            checksum: metadata.map.map_checksum.clone(),
            checksum_sha1: metadata.map.map_checksum_sha1.clone(),
        },
        players: basic_players,
        expansion: subheader.game_identifier == "PX3W",
    };
    handler.on_basic_replay_information(&basic_info);

    let player_action_track_interval: u32 = 60000;

    let mut total_time_tracker: u32 = 0;
    let mut time_segment_tracker: u32 = 0;
    let mut ms_elapsed: u32 = 0;
    let mut chat_log: Vec<ChatMessage> = Vec::new();
    let mut leave_events: Vec<crate::parser::game_data::LeaveGameBlock> = Vec::new();

    for block in &game_data_blocks {
        handler.on_gamedatablock(block);
        match block {
            GameDataBlock::Timeslot(ts) => {
                total_time_tracker += ts.time_increment as u32;
                time_segment_tracker += ts.time_increment as u32;
                ms_elapsed += ts.time_increment as u32;

                if time_segment_tracker > player_action_track_interval {
                    for p in players.values_mut() {
                        p.new_action_tracking_segment(player_action_track_interval);
                    }
                    time_segment_tracker = 0;
                }

                for cmd in &ts.command_blocks {
                    if !known_player_ids.contains(&cmd.player_id) {
                        eprintln!(
                            "detected unknown playerId in CommandBlock: {} - time elapsed: {}",
                            cmd.player_id, total_time_tracker
                        );
                        continue;
                    }
                    let player = match players.get_mut(&cmd.player_id) {
                        Some(p) => p,
                        None => continue,
                    };
                    player.current_time_played = total_time_tracker;
                    player._last_action_was_deselect = false;

                    let actions_clone = cmd.actions.clone();
                    for action in &actions_clone {
                        process_action(
                            action,
                            player,
                            total_time_tracker,
                            &slot_to_player_id,
                            &mut HashMap::new(), // player_names will be looked up separately
                        );
                    }
                }
            }
            GameDataBlock::ChatMessage(chat) => {
                if let Some(p) = players.get(&chat.player_id) {
                    let mode = match chat.mode {
                        0x00 => ChatMode::All,
                        0x01 => ChatMode::Team,
                        0x02 => ChatMode::Observers,
                        _ => ChatMode::Private,
                    };
                    chat_log.push(ChatMessage {
                        player_name: p.name.clone(),
                        player_id: chat.player_id,
                        message: chat.message.clone(),
                        mode,
                        time_ms: total_time_tracker,
                    });
                }
            }
            GameDataBlock::LeaveGame(lg) => {
                leave_events.push(lg.clone());
            }
        }
    }

    // Cleanup players - build player_names for resource transfers
    let player_names: HashMap<u8, String> = players.iter()
        .map(|(id, p)| (*id, p.name.clone()))
        .collect();

    // Re-process resource transfers with player names (we deferred them above)
    // Actually let's do it properly - we need to pass player names during action processing
    // Since we can't borrow players twice, re-process transfers in a second pass
    // Instead, let's rebuild transfers with proper names now:
    for p in players.values_mut() {
        for tr in p.resource_transfers.iter_mut() {
            if let Some(name) = player_names.get(&tr.player_id) {
                tr.player_name = name.clone();
            }
        }
    }

    // Determine version for observer team
    let version_num = subheader.version;

    // Separate observers from players (sorted by ID for determinism)
    let mut observers: Vec<String> = Vec::new();
    let mut final_players: HashMap<u8, Player> = HashMap::new();

    let mut sorted_player_ids: Vec<u8> = players.keys().cloned().collect();
    sorted_player_ids.sort();

    for id in sorted_player_ids {
        let mut player = players.remove(&id).unwrap();
        let is_obs = is_observer(&player, version_num);
        player.cleanup(player_action_track_interval);

        if is_obs {
            observers.push(player.name.clone());
        } else {
            final_players.insert(id, player);
        }
    }

    // Remove observer teams
    if version_num >= 29 {
        teams.remove(&24);
    } else {
        teams.remove(&12);
    }

    // Determine matchup and gametype
    let (gametype, matchup) = determine_matchup(&final_players, version_num);

    // Determine winning team (1on1)
    let winning_team_id = determine_winning_team(&gametype, &leave_events, &final_players, version_num);

    // Generate ID
    let id = generate_id(metadata.random_seed, &final_players, &metadata.game_name);

    // Sort players
    let mut sorted_players: Vec<Player> = final_players.into_values().collect();
    sorted_players.sort_by(|a, b| {
        if a.teamid != b.teamid {
            a.teamid.cmp(&b.teamid)
        } else {
            a.id.cmp(&b.id)
        }
    });

    let settings = build_settings(&metadata.map, subheader.version);

    let elapsed_ms = start.elapsed().as_millis() as u64;

    Some(ParserOutput {
        id,
        gamename: metadata.game_name,
        randomseed: metadata.random_seed,
        start_spots: metadata.start_spot_count,
        observers,
        players: sorted_players,
        matchup,
        creator: metadata.map.creator.clone(),
        game_type: gametype,
        chat: chat_log,
        apm: ApmConfig { tracking_interval: player_action_track_interval },
        map: MapInfo {
            path: metadata.map.map_name.clone(),
            file: map_filename(&metadata.map.map_name),
            checksum: metadata.map.map_checksum.clone(),
            checksum_sha1: metadata.map.map_checksum_sha1.clone(),
        },
        build_number: subheader.build_no,
        version: game_version(subheader.version),
        duration: subheader.replay_length_ms,
        expansion: subheader.game_identifier == "PX3W",
        parse_time: elapsed_ms,
        winning_team_id,
        settings,
    })
}

fn process_action(
    action: &Action,
    player: &mut Player,
    total_time: u32,
    slot_to_player_id: &HashMap<usize, u8>,
    _player_names: &mut HashMap<u8, String>,
) {
    match action {
        Action::UnitAbilityNoParams { order_id, .. } => {
            let fmt = format_object_id(order_id);
            // Check for retraining (tert or tret)
            if let ObjectId::StringEncoded(s) = &fmt {
                if s == "tert" || s == "tret" {
                    player.handle_retraining(total_time);
                }
            }
            player.handle_0x10(&fmt, total_time);
        }
        Action::UnitAbilityTargetPos { order_id, .. } => {
            let fmt = format_object_id(order_id);
            player.handle_0x11(&fmt, total_time);
        }
        Action::UnitAbilityTargetObj { order_id, .. } => {
            let fmt = format_object_id(order_id);
            player.handle_0x12(&fmt, total_time);
        }
        Action::GiveItemToUnit { .. } => {
            player.handle_0x13();
        }
        Action::UnitAbilityTwoTargets { order_id1, .. } => {
            let fmt = format_object_id(order_id1);
            player.handle_0x14(&fmt);
        }
        Action::UnitAbilityTwoTargetsItem { order_id1, .. } => {
            let fmt = format_object_id(order_id1);
            player.handle_0x14(&fmt);
        }
        Action::ChangeSelection { select_mode, .. } => {
            if *select_mode == 0x02 {
                player._last_action_was_deselect = true;
                player.handle_0x16(*select_mode, true);
            } else {
                if !player._last_action_was_deselect {
                    player.handle_0x16(*select_mode, true);
                }
                player._last_action_was_deselect = false;
            }
        }
        Action::AssignGroupHotkey { .. } |
        Action::SelectGroupHotkey { .. } |
        Action::SelectGroundItem { .. } |
        Action::CancelHeroRevival { .. } |
        Action::RemoveUnitFromQueue { .. } |
        Action::EscPressed |
        Action::ChooseHeroSkillSubmenu |
        Action::EnterBuildingSubmenu => {
            player.handle_other(action);
        }
        Action::TransferResources { slot, gold, lumber } => {
            if let Some(&pid) = slot_to_player_id.get(&(*slot as usize)) {
                // We'll fix the name later in a second pass
                player.handle_0x51(*slot, pid, String::new(), *gold, *lumber);
            }
        }
        _ => {}
    }
}

fn is_observer(player: &Player, version: u32) -> bool {
    if version >= 29 {
        player.teamid == 24
    } else {
        player.teamid == 12
    }
}

fn determine_matchup(players: &HashMap<u8, Player>, version: u32) -> (String, String) {
    let mut team_races: HashMap<u8, Vec<String>> = HashMap::new();

    for p in players.values() {
        let race = if p.race_detected.is_empty() { p.race.clone() } else { p.race_detected.clone() };
        team_races.entry(p.teamid).or_default().push(race);
    }

    let mut sizes: Vec<usize> = team_races.values().map(|v| v.len()).collect();
    sizes.sort_unstable();
    let gametype = sizes.iter().map(|s| s.to_string()).collect::<Vec<_>>().join("on");

    let mut race_groups: Vec<String> = team_races.values()
        .map(|races| {
            let mut sorted = races.clone();
            sorted.sort();
            sorted.join("")
        })
        .collect();
    race_groups.sort();
    let matchup = race_groups.join("v");

    (gametype, matchup)
}

fn determine_winning_team(
    gametype: &str,
    leave_events: &[crate::parser::game_data::LeaveGameBlock],
    players: &HashMap<u8, Player>,
    _version: u32,
) -> i32 {
    if gametype != "1on1" {
        return -1;
    }

    // Filter to non-observer leave events (observers are not in the players map)
    let non_obs_leaves: Vec<_> = leave_events.iter()
        .filter(|e| players.contains_key(&e.player_id))
        .collect();

    // Tier 1: player left with victory result code
    if let Some(event) = non_obs_leaves.iter().find(|e| e.result == "09000000") {
        return players[&event.player_id].teamid as i32;
    }

    // Tier 2: WC3 ended the game for the winner (game-over reason)
    if let Some(event) = non_obs_leaves.iter().find(|e| e.reason == "0c000000") {
        return players[&event.player_id].teamid as i32;
    }

    // Tier 3: first non-observer to leave is the loser; the other player wins
    if let Some(first) = non_obs_leaves.first() {
        let loser_team = players[&first.player_id].teamid;
        if let Some(winner) = players.values().find(|p| p.teamid != loser_team) {
            return winner.teamid as i32;
        }
    }

    -1
}

fn generate_id(random_seed: u32, players: &HashMap<u8, Player>, game_name: &str) -> String {
    let mut sorted_ids: Vec<u8> = players.keys().copied().collect();
    sorted_ids.sort_unstable();
    let names: String = sorted_ids.iter()
        .map(|id| players[id].name.as_str())
        .collect();
    let id_base = format!("{}{}{}", random_seed, names, game_name);
    let mut hasher = Sha256::new();
    hasher.update(id_base.as_bytes());
    hex::encode(hasher.finalize())
}

fn race_flag_to_string(flag: u8) -> String {
    match flag {
        0x01 | 0x41 => "H",
        0x02 | 0x42 => "O",
        0x04 | 0x44 => "N",
        0x08 | 0x48 => "U",
        0x20 | 0x60 => "R",
        _ => "R",
    }.to_string()
}

fn build_settings(map: &crate::parser::metadata::MapMetadata, version: u32) -> Settings {
    let observer_mode = get_observer_mode(map.referees, map.observer_mode);
    Settings {
        observer_mode,
        referees: map.referees,
        fixed_teams: map.fixed_teams,
        full_shared_unit_control: map.full_shared_unit_control,
        always_visible: map.always_visible,
        hide_terrain: map.hide_terrain,
        map_explored: map.map_explored,
        teams_together: map.teams_together,
        random_hero: map.random_hero,
        random_races: map.random_races,
        speed: map.speed,
    }
}

fn get_observer_mode(referee_flag: bool, observer_mode: u8) -> ObserverMode {
    if (observer_mode == 3 || observer_mode == 0) && referee_flag {
        ObserverMode::Referees
    } else if observer_mode == 2 {
        ObserverMode::OnDefeat
    } else if observer_mode == 3 {
        ObserverMode::Full
    } else {
        ObserverMode::None
    }
}
