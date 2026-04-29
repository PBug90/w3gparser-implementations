use super::action::{parse_actions, Action};

#[derive(Debug, Clone)]
pub struct LeaveGameBlock {
    pub player_id: u8,
    pub reason: String,
    pub result: String,
}

#[derive(Debug, Clone)]
pub struct CommandBlock {
    pub player_id: u8,
    pub actions: Vec<Action>,
}

#[derive(Debug, Clone)]
pub struct TimeslotBlock {
    pub time_increment: u16,
    pub command_blocks: Vec<CommandBlock>,
}

#[derive(Debug, Clone)]
pub struct ChatMessageBlock {
    pub player_id: u8,
    pub mode: u32,
    pub message: String,
}

#[derive(Debug, Clone)]
pub enum GameDataBlock {
    LeaveGame(LeaveGameBlock),
    Timeslot(TimeslotBlock),
    ChatMessage(ChatMessageBlock),
}

pub fn parse_game_data(data: &[u8], is_post_202: bool) -> Vec<GameDataBlock> {
    let mut blocks = Vec::new();
    let mut pos = 0usize;

    while pos < data.len() {
        let id = match data.get(pos) {
            Some(&v) => { pos += 1; v }
            None => break,
        };

        match id {
            0x17 => {
                if let Some(b) = parse_leave_game(data, &mut pos) {
                    blocks.push(GameDataBlock::LeaveGame(b));
                }
            }
            0x1a | 0x1b | 0x1c => {
                // skip 4 bytes
                pos = (pos + 4).min(data.len());
            }
            0x1e | 0x1f => {
                if let Some(b) = parse_timeslot(data, &mut pos, is_post_202) {
                    blocks.push(GameDataBlock::Timeslot(b));
                }
            }
            0x20 => {
                if let Some(b) = parse_chat_message(data, &mut pos) {
                    blocks.push(GameDataBlock::ChatMessage(b));
                }
            }
            0x22 => {
                let len = data.get(pos).copied().unwrap_or(0) as usize;
                pos += 1;
                pos = (pos + len).min(data.len());
            }
            0x23 => {
                pos = (pos + 10).min(data.len());
            }
            0x2f => {
                pos = (pos + 8).min(data.len());
            }
            _ => {
                // unknown block id - skip and continue (mirrors w3gjs behavior)
                continue;
            }
        }
    }
    blocks
}

fn parse_leave_game(data: &[u8], pos: &mut usize) -> Option<LeaveGameBlock> {
    let reason = read_hex(data, pos, 4);
    let player_id = ru8(data, pos)?;
    let result = read_hex(data, pos, 4);
    *pos += 4; // skip 4
    Some(LeaveGameBlock { player_id, reason, result })
}

fn parse_timeslot(data: &[u8], pos: &mut usize, is_post_202: bool) -> Option<TimeslotBlock> {
    let byte_count = ru16(data, pos)? as usize;
    let time_increment = ru16(data, pos)?;
    let action_block_last_offset = *pos + byte_count.saturating_sub(2);

    let mut command_blocks = Vec::new();

    while *pos < action_block_last_offset && *pos < data.len() {
        let player_id = ru8(data, pos)?;
        let action_block_length = ru16(data, pos)? as usize;
        let end = (*pos + action_block_length).min(data.len());
        let action_data = &data[*pos..end];
        let actions = parse_actions(action_data, is_post_202);
        *pos = (*pos + action_block_length).min(data.len());
        command_blocks.push(CommandBlock { player_id, actions });
    }

    Some(TimeslotBlock { time_increment, command_blocks })
}

fn parse_chat_message(data: &[u8], pos: &mut usize) -> Option<ChatMessageBlock> {
    let player_id = ru8(data, pos)?;
    *pos += 2; // byteCount
    let flags = ru8(data, pos)?;
    let mut mode = 0u32;
    if flags == 0x20 {
        mode = ru32(data, pos)?;
    }
    let message = read_zts(data, pos);
    Some(ChatMessageBlock { player_id, mode, message })
}

// helpers
fn ru8(data: &[u8], pos: &mut usize) -> Option<u8> {
    let v = *data.get(*pos)?;
    *pos += 1;
    Some(v)
}

fn ru16(data: &[u8], pos: &mut usize) -> Option<u16> {
    if *pos + 2 > data.len() { return None; }
    let v = u16::from_le_bytes([data[*pos], data[*pos+1]]);
    *pos += 2;
    Some(v)
}

fn ru32(data: &[u8], pos: &mut usize) -> Option<u32> {
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

fn read_zts(data: &[u8], pos: &mut usize) -> String {
    let start = *pos;
    while *pos < data.len() && data[*pos] != 0 { *pos += 1; }
    let s = String::from_utf8_lossy(&data[start..*pos]).into_owned();
    if *pos < data.len() { *pos += 1; }
    s
}
