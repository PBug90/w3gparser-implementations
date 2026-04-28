//! Basic example: parse a replay and print a summary.
//!
//! Demonstrates both the simple API and the event-based handler API.
//!
//! Run with:
//!   cargo run --example parse_replay -- path/to/replay.w3g

use w3grs::{parse_file, parse_file_with_handler, GameDataBlock, ReplayHandler, BasicReplayInfo};

// --- Simple API ---

fn print_summary(path: &str) {
    let out = parse_file(path).unwrap_or_else(|| {
        eprintln!("error: failed to parse {}", path);
        std::process::exit(1);
    });

    println!("Map:      {}", out.map.file);
    println!("Matchup:  {}", out.matchup);
    println!("Duration: {}s", out.duration / 1000);
    println!("Players:");
    for p in &out.players {
        println!("  [{:>7}] {} — APM {}", p.race, p.name, p.apm);
        if !p.resource_transfers.is_empty() {
            println!("    {} resource transfer(s)", p.resource_transfers.len());
        }
    }
    if !out.chat.is_empty() {
        println!("Chat ({} messages):", out.chat.len());
        for msg in &out.chat {
            println!("  <{}> {}", msg.player_name, msg.message);
        }
    }
}

// --- Event-based API ---

struct StatsHandler {
    timeslots: u32,
    total_actions: u32,
}

impl ReplayHandler for StatsHandler {
    fn on_basic_replay_information(&mut self, info: &BasicReplayInfo) {
        println!("[event] map: {}  players: {}", info.map.file, info.players.len());
    }

    fn on_gamedatablock(&mut self, block: &GameDataBlock) {
        if let GameDataBlock::Timeslot(ts) = block {
            self.timeslots += 1;
            for cmd in &ts.command_blocks {
                self.total_actions += cmd.actions.len() as u32;
            }
        }
    }
}

fn count_with_handler(path: &str) {
    let mut handler = StatsHandler { timeslots: 0, total_actions: 0 };
    parse_file_with_handler(path, &mut handler);
    println!("[event] {} timeslots, {} actions", handler.timeslots, handler.total_actions);
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 2 {
        eprintln!("usage: cargo run --example parse_replay -- <replay.w3g>");
        std::process::exit(1);
    }

    println!("=== simple API ===");
    print_summary(&args[1]);

    println!("\n=== handler API ===");
    count_with_handler(&args[1]);
}
