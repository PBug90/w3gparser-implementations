pub fn player_color(color: u8) -> String {
    match color {
        0 => "#ff0303",
        1 => "#0042ff",
        2 => "#1ce6b9",
        3 => "#540081",
        4 => "#fffc00",
        5 => "#fe8a0e",
        6 => "#20c000",
        7 => "#e55bb0",
        8 => "#959697",
        9 => "#7ebff1",
        10 => "#106246",
        11 => "#4a2a04",
        12 => "#9b0000",
        13 => "#0000c3",
        14 => "#00eaff",
        15 => "#be00fe",
        16 => "#ebcd87",
        17 => "#f8a48b",
        18 => "#bfff80",
        19 => "#dcb9eb",
        20 => "#282828",
        21 => "#ebf0ff",
        22 => "#00781e",
        23 => "#a46f33",
        _ => "000000",
    }
    .to_string()
}

pub fn game_version(version: u32) -> String {
    if version == 10030 {
        "1.30.2+".to_string()
    } else if version > 10030 && version < 10100 {
        let s = version.to_string();
        format!("1.{}", &s[s.len() - 2..])
    } else if version >= 10100 {
        let s = version.to_string();
        format!("2.{}", &s[s.len() - 2..])
    } else {
        format!("1.{}", version)
    }
}

pub fn map_filename(map_path: &str) -> String {
    // Matches: filename (optionally followed by (N)) with .w3x or .w3m extension
    // Using manual search to avoid regex dependency at runtime
    let path = map_path.replace('\\', "/");
    if let Some(last_sep) = path.rfind('/') {
        let filename = &path[last_sep + 1..];
        if filename.ends_with(".w3x") || filename.ends_with(".w3m") {
            return filename.to_string();
        }
    } else if path.ends_with(".w3x") || path.ends_with(".w3m") {
        return path.clone();
    }
    // fallback: regex-style search for the pattern
    let bytes = map_path.as_bytes();
    let len = bytes.len();
    // scan backwards for last \ or /
    let mut start = 0;
    for i in (0..len).rev() {
        if bytes[i] == b'\\' || bytes[i] == b'/' {
            start = i + 1;
            break;
        }
    }
    let candidate = &map_path[start..];
    if candidate.ends_with(".w3x") || candidate.ends_with(".w3m") {
        candidate.to_string()
    } else {
        String::new()
    }
}
