/// Minimal protobuf decoder for Reforged player metadata.

pub struct ProtoDecoder<'a> {
    data: &'a [u8],
    pos: usize,
}

impl<'a> ProtoDecoder<'a> {
    pub fn new(data: &'a [u8]) -> Self {
        ProtoDecoder { data, pos: 0 }
    }

    fn remaining(&self) -> usize {
        self.data.len().saturating_sub(self.pos)
    }

    pub fn read_varint(&mut self) -> Option<u64> {
        let mut result: u64 = 0;
        let mut shift = 0u32;
        loop {
            if self.pos >= self.data.len() { return None; }
            let byte = self.data[self.pos];
            self.pos += 1;
            result |= ((byte & 0x7F) as u64) << shift;
            if byte & 0x80 == 0 { break; }
            shift += 7;
            if shift >= 64 { return None; }
        }
        Some(result)
    }

    fn read_length_delimited(&mut self) -> Option<&'a [u8]> {
        let len = self.read_varint()? as usize;
        if self.pos + len > self.data.len() { return None; }
        let s = &self.data[self.pos..self.pos + len];
        self.pos += len;
        Some(s)
    }

    fn skip_field(&mut self, wire_type: u64) {
        match wire_type {
            0 => { self.read_varint(); }
            1 => { self.pos += 8; }
            2 => { self.read_length_delimited(); }
            5 => { self.pos += 4; }
            _ => {}
        }
    }

    /// Decode a flat player message: field1=playerId(varint), field2=battleTag(str), field3=clan(str)
    pub fn decode_player_flat(&mut self) -> Option<(u32, String, String)> {
        let mut player_id: u32 = 0;
        let mut battle_tag = String::new();
        let mut clan = String::new();

        while self.remaining() > 0 {
            let tag = self.read_varint()?;
            let field_number = tag >> 3;
            let wire_type = tag & 0x7;
            match field_number {
                1 => { player_id = self.read_varint()? as u32; }
                2 => {
                    let bytes = self.read_length_delimited()?;
                    battle_tag = String::from_utf8_lossy(bytes).into_owned();
                }
                3 => {
                    let bytes = self.read_length_delimited()?;
                    clan = String::from_utf8_lossy(bytes).into_owned();
                }
                _ => { self.skip_field(wire_type); }
            }
        }
        Some((player_id, battle_tag, clan))
    }

    /// Decode nested format: outer blob has repeated field-1 length-delimited sub-messages
    pub fn decode_player_nested(&mut self) -> Vec<(u32, String, String)> {
        let mut players = Vec::new();
        while self.remaining() > 0 {
            let tag = match self.read_varint() { Some(v) => v, None => break };
            let field_number = tag >> 3;
            let wire_type = tag & 0x7;
            if field_number == 1 && wire_type == 2 {
                if let Some(sub_data) = self.read_length_delimited() {
                    let mut sub = ProtoDecoder::new(sub_data);
                    if let Some(player) = sub.decode_player_flat() {
                        players.push(player);
                    }
                }
            } else {
                self.skip_field(wire_type);
            }
        }
        players
    }
}

/// Decode player metadata blob (subtype 0x03).
/// Returns list of (player_id, battle_tag, clan).
pub fn decode_player_metadata(data: &[u8]) -> Vec<(u32, String, String)> {
    if data.is_empty() { return Vec::new(); }

    let first_tag = {
        let mut d = ProtoDecoder::new(data);
        d.read_varint().unwrap_or(0)
    };
    let field_number = first_tag >> 3;
    let wire_type = first_tag & 0x7;

    if field_number == 1 && wire_type == 2 {
        // Nested format (build 6105+)
        let mut d = ProtoDecoder::new(data);
        d.decode_player_nested()
    } else {
        // Flat format
        let mut d = ProtoDecoder::new(data);
        if let Some(player) = d.decode_player_flat() {
            vec![player]
        } else {
            Vec::new()
        }
    }
}
