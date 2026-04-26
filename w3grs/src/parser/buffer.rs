pub struct BufParser<'a> {
    pub data: &'a [u8],
    pub pos: usize,
}

impl<'a> BufParser<'a> {
    pub fn new(data: &'a [u8]) -> Self {
        BufParser { data, pos: 0 }
    }

    pub fn remaining(&self) -> usize {
        self.data.len().saturating_sub(self.pos)
    }

    pub fn skip(&mut self, n: usize) {
        self.pos += n;
    }

    pub fn set_pos(&mut self, p: usize) {
        self.pos = p;
    }

    pub fn peek_u8(&self) -> Option<u8> {
        self.data.get(self.pos).copied()
    }

    pub fn read_u8(&mut self) -> Option<u8> {
        let v = self.data.get(self.pos).copied()?;
        self.pos += 1;
        Some(v)
    }

    pub fn read_u16_le(&mut self) -> Option<u16> {
        if self.pos + 2 > self.data.len() { return None; }
        let v = u16::from_le_bytes([self.data[self.pos], self.data[self.pos + 1]]);
        self.pos += 2;
        Some(v)
    }

    pub fn read_u32_le(&mut self) -> Option<u32> {
        if self.pos + 4 > self.data.len() { return None; }
        let v = u32::from_le_bytes([
            self.data[self.pos], self.data[self.pos+1],
            self.data[self.pos+2], self.data[self.pos+3],
        ]);
        self.pos += 4;
        Some(v)
    }

    pub fn read_f32_le(&mut self) -> Option<f32> {
        if self.pos + 4 > self.data.len() { return None; }
        let v = f32::from_le_bytes([
            self.data[self.pos], self.data[self.pos+1],
            self.data[self.pos+2], self.data[self.pos+3],
        ]);
        self.pos += 4;
        Some(v)
    }

    pub fn read_bytes(&mut self, n: usize) -> Option<&'a [u8]> {
        if self.pos + n > self.data.len() { return None; }
        let s = &self.data[self.pos..self.pos + n];
        self.pos += n;
        Some(s)
    }

    pub fn read_zero_term_string(&mut self) -> String {
        let start = self.pos;
        while self.pos < self.data.len() && self.data[self.pos] != 0 {
            self.pos += 1;
        }
        let s = String::from_utf8_lossy(&self.data[start..self.pos]).into_owned();
        if self.pos < self.data.len() {
            self.pos += 1; // consume null terminator
        }
        s
    }

    pub fn read_string_of_length_as_hex(&mut self, n: usize) -> String {
        if let Some(bytes) = self.read_bytes(n) {
            hex::encode(bytes)
        } else {
            String::new()
        }
    }

    pub fn read_string_of_length_utf8(&mut self, n: usize) -> String {
        if let Some(bytes) = self.read_bytes(n) {
            String::from_utf8_lossy(bytes).into_owned()
        } else {
            String::new()
        }
    }

    pub fn read_four_cc(&mut self) -> [u8; 4] {
        let mut arr = [0u8; 4];
        for a in &mut arr {
            *a = self.read_u8().unwrap_or(0);
        }
        arr
    }
}
