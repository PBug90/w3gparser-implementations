use super::buffer::BufParser;
use super::protobuf::decode_player_metadata;

#[derive(Debug, Clone)]
pub struct PlayerRecord {
    pub player_id: u8,
    pub player_name: String,
}

#[derive(Debug, Clone)]
pub struct SlotRecord {
    pub player_id: u8,
    pub download_progress: u8,
    pub slot_status: u8,
    pub computer_flag: u8,
    pub team_id: u8,
    pub color: u8,
    pub race_flag: u8,
    pub ai_strength: u8,
    pub handicap_flag: u8,
}

#[derive(Debug, Clone)]
pub struct ReforgedPlayerMetadata {
    pub player_id: u32,
    pub name: String,
    pub clan: String,
}

#[derive(Debug, Clone)]
pub struct MapMetadata {
    pub speed: u8,
    pub hide_terrain: bool,
    pub map_explored: bool,
    pub always_visible: bool,
    pub default: bool,
    pub observer_mode: u8,
    pub teams_together: bool,
    pub fixed_teams: bool,
    pub full_shared_unit_control: bool,
    pub random_hero: bool,
    pub random_races: bool,
    pub referees: bool,
    pub map_checksum: String,
    pub map_checksum_sha1: String,
    pub map_name: String,
    pub creator: String,
}

#[derive(Debug)]
pub struct Metadata {
    pub game_data: Vec<u8>,
    pub map: MapMetadata,
    pub player_count: u32,
    pub game_type: String,
    pub locale_hash: String,
    pub player_records: Vec<PlayerRecord>,
    pub slot_records: Vec<SlotRecord>,
    pub reforged_player_metadata: Vec<ReforgedPlayerMetadata>,
    pub random_seed: u32,
    pub select_mode: String,
    pub game_name: String,
    pub start_spot_count: u8,
    pub is_post_202: bool,
}

/// Read a null-terminated sequence of raw bytes from parser (for binary data).
fn read_zero_term_raw(data: &[u8], pos: &mut usize) -> Vec<u8> {
    let start = *pos;
    while *pos < data.len() && data[*pos] != 0 {
        *pos += 1;
    }
    let result = data[start..*pos].to_vec();
    if *pos < data.len() {
        *pos += 1; // consume null
    }
    result
}

pub fn parse(data: &[u8]) -> Option<Metadata> {
    let mut pos = 0usize;

    // skip 5 bytes
    pos += 5;

    // Parse host record
    let mut player_records = Vec::new();
    player_records.push(parse_host_record(data, &mut pos)?);

    // game name (UTF-8 null-terminated)
    let game_name = read_zts_utf8(data, &mut pos);

    // private string (discard)
    read_zts_utf8(data, &mut pos);

    // encoded string - read as raw bytes (then hex-encoded in JS, but we keep raw)
    let encoded_raw = read_zero_term_raw(data, &mut pos);

    // Decode game meta string
    let map_meta_decoded = decode_game_meta_string(&encoded_raw);
    let map_metadata = parse_encoded_map_meta_string(&map_meta_decoded)?;

    // player count
    let player_count = read_u32_le(data, &mut pos)?;

    // game type (4 bytes as hex)
    let game_type = read_hex(data, &mut pos, 4);

    // locale hash (4 bytes as hex)
    let locale_hash = read_hex(data, &mut pos, 4);

    // Additional player list
    let additional = parse_player_list(data, &mut pos);

    // JS: playerListFinal = playerRecords.concat(playerRecords, parsePlayerList())
    let mut player_records_final = player_records.clone();
    player_records_final.extend(player_records.clone());
    player_records_final.extend(additional);

    let mut reforged_player_metadata: Vec<ReforgedPlayerMetadata> = Vec::new();
    let mut is_post_202 = false;

    // Reforged metadata blocks (0x38 or 0x39)
    while pos < data.len() && (data[pos] == 0x38 || data[pos] == 0x39) {
        let record_type = data[pos]; pos += 1;
        if record_type == 0x38 { is_post_202 = true; }
        let subtype = read_u8(data, &mut pos)?;
        let following_bytes = read_u32_le(data, &mut pos)? as usize;
        let blob_end = pos + following_bytes;
        let blob = if blob_end <= data.len() { &data[pos..blob_end] } else { break };

        if subtype == 0x03 {
            let players = decode_player_metadata(blob);
            for (pid, battle_tag, clan) in players {
                reforged_player_metadata.push(ReforgedPlayerMetadata {
                    player_id: pid,
                    name: battle_tag,
                    clan,
                });
            }
        }
        pos += following_bytes;
    }

    // Expect 0x19 = 25
    let _check = read_u8(data, &mut pos);
    let _remaining_bytes = read_u16_le(data, &mut pos)?;
    let slot_record_count = read_u8(data, &mut pos)? as usize;

    let slot_records = parse_slot_records(data, &mut pos, slot_record_count);
    let random_seed = read_u32_le(data, &mut pos)?;
    let select_mode = read_hex(data, &mut pos, 1);
    let start_spot_count = read_u8(data, &mut pos)?;

    let game_data = data[pos..].to_vec();

    Some(Metadata {
        game_data,
        map: map_metadata,
        player_count,
        game_type,
        locale_hash,
        player_records: player_records_final,
        slot_records,
        reforged_player_metadata,
        random_seed,
        select_mode,
        game_name,
        start_spot_count,
        is_post_202,
    })
}

fn parse_host_record(data: &[u8], pos: &mut usize) -> Option<PlayerRecord> {
    let player_id = read_u8(data, pos)?;
    let player_name = read_zts_utf8(data, pos);
    let add_data = read_u8(data, pos)? as usize;
    *pos += add_data;
    Some(PlayerRecord { player_id, player_name })
}

fn parse_player_list(data: &[u8], pos: &mut usize) -> Vec<PlayerRecord> {
    let mut list = Vec::new();
    while *pos < data.len() && data[*pos] == 22 {
        *pos += 1; // consume the 22
        if let Some(record) = parse_host_record(data, pos) {
            list.push(record);
        }
        *pos += 4;
    }
    list
}

fn parse_slot_records(data: &[u8], pos: &mut usize, count: usize) -> Vec<SlotRecord> {
    let mut slots = Vec::new();
    for _ in 0..count {
        slots.push(SlotRecord {
            player_id: read_u8(data, pos).unwrap_or(0),
            download_progress: read_u8(data, pos).unwrap_or(0),
            slot_status: read_u8(data, pos).unwrap_or(0),
            computer_flag: read_u8(data, pos).unwrap_or(0),
            team_id: read_u8(data, pos).unwrap_or(0),
            color: read_u8(data, pos).unwrap_or(0),
            race_flag: read_u8(data, pos).unwrap_or(0),
            ai_strength: read_u8(data, pos).unwrap_or(0),
            handicap_flag: read_u8(data, pos).unwrap_or(0),
        });
    }
    slots
}

fn decode_game_meta_string(bytes: &[u8]) -> Vec<u8> {
    let mut decoded = Vec::new();
    let mut mask: u8 = 0;
    for (i, &b) in bytes.iter().enumerate() {
        if i % 8 == 0 {
            mask = b;
        } else {
            let bit_pos = i % 8;
            if (mask & (1u8 << bit_pos)) == 0 {
                decoded.push(b.wrapping_sub(1));
            } else {
                decoded.push(b);
            }
        }
    }
    decoded
}

fn parse_encoded_map_meta_string(data: &[u8]) -> Option<MapMetadata> {
    let mut pos = 0usize;
    let speed = read_u8(data, &mut pos)?;
    let second_byte = read_u8(data, &mut pos)?;
    let third_byte = read_u8(data, &mut pos)?;
    let fourth_byte = read_u8(data, &mut pos)?;
    pos += 5; // skip 5
    let checksum = read_hex(data, &mut pos, 4);
    let map_name = read_zts_utf8(data, &mut pos);
    let creator = read_zts_utf8(data, &mut pos);
    pos += 1; // skip 1
    let checksum_sha1 = read_hex(data, &mut pos, 20);

    Some(MapMetadata {
        speed,
        hide_terrain:              (second_byte & 0b00000001) != 0,
        map_explored:              (second_byte & 0b00000010) != 0,
        always_visible:            (second_byte & 0b00000100) != 0,
        default:                   (second_byte & 0b00001000) != 0,
        observer_mode:             (second_byte & 0b00110000) >> 4,
        teams_together:            (second_byte & 0b01000000) != 0,
        fixed_teams:               (third_byte  & 0b00000110) != 0,
        full_shared_unit_control:  (fourth_byte & 0b00000001) != 0,
        random_hero:               (fourth_byte & 0b00000010) != 0,
        random_races:              (fourth_byte & 0b00000100) != 0,
        referees:                  (fourth_byte & 0b01000000) != 0,
        map_checksum: checksum,
        map_checksum_sha1: checksum_sha1,
        map_name,
        creator,
    })
}

// --- helpers ---

fn read_u8(data: &[u8], pos: &mut usize) -> Option<u8> {
    let v = *data.get(*pos)?;
    *pos += 1;
    Some(v)
}

fn read_u16_le(data: &[u8], pos: &mut usize) -> Option<u16> {
    if *pos + 2 > data.len() { return None; }
    let v = u16::from_le_bytes([data[*pos], data[*pos+1]]);
    *pos += 2;
    Some(v)
}

fn read_u32_le(data: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > data.len() { return None; }
    let v = u32::from_le_bytes([data[*pos], data[*pos+1], data[*pos+2], data[*pos+3]]);
    *pos += 4;
    Some(v)
}

fn read_hex(data: &[u8], pos: &mut usize, n: usize) -> String {
    if *pos + n > data.len() { return String::new(); }
    let s = hex::encode(&data[*pos..*pos+n]);
    *pos += n;
    s
}

fn read_zts_utf8(data: &[u8], pos: &mut usize) -> String {
    let start = *pos;
    while *pos < data.len() && data[*pos] != 0 {
        *pos += 1;
    }
    let s = String::from_utf8_lossy(&data[start..*pos]).into_owned();
    if *pos < data.len() { *pos += 1; }
    s
}
