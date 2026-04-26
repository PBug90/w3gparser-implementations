use crate::player::{Player, HeroInfo, Summary, ResourceTransfer, GroupHotkey, Actions, OrderEntry};
use crate::infer_hero_abilities::RetrainingSnapshot;
use std::collections::HashMap;
use serde::ser::{Serializer, Serialize};

#[derive(Debug, Clone)]
pub struct MapInfo {
    pub path: String,
    pub file: String,
    pub checksum: String,
    pub checksum_sha1: String,
}

impl Serialize for MapInfo {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        use serde::ser::SerializeMap;
        let mut m = s.serialize_map(Some(4))?;
        m.serialize_entry("path", &self.path)?;
        m.serialize_entry("file", &self.file)?;
        m.serialize_entry("checksum", &self.checksum)?;
        m.serialize_entry("checksumSha1", &self.checksum_sha1)?;
        m.end()
    }
}

#[derive(Debug, Clone)]
pub struct ApmConfig {
    pub tracking_interval: u32,
}

impl Serialize for ApmConfig {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        use serde::ser::SerializeMap;
        let mut m = s.serialize_map(Some(1))?;
        m.serialize_entry("trackingInterval", &self.tracking_interval)?;
        m.end()
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum ObserverMode {
    OnDefeat,
    Full,
    Referees,
    None,
}

impl ObserverMode {
    pub fn as_str(&self) -> &'static str {
        match self {
            ObserverMode::OnDefeat => "ON_DEFEAT",
            ObserverMode::Full => "FULL",
            ObserverMode::Referees => "REFEREES",
            ObserverMode::None => "NONE",
        }
    }
}

impl Serialize for ObserverMode {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        s.serialize_str(self.as_str())
    }
}

#[derive(Debug, Clone)]
pub struct Settings {
    pub observer_mode: ObserverMode,
    pub referees: bool,
    pub fixed_teams: bool,
    pub full_shared_unit_control: bool,
    pub always_visible: bool,
    pub hide_terrain: bool,
    pub map_explored: bool,
    pub teams_together: bool,
    pub random_hero: bool,
    pub random_races: bool,
    pub speed: u8,
}

impl Serialize for Settings {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        use serde::ser::SerializeMap;
        let mut m = s.serialize_map(Some(11))?;
        m.serialize_entry("observerMode", &self.observer_mode)?;
        m.serialize_entry("referees", &self.referees)?;
        m.serialize_entry("fixedTeams", &self.fixed_teams)?;
        m.serialize_entry("fullSharedUnitControl", &self.full_shared_unit_control)?;
        m.serialize_entry("alwaysVisible", &self.always_visible)?;
        m.serialize_entry("hideTerrain", &self.hide_terrain)?;
        m.serialize_entry("mapExplored", &self.map_explored)?;
        m.serialize_entry("teamsTogether", &self.teams_together)?;
        m.serialize_entry("randomHero", &self.random_hero)?;
        m.serialize_entry("randomRaces", &self.random_races)?;
        m.serialize_entry("speed", &self.speed)?;
        m.end()
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum ChatMode {
    All,
    Team,
    Observers,
    Private,
}

impl ChatMode {
    pub fn as_str(&self) -> &'static str {
        match self {
            ChatMode::All => "All",
            ChatMode::Team => "Team",
            ChatMode::Observers => "Observers",
            ChatMode::Private => "Private",
        }
    }
}

impl Serialize for ChatMode {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        s.serialize_str(self.as_str())
    }
}

#[derive(Debug, Clone)]
pub struct ChatMessage {
    pub player_name: String,
    pub player_id: u8,
    pub mode: ChatMode,
    pub time_ms: u32,
    pub message: String,
}

impl Serialize for ChatMessage {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        use serde::ser::SerializeMap;
        let mut m = s.serialize_map(Some(5))?;
        m.serialize_entry("playerName", &self.player_name)?;
        m.serialize_entry("playerId", &self.player_id)?;
        m.serialize_entry("mode", &self.mode)?;
        m.serialize_entry("timeMS", &self.time_ms)?;
        m.serialize_entry("message", &self.message)?;
        m.end()
    }
}

/// Final output of the parser - mirrors ParserOutput from w3gjs
#[derive(Debug)]
pub struct ParserOutput {
    pub id: String,
    pub gamename: String,
    pub randomseed: u32,
    pub start_spots: u8,
    pub observers: Vec<String>,
    pub players: Vec<Player>,
    pub matchup: String,
    pub creator: String,
    pub game_type: String,
    pub chat: Vec<ChatMessage>,
    pub apm: ApmConfig,
    pub map: MapInfo,
    pub build_number: u16,
    pub version: String,
    pub duration: u32,
    pub expansion: bool,
    pub parse_time: u64,
    pub winning_team_id: i32,
    pub settings: Settings,
}

impl Serialize for ParserOutput {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        use serde::ser::SerializeMap;
        let mut m = s.serialize_map(Some(19))?;
        m.serialize_entry("id", &self.id)?;
        m.serialize_entry("gamename", &self.gamename)?;
        m.serialize_entry("randomseed", &self.randomseed)?;
        m.serialize_entry("startSpots", &self.start_spots)?;
        m.serialize_entry("observers", &self.observers)?;
        m.serialize_entry("players", &self.players)?;
        m.serialize_entry("matchup", &self.matchup)?;
        m.serialize_entry("creator", &self.creator)?;
        m.serialize_entry("type", &self.game_type)?;
        m.serialize_entry("chat", &self.chat)?;
        m.serialize_entry("apm", &self.apm)?;
        m.serialize_entry("map", &self.map)?;
        m.serialize_entry("buildNumber", &self.build_number)?;
        m.serialize_entry("version", &self.version)?;
        m.serialize_entry("duration", &self.duration)?;
        m.serialize_entry("expansion", &self.expansion)?;
        m.serialize_entry("parseTime", &self.parse_time)?;
        m.serialize_entry("winningTeamId", &self.winning_team_id)?;
        m.serialize_entry("settings", &self.settings)?;
        m.end()
    }
}
