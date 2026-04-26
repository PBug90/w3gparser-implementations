use super::buffer::BufParser;
use flate2::read::ZlibDecoder;
use std::io::Read;

#[derive(Debug, Clone)]
pub struct Header {
    pub compressed_size: u32,
    pub header_version: String,
    pub decompressed_size: u32,
    pub block_count: u32,
}

#[derive(Debug, Clone)]
pub struct SubHeader {
    pub game_identifier: String,
    pub version: u32,
    pub build_no: u16,
    pub replay_length_ms: u32,
}

#[derive(Debug)]
pub struct DataBlock {
    pub block_size: u16,
    pub block_decompressed_size: u16,
    pub block_content: Vec<u8>,
}

pub fn parse(input: &[u8]) -> Option<(Header, SubHeader, Vec<DataBlock>)> {
    let mut p = BufParser::new(input);

    // Find "Warcraft III recorded game\0"
    let magic = b"Warcraft III recorded game";
    let start = find_subsequence(input, magic)?;
    p.set_pos(start);
    p.read_zero_term_string(); // consume the null-terminated magic string
    p.skip(4); // unknown 4 bytes

    let compressed_size = p.read_u32_le()?;
    let header_version = p.read_string_of_length_as_hex(4);
    let decompressed_size = p.read_u32_le()?;
    let block_count = p.read_u32_le()?;

    let header = Header {
        compressed_size,
        header_version,
        decompressed_size,
        block_count,
    };

    // Subheader
    let game_identifier = p.read_string_of_length_utf8(4);
    let version = p.read_u32_le()?;
    let build_no = p.read_u16_le()?;
    p.skip(2);
    let replay_length_ms = p.read_u32_le()?;
    p.skip(4);

    let subheader = SubHeader {
        game_identifier,
        version,
        build_no,
        replay_length_ms,
    };

    let is_reforged = build_no >= 6089;

    let mut blocks = Vec::new();
    while p.remaining() > 0 {
        let block_size = match p.read_u16_le() { Some(v) => v, None => break };
        if is_reforged { p.skip(2); }
        let block_decompressed_size = match p.read_u16_le() { Some(v) => v, None => break };
        if is_reforged { p.skip(6); } else { p.skip(4); }
        let content = match p.read_bytes(block_size as usize) { Some(b) => b.to_vec(), None => break };

        if block_decompressed_size == 8192 {
            blocks.push(DataBlock {
                block_size,
                block_decompressed_size,
                block_content: content,
            });
        }
    }

    Some((header, subheader, blocks))
}

pub fn decompress_blocks(blocks: &[DataBlock]) -> Vec<u8> {
    let mut result = Vec::new();
    for block in blocks {
        if block.block_content.is_empty() {
            continue;
        }
        let mut decoder = ZlibDecoder::new(block.block_content.as_slice());
        let mut buf = Vec::new();
        if decoder.read_to_end(&mut buf).is_ok() && !buf.is_empty() {
            result.extend_from_slice(&buf);
        }
    }
    result
}

fn find_subsequence(haystack: &[u8], needle: &[u8]) -> Option<usize> {
    haystack.windows(needle.len()).position(|w| w == needle)
}
