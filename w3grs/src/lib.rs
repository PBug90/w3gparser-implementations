pub mod convert;
pub mod mappings;
pub mod types;
pub mod parser;
pub mod player;
pub mod detect_retraining;
pub mod infer_hero_abilities;
pub mod replay;

pub use replay::{parse_file, parse_bytes, parse_file_with_handler, parse_bytes_with_handler};
pub use types::{ParserOutput, BasicReplayInfo, BasicPlayerInfo};
pub use parser::action::{Action, Cache, ObjectId, format_object_id};
pub use parser::game_data::{GameDataBlock, TimeslotBlock, CommandBlock, ChatMessageBlock, LeaveGameBlock};

/// Implement this trait to receive events during replay parsing.
///
/// All methods have default no-op implementations; override only what you need.
///
/// # Example
///
/// ```no_run
/// struct MyHandler;
///
/// impl w3grs::ReplayHandler for MyHandler {
///     fn on_basic_replay_information(&mut self, info: &w3grs::BasicReplayInfo) {
///         println!("Map: {}", info.map.file);
///     }
///     fn on_gamedatablock(&mut self, block: &w3grs::GameDataBlock) {
///         if let w3grs::GameDataBlock::Timeslot(ts) = block {
///             println!("{} command(s) at +{}ms", ts.command_blocks.len(), ts.time_increment);
///         }
///     }
/// }
/// ```
pub trait ReplayHandler {
    /// Called once after metadata is parsed, before game data blocks are processed.
    fn on_basic_replay_information(&mut self, _info: &BasicReplayInfo) {}
    /// Called for each game data block before the parser processes it internally.
    fn on_gamedatablock(&mut self, _block: &GameDataBlock) {}
}
