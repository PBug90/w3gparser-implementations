# w3grs

A fast Warcraft III replay (`.w3g`) parser written in Rust. Parses replay files into structured data: players, heroes, actions, resource transfers, chat, APM, map info, and more.

Output shape mirrors [w3gjs](https://github.com/PBug90/w3gjs) for cross-implementation compatibility.

## Install

Add to your `Cargo.toml`:

```toml
[dependencies]
w3grs = "0.1"
```

## Usage

### Parse from file path

```rust
use w3grs::parse_file;

fn main() {
    let out = parse_file("replay.w3g").expect("failed to parse replay");

    println!("Map:      {}", out.map.file);
    println!("Matchup:  {}", out.matchup);
    println!("Duration: {}s", out.duration / 1000);

    for player in &out.players {
        println!(
            "  {} ({}) — APM {}",
            player.name, player.race, player.apm
        );
    }
}
```

### Parse from bytes

```rust
use w3grs::parse_bytes;

let data = std::fs::read("replay.w3g").unwrap();
let out = parse_bytes(&data).expect("failed to parse replay");
```

### Serialize to JSON

`ParserOutput` implements `serde::Serialize`:

```rust
use w3grs::parse_file;

let out = parse_file("replay.w3g").unwrap();
let json = serde_json::to_string_pretty(&out).unwrap();
println!("{}", json);
```

## CLI

The crate also ships a binary that prints JSON to stdout:

```sh
cargo install w3grs
w3grs replay.w3g
```

## License

MIT
