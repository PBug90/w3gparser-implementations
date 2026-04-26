/// Action parser - mirrors ActionParser.ts

#[derive(Debug, Clone)]
pub struct Cache {
    pub filename: String,
    pub mission_key: String,
    pub key: String,
}

#[derive(Debug, Clone)]
pub enum Action {
    UnitAbilityNoParams        { ability_flags: u16, order_id: [u8; 4] },
    UnitAbilityTargetPos       { ability_flags: u16, order_id: [u8; 4], target: [f32; 2] },
    UnitAbilityTargetObj       { ability_flags: u16, order_id: [u8; 4], target: [f32; 2], object: [u32; 2] },
    GiveItemToUnit             { ability_flags: u16, order_id: [u8; 4], target: [f32; 2], unit: [u32; 2], item: [u32; 2] },
    UnitAbilityTwoTargets      { ability_flags: u16, order_id1: [u8; 4], target_a: [f32; 2], order_id2: [u8; 4], flags: u32, category: u32, owner: u8, target_b: [f32; 2] },
    UnitAbilityTwoTargetsItem  { ability_flags: u16, order_id1: [u8; 4], target_a: [f32; 2], order_id2: [u8; 4], flags: u32, category: u32, owner: u8, target_b: [f32; 2], object: [u32; 2] },
    ChangeSelection            { select_mode: u8, number_units: u16 },
    AssignGroupHotkey          { group_number: u8, number_units: u16 },
    SelectGroupHotkey          { group_number: u8 },
    SelectSubgroup             { item_id: [u8; 4], object: [u32; 2] },
    PreSubselection,
    SelectUnit                 { object: [u32; 2] },
    SelectGroundItem           { item: [u32; 2] },
    CancelHeroRevival          { hero: [u32; 2] },
    RemoveUnitFromQueue        { slot_number: u8, item_id: [u8; 4] },
    TransferResources          { slot: u8, gold: u32, lumber: u32 },
    EscPressed,
    ChooseHeroSkillSubmenu,
    EnterBuildingSubmenu,
    W3MMDStoreInt              { cache: Cache, value: u32 },
    W3MMDStoreReal             { cache: Cache, value: f32 },
    W3MMDStoreBool             { cache: Cache, value: u8 },
    W3MMDClearInt              { cache: Cache },
    W3MMDClearReal             { cache: Cache },
    W3MMDClearBool             { cache: Cache },
    W3MMDClearUnit             { cache: Cache },
    AllyPing                   { pos: [f32; 2], duration: f32 },
    ArrowKey                   { arrow_key: u8 },
    SetGameSpeed               { game_speed: u8 },
    TrackableHit               { object: [u32; 2] },
    TrackableTrack             { object: [u32; 2] },
    BlzSync                    { identifier: String, value: String },
    CommandFrame               { event_id: u32, val: f32, text: String },
    MouseAction                { event_id: u8, pos: [f32; 2], button: u8 },
    W3Api                      { command_id: u32, data: u32 },
}

pub fn parse_actions(data: &[u8], is_post_202: bool) -> Vec<Action> {
    let mut actions = Vec::new();
    let mut pos = 0usize;

    while pos < data.len() {
        let action_id_raw = match data.get(pos) {
            Some(&v) => { pos += 1; v }
            None => break,
        };

        let action_id = if is_post_202 && action_id_raw > 0x77 {
            action_id_raw + 1
        } else {
            action_id_raw
        };

        match parse_single_action(action_id, data, &mut pos) {
            Some(a) => actions.push(a),
            None => {} // skip or unknown - already advanced pos inside
        }
    }
    actions
}

fn parse_single_action(id: u8, data: &[u8], pos: &mut usize) -> Option<Action> {
    match id {
        0x01 => { advance(data, pos, 1); None }
        0x02 | 0x04 | 0x05 => None,
        0x03 => {
            let gs = ru8(data, pos)?;
            Some(Action::SetGameSpeed { game_speed: gs })
        }
        0x06 => {
            rzts(data, pos); rzts(data, pos); advance(data, pos, 1); None
        }
        0x07 => { advance(data, pos, 4); None }
        0x10 => {
            let ability_flags = ru16(data, pos)?;
            let order_id = rfourcc(data, pos)?;
            advance(data, pos, 8);
            Some(Action::UnitAbilityNoParams { ability_flags, order_id })
        }
        0x11 => {
            let ability_flags = ru16(data, pos)?;
            let order_id = rfourcc(data, pos)?;
            advance(data, pos, 8);
            let target = rvec2(data, pos)?;
            Some(Action::UnitAbilityTargetPos { ability_flags, order_id, target })
        }
        0x12 => {
            let ability_flags = ru16(data, pos)?;
            let order_id = rfourcc(data, pos)?;
            advance(data, pos, 8);
            let target = rvec2(data, pos)?;
            let object = rnettag(data, pos)?;
            Some(Action::UnitAbilityTargetObj { ability_flags, order_id, target, object })
        }
        0x13 => {
            let ability_flags = ru16(data, pos)?;
            let order_id = rfourcc(data, pos)?;
            advance(data, pos, 8);
            let target = rvec2(data, pos)?;
            let unit = rnettag(data, pos)?;
            let item = rnettag(data, pos)?;
            Some(Action::GiveItemToUnit { ability_flags, order_id, target, unit, item })
        }
        0x14 => {
            let ability_flags = ru16(data, pos)?;
            let order_id1 = rfourcc(data, pos)?;
            advance(data, pos, 8);
            let target_a = rvec2(data, pos)?;
            let order_id2 = rfourcc(data, pos)?;
            let flags = ru32(data, pos)?;
            let category = ru32(data, pos)?;
            let owner = ru8(data, pos)?;
            let target_b = rvec2(data, pos)?;
            Some(Action::UnitAbilityTwoTargets { ability_flags, order_id1, target_a, order_id2, flags, category, owner, target_b })
        }
        0x15 => {
            let ability_flags = ru16(data, pos)?;
            let order_id1 = rfourcc(data, pos)?;
            advance(data, pos, 8);
            let target_a = rvec2(data, pos)?;
            let order_id2 = rfourcc(data, pos)?;
            let flags = ru32(data, pos)?;
            let category = ru32(data, pos)?;
            let owner = ru8(data, pos)?;
            let target_b = rvec2(data, pos)?;
            let object = rnettag(data, pos)?;
            Some(Action::UnitAbilityTwoTargetsItem { ability_flags, order_id1, target_a, order_id2, flags, category, owner, target_b, object })
        }
        0x16 => {
            let select_mode = ru8(data, pos)?;
            let number_units = ru16(data, pos)?;
            // skip units: each nettag is 8 bytes
            advance(data, pos, number_units as usize * 8);
            Some(Action::ChangeSelection { select_mode, number_units })
        }
        0x17 => {
            let group_number = ru8(data, pos)?;
            let number_units = ru16(data, pos)?;
            advance(data, pos, number_units as usize * 8);
            Some(Action::AssignGroupHotkey { group_number, number_units })
        }
        0x18 => {
            let group_number = ru8(data, pos)?;
            advance(data, pos, 1);
            Some(Action::SelectGroupHotkey { group_number })
        }
        0x19 => {
            let item_id = rfourcc(data, pos)?;
            let object = rnettag(data, pos)?;
            Some(Action::SelectSubgroup { item_id, object })
        }
        0x1a => Some(Action::PreSubselection),
        0x1b => {
            advance(data, pos, 1);
            let object = rnettag(data, pos)?;
            Some(Action::SelectUnit { object })
        }
        0x1c => {
            advance(data, pos, 1);
            let item = rnettag(data, pos)?;
            Some(Action::SelectGroundItem { item })
        }
        0x1d => {
            let hero = rnettag(data, pos)?;
            Some(Action::CancelHeroRevival { hero })
        }
        0x1e | 0x1f => {
            let slot_number = ru8(data, pos)?;
            let item_id = rfourcc(data, pos)?;
            Some(Action::RemoveUnitFromQueue { slot_number, item_id })
        }
        // cheat actions 0x20-0x4f
        0x20 => None,
        0x21 => { advance(data, pos, 8); None }
        0x22..=0x26 => None,
        0x27 | 0x28 => { advance(data, pos, 5); None }
        0x29..=0x2c => None,
        0x2d => { advance(data, pos, 5); None }
        0x2e => { advance(data, pos, 4); None }
        0x2f => None,
        0x50 => { advance(data, pos, 1); advance(data, pos, 4); None }
        0x51 => {
            let slot = ru8(data, pos)?;
            let gold = ru32(data, pos)?;
            let lumber = ru32(data, pos)?;
            Some(Action::TransferResources { slot, gold, lumber })
        }
        0x60 => { advance(data, pos, 8); rzts(data, pos); None }
        0x61 => Some(Action::EscPressed),
        0x62 => { advance(data, pos, 12); None }
        0x63 => { advance(data, pos, 8); None }
        0x64 | 0x65 => {
            let object = rnettag(data, pos)?;
            Some(Action::TrackableHit { object })
        }
        0x66 => Some(Action::ChooseHeroSkillSubmenu),
        0x67 => Some(Action::EnterBuildingSubmenu),
        0x68 => {
            let pos2 = rvec2(data, pos)?;
            let duration = rf32(data, pos)?;
            Some(Action::AllyPing { pos: pos2, duration })
        }
        0x69 | 0x6a => { advance(data, pos, 16); None }
        0x6b => {
            let cache = rcache(data, pos)?;
            let value = ru32(data, pos)?;
            Some(Action::W3MMDStoreInt { cache, value })
        }
        0x6c => {
            let cache = rcache(data, pos)?;
            let value = rf32(data, pos)?;
            Some(Action::W3MMDStoreReal { cache, value })
        }
        0x6d => {
            let cache = rcache(data, pos)?;
            let value = ru8(data, pos)?;
            Some(Action::W3MMDStoreBool { cache, value })
        }
        0x6e => {
            let cache = rcache(data, pos)?;
            // read cache unit: unitId(4) + itemsCount(4) + items(12*count) + heroData(variable)
            read_cache_unit(data, pos)?;
            // We discard the unit data; just return None to not expose it
            Some(Action::W3MMDClearUnit { cache })
        }
        0x70 => { let cache = rcache(data, pos)?; Some(Action::W3MMDClearInt { cache }) }
        0x71 => { let cache = rcache(data, pos)?; Some(Action::W3MMDClearReal { cache }) }
        0x72 => { let cache = rcache(data, pos)?; Some(Action::W3MMDClearBool { cache }) }
        0x73 => { let cache = rcache(data, pos)?; Some(Action::W3MMDClearUnit { cache }) }
        0x75 => {
            let arrow_key = ru8(data, pos)?;
            Some(Action::ArrowKey { arrow_key })
        }
        0x76 => {
            let event_id = ru8(data, pos)?;
            let pos2 = rvec2(data, pos)?;
            let button = ru8(data, pos)?;
            Some(Action::MouseAction { event_id, pos: pos2, button })
        }
        0x77 => {
            let command_id = ru32(data, pos)?;
            let data_val = ru32(data, pos)?;
            let buff_len = ru32(data, pos)? as usize;
            advance(data, pos, buff_len);
            Some(Action::W3Api { command_id, data: data_val })
        }
        0x78 => {
            let identifier = rzts(data, pos);
            let value = rzts(data, pos);
            advance(data, pos, 4);
            Some(Action::BlzSync { identifier, value })
        }
        0x79 => {
            advance(data, pos, 8);
            let event_id = ru32(data, pos)?;
            let val = rf32(data, pos)?;
            let text = rzts(data, pos);
            Some(Action::CommandFrame { event_id, val, text })
        }
        0x7a => { advance(data, pos, 20); None }
        0x7b => { advance(data, pos, 16); None }
        0xa0 => { advance(data, pos, 14); None }
        0xa1 => { advance(data, pos, 9); None }
        _ => {
            // unknown action - can't reliably skip, so return None
            None
        }
    }
}

fn read_cache_unit(data: &[u8], pos: &mut usize) -> Option<()> {
    advance(data, pos, 4); // unitId
    let items_count = ru32(data, pos)? as usize;
    advance(data, pos, items_count * 12); // each item: 4+4+4 bytes
    // hero data
    advance(data, pos, 4*4 + 2*4); // xp, level, skillPoints, properNameId, str, strBonus(float)
    // agi, speedMod(f), cooldownMod(f), agiBonus(f)
    advance(data, pos, 4 + 4 + 4 + 4);
    // intel, intBonus(f)
    advance(data, pos, 4 + 4);
    let hero_abil_count = ru32(data, pos)? as usize;
    advance(data, pos, hero_abil_count * 8); // each ability: id(4) + level(4)
    // maxLife(f), maxMana(f), sight(f)
    advance(data, pos, 12);
    let damage_count = ru32(data, pos)? as usize;
    advance(data, pos, damage_count * 4);
    // defense(f), controlGroups(u16)
    advance(data, pos, 4 + 2);
    Some(())
}

// --- tiny helpers ---

#[inline]
fn advance(data: &[u8], pos: &mut usize, n: usize) {
    let new = (*pos + n).min(data.len());
    *pos = new;
}

#[inline]
fn ru8(data: &[u8], pos: &mut usize) -> Option<u8> {
    let v = *data.get(*pos)?;
    *pos += 1;
    Some(v)
}

#[inline]
fn ru16(data: &[u8], pos: &mut usize) -> Option<u16> {
    if *pos + 2 > data.len() { return None; }
    let v = u16::from_le_bytes([data[*pos], data[*pos+1]]);
    *pos += 2;
    Some(v)
}

#[inline]
fn ru32(data: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > data.len() { return None; }
    let v = u32::from_le_bytes([data[*pos], data[*pos+1], data[*pos+2], data[*pos+3]]);
    *pos += 4;
    Some(v)
}

#[inline]
fn rf32(data: &[u8], pos: &mut usize) -> Option<f32> {
    if *pos + 4 > data.len() { return None; }
    let v = f32::from_le_bytes([data[*pos], data[*pos+1], data[*pos+2], data[*pos+3]]);
    *pos += 4;
    Some(v)
}

#[inline]
fn rfourcc(data: &[u8], pos: &mut usize) -> Option<[u8; 4]> {
    if *pos + 4 > data.len() { return None; }
    let arr = [data[*pos], data[*pos+1], data[*pos+2], data[*pos+3]];
    *pos += 4;
    Some(arr)
}

#[inline]
fn rnettag(data: &[u8], pos: &mut usize) -> Option<[u32; 2]> {
    let a = ru32(data, pos)?;
    let b = ru32(data, pos)?;
    Some([a, b])
}

#[inline]
fn rvec2(data: &[u8], pos: &mut usize) -> Option<[f32; 2]> {
    let x = rf32(data, pos)?;
    let y = rf32(data, pos)?;
    Some([x, y])
}

fn rzts(data: &[u8], pos: &mut usize) -> String {
    let start = *pos;
    while *pos < data.len() && data[*pos] != 0 { *pos += 1; }
    let s = String::from_utf8_lossy(&data[start..*pos]).into_owned();
    if *pos < data.len() { *pos += 1; }
    s
}

fn rcache(data: &[u8], pos: &mut usize) -> Option<Cache> {
    let filename = rzts(data, pos);
    let mission_key = rzts(data, pos);
    let key = rzts(data, pos);
    Some(Cache { filename, mission_key, key })
}

/// Convert a 4-byte order ID to a formatted ItemID string (mirrors objectIdFormatter).
pub fn format_object_id(arr: &[u8; 4]) -> ObjectId {
    if arr[3] >= 0x41 && arr[3] <= 0x7a {
        // string encoded: reverse array and convert each byte to char
        let s: String = arr.iter().rev()
            .map(|&b| b as char)
            .collect();
        ObjectId::StringEncoded(s)
    } else {
        ObjectId::Alphanumeric(*arr)
    }
}

#[derive(Debug, Clone)]
pub enum ObjectId {
    StringEncoded(String),
    Alphanumeric([u8; 4]),
}

impl ObjectId {
    pub fn as_str(&self) -> Option<&str> {
        match self {
            ObjectId::StringEncoded(s) => Some(s.as_str()),
            ObjectId::Alphanumeric(_) => None,
        }
    }

    pub fn first_char(&self) -> Option<char> {
        match self {
            ObjectId::StringEncoded(s) => s.chars().next(),
            ObjectId::Alphanumeric(_) => None,
        }
    }

    pub fn is_alphanumeric(&self) -> bool {
        matches!(self, ObjectId::Alphanumeric(_))
    }
}
