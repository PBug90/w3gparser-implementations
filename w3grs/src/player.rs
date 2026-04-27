use std::collections::HashMap;
use crate::convert::player_color;
use crate::mappings::{ITEMS, UNITS, BUILDINGS, UPGRADES, ABILITY_TO_HERO};
use crate::detect_retraining::get_retraining_index;
use crate::infer_hero_abilities::{infer_hero_ability_levels, RetrainingSnapshot};
use crate::parser::action::{Action, ObjectId, format_object_id};
use serde::ser::{Serialize, Serializer, SerializeMap};

#[derive(Debug, Clone)]
pub enum AbilityOrderEntry {
    Ability { time: u32, value: String },
    Retraining { time: u32 },
}

impl Serialize for AbilityOrderEntry {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        match self {
            AbilityOrderEntry::Ability { time, value } => {
                let mut m = s.serialize_map(Some(3))?;
                m.serialize_entry("type", "ability")?;
                m.serialize_entry("time", time)?;
                m.serialize_entry("value", value)?;
                m.end()
            }
            AbilityOrderEntry::Retraining { time } => {
                let mut m = s.serialize_map(Some(2))?;
                m.serialize_entry("type", "retraining")?;
                m.serialize_entry("time", time)?;
                m.end()
            }
        }
    }
}

impl AbilityOrderEntry {
    pub fn time(&self) -> u32 {
        match self {
            AbilityOrderEntry::Ability { time, .. } => *time,
            AbilityOrderEntry::Retraining { time } => *time,
        }
    }
}

#[derive(Debug, Clone)]
pub struct HeroCollectorEntry {
    pub id: String,
    pub order: usize,
    pub ability_order: Vec<AbilityOrderEntry>,
}

#[derive(Debug, Clone)]
pub struct HeroInfo {
    pub id: String,
    pub level: u32,
    pub abilities: HashMap<String, u32>,
    pub retraining_history: Vec<RetrainingSnapshot>,
    pub ability_order: Vec<AbilityOrderEntry>,
}

impl Serialize for HeroInfo {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(5))?;
        m.serialize_entry("id", &self.id)?;
        m.serialize_entry("level", &self.level)?;
        m.serialize_entry("abilities", &self.abilities)?;
        m.serialize_entry("retrainingHistory", &self.retraining_history)?;
        m.serialize_entry("abilityOrder", &self.ability_order)?;
        m.end()
    }
}

#[derive(Debug, Clone)]
pub struct OrderEntry {
    pub id: String,
    pub ms: u32,
}

impl Serialize for OrderEntry {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(2))?;
        m.serialize_entry("id", &self.id)?;
        m.serialize_entry("ms", &self.ms)?;
        m.end()
    }
}

#[derive(Debug, Clone)]
pub struct Summary {
    pub summary: HashMap<String, u32>,
    pub order: Vec<OrderEntry>,
}

impl Summary {
    pub fn new() -> Self {
        Summary { summary: HashMap::new(), order: Vec::new() }
    }

    pub fn add(&mut self, id: &str, ms: u32) {
        *self.summary.entry(id.to_string()).or_insert(0) += 1;
        self.order.push(OrderEntry { id: id.to_string(), ms });
    }
}

impl Serialize for Summary {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(2))?;
        m.serialize_entry("summary", &self.summary)?;
        m.serialize_entry("order", &self.order)?;
        m.end()
    }
}

#[derive(Debug, Clone, Default)]
pub struct Actions {
    pub timed: Vec<u32>,
    pub assign_group: u32,
    pub right_click: u32,
    pub basic: u32,
    pub build_train: u32,
    pub ability: u32,
    pub item: u32,
    pub select: u32,
    pub remove_unit: u32,
    pub subgroup: u32,
    pub select_hotkey: u32,
    pub esc: u32,
}

impl Serialize for Actions {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(12))?;
        m.serialize_entry("timed", &self.timed)?;
        m.serialize_entry("assigngroup", &self.assign_group)?;
        m.serialize_entry("rightclick", &self.right_click)?;
        m.serialize_entry("basic", &self.basic)?;
        m.serialize_entry("buildtrain", &self.build_train)?;
        m.serialize_entry("ability", &self.ability)?;
        m.serialize_entry("item", &self.item)?;
        m.serialize_entry("select", &self.select)?;
        m.serialize_entry("removeunit", &self.remove_unit)?;
        m.serialize_entry("subgroup", &self.subgroup)?;
        m.serialize_entry("selecthotkey", &self.select_hotkey)?;
        m.serialize_entry("esc", &self.esc)?;
        m.end()
    }
}

#[derive(Debug, Clone, Default)]
pub struct GroupHotkey {
    pub assigned: u32,
    pub used: u32,
}

impl Serialize for GroupHotkey {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(2))?;
        m.serialize_entry("assigned", &self.assigned)?;
        m.serialize_entry("used", &self.used)?;
        m.end()
    }
}

#[derive(Debug, Clone)]
pub struct ResourceTransfer {
    pub slot: u8,
    pub player_id: u8,
    pub player_name: String,
    pub gold: u32,
    pub lumber: u32,
    pub ms_elapsed: u32,
}

impl Serialize for ResourceTransfer {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(6))?;
        m.serialize_entry("slot", &self.slot)?;
        m.serialize_entry("playerId", &self.player_id)?;
        m.serialize_entry("playerName", &self.player_name)?;
        m.serialize_entry("gold", &self.gold)?;
        m.serialize_entry("lumber", &self.lumber)?;
        m.serialize_entry("msElapsed", &self.ms_elapsed)?;
        m.end()
    }
}

#[derive(Debug, Clone)]
pub struct Player {
    pub id: u8,
    pub name: String,
    pub teamid: u8,
    pub color: String,
    pub race: String,       // declared race (slot flag)
    pub race_detected: String,
    pub units: Summary,
    pub upgrades: Summary,
    pub items: Summary,
    pub buildings: Summary,
    pub heroes: Vec<HeroInfo>,
    pub hero_collector: HashMap<String, HeroCollectorEntry>,
    pub hero_count: usize,
    pub actions: Actions,
    pub group_hotkeys: HashMap<u8, GroupHotkey>,
    pub resource_transfers: Vec<ResourceTransfer>,
    pub apm: u32,
    pub current_time_played: u32,

    // internal
    pub _currently_tracked_apm: u32,
    pub _last_action_was_deselect: bool,
    pub _last_retraining_time: u32,
}

/// Serializes only the public output fields, matching Go's PlayerOutput struct.
impl Serialize for Player {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(16))?;
        m.serialize_entry("id", &self.id)?;
        m.serialize_entry("name", &self.name)?;
        m.serialize_entry("teamid", &self.teamid)?;
        m.serialize_entry("color", &self.color)?;
        m.serialize_entry("race", &self.race)?;
        m.serialize_entry("raceDetected", &self.race_detected)?;
        m.serialize_entry("units", &self.units)?;
        m.serialize_entry("buildings", &self.buildings)?;
        m.serialize_entry("items", &self.items)?;
        m.serialize_entry("upgrades", &self.upgrades)?;
        m.serialize_entry("heroes", &self.heroes)?;
        m.serialize_entry("actions", &self.actions)?;
        // Sort group hotkey keys for deterministic output
        let mut gk: Vec<(&u8, &GroupHotkey)> = self.group_hotkeys.iter().collect();
        gk.sort_by_key(|(k, _)| *k);
        let gk_map: std::collections::BTreeMap<u8, &GroupHotkey> =
            gk.into_iter().map(|(k, v)| (*k, v)).collect();
        m.serialize_entry("groupHotkeys", &gk_map)?;
        m.serialize_entry("resourceTransfers", &self.resource_transfers)?;
        m.serialize_entry("apm", &self.apm)?;
        m.serialize_entry("currentTimePlayed", &self.current_time_played)?;
        m.end()
    }
}

impl Player {
    pub fn new(id: u8, name: String, teamid: u8, color: u8, race: String) -> Self {
        let mut group_hotkeys = HashMap::new();
        for i in 0u8..10 {
            group_hotkeys.insert(i, GroupHotkey::default());
        }
        Player {
            id, name, teamid,
            color: player_color(color),
            race,
            race_detected: String::new(),
            units: Summary::new(),
            upgrades: Summary::new(),
            items: Summary::new(),
            buildings: Summary::new(),
            heroes: Vec::new(),
            hero_collector: HashMap::new(),
            hero_count: 0,
            actions: Actions::default(),
            group_hotkeys,
            resource_transfers: Vec::new(),
            apm: 0,
            current_time_played: 0,
            _currently_tracked_apm: 0,
            _last_action_was_deselect: false,
            _last_retraining_time: 0,
        }
    }

    pub fn new_action_tracking_segment(&mut self, tracking_interval: u32) {
        let scaled = (self._currently_tracked_apm as f64 * (60000.0 / tracking_interval as f64)) as u32;
        self.actions.timed.push(scaled);
        self._currently_tracked_apm = 0;
    }

    pub fn detect_race_by_action_id(&mut self, action_id: &str) {
        if self.race_detected.is_empty() {
            if let Some(first) = action_id.chars().next() {
                match first {
                    'e' => self.race_detected = "N".to_string(),
                    'o' => self.race_detected = "O".to_string(),
                    'h' => self.race_detected = "H".to_string(),
                    'u' => self.race_detected = "U".to_string(),
                    _ => {}
                }
            }
        }
    }

    fn handle_stringencoded_item_id(&mut self, action_id: &str, gametime: u32) {
        if UNITS.contains_key(action_id) {
            self.units.add(action_id, gametime);
        } else if ITEMS.contains_key(action_id) {
            self.items.add(action_id, gametime);
        } else if BUILDINGS.contains_key(action_id) {
            self.buildings.add(action_id, gametime);
        } else if UPGRADES.contains_key(action_id) {
            self.upgrades.add(action_id, gametime);
        }
    }

    fn handle_hero_skill(&mut self, action_id: &str, gametime: u32) {
        let hero_id = match ABILITY_TO_HERO.get(action_id) {
            Some(&h) => h,
            None => return,
        };

        if !self.hero_collector.contains_key(hero_id) {
            self.hero_count += 1;
            self.hero_collector.insert(hero_id.to_string(), HeroCollectorEntry {
                id: hero_id.to_string(),
                order: self.hero_count,
                ability_order: Vec::new(),
            });
        }

        let entry = self.hero_collector.get_mut(hero_id).unwrap();
        entry.ability_order.push(AbilityOrderEntry::Ability {
            time: gametime,
            value: action_id.to_string(),
        });

        if self._last_retraining_time > 0 {
            let lrt = self._last_retraining_time;
            let idx = get_retraining_index(&entry.ability_order, lrt);
            if idx >= 0 {
                entry.ability_order.insert(idx as usize, AbilityOrderEntry::Retraining { time: lrt });
                self._last_retraining_time = 0;
            }
        }
    }

    pub fn handle_retraining(&mut self, gametime: u32) {
        self._last_retraining_time = gametime;
    }

    pub fn handle_0x10(&mut self, item_id: &ObjectId, gametime: u32) {
        match item_id {
            ObjectId::StringEncoded(s) => {
                match s.chars().next() {
                    Some('A') => self.handle_hero_skill(s, gametime),
                    Some('R') => self.handle_stringencoded_item_id(s, gametime),
                    Some('u') | Some('e') | Some('h') | Some('o') => {
                        if self.race_detected.is_empty() {
                            self.detect_race_by_action_id(s);
                        }
                        self.handle_stringencoded_item_id(s, gametime);
                    }
                    _ => self.handle_stringencoded_item_id(s, gametime),
                }
                // build/train action
                if s.chars().next().map_or(false, |c| c != '0') {
                    self.actions.build_train += 1;
                } else {
                    self.actions.ability += 1;
                }
            }
            ObjectId::Alphanumeric(arr) => {
                // first byte is '0' if arr[0] == 0x30? check JS: itemid.value[0] !== "0"
                // In alphanumeric case, value[0] is arr[0] as number
                // JS: if (itemid.value[0] !== "0") build_train++ else ability++
                // In our case we check arr[0] != 0 (numeric comparison, not char)
                // Actually JS code: value[0] is a number (array element)
                // "0" as string is falsy for === comparison to number
                // So numeric arr[0] != '0' (string) is always true since type mismatch.
                // This means for alphanumeric, always build_train++ (the condition always true).
                self.actions.build_train += 1;
            }
        }
        self._currently_tracked_apm += 1;
    }

    pub fn handle_0x11(&mut self, item_id: &ObjectId, gametime: u32) {
        self._currently_tracked_apm += 1;
        match item_id {
            ObjectId::Alphanumeric(arr) => {
                if arr[0] <= 0x19 && arr[1] == 0 {
                    self.actions.basic += 1;
                } else {
                    self.actions.ability += 1;
                }
            }
            ObjectId::StringEncoded(s) => {
                // JS: just calls handleStringencodedItemID, no ability++
                self.handle_stringencoded_item_id(s, gametime);
            }
        }
    }

    pub fn handle_0x12(&mut self, item_id: &ObjectId, gametime: u32) {
        match item_id {
            ObjectId::Alphanumeric(arr) => {
                if arr[0] == 0x03 && arr[1] == 0 {
                    self.actions.right_click += 1;
                } else if arr[0] <= 0x19 && arr[1] == 0 {
                    self.actions.basic += 1;
                } else {
                    self.actions.ability += 1;
                }
            }
            ObjectId::StringEncoded(s) => {
                self.handle_stringencoded_item_id(s, gametime);
                self.actions.ability += 1;
            }
        }
        self._currently_tracked_apm += 1;
    }

    pub fn handle_0x13(&mut self) {
        self.actions.item += 1;
        self._currently_tracked_apm += 1;
    }

    pub fn handle_0x14(&mut self, item_id: &ObjectId) {
        match item_id {
            ObjectId::Alphanumeric(arr) => {
                if arr[0] == 0x03 && arr[1] == 0 {
                    self.actions.right_click += 1;
                } else if arr[0] <= 0x19 && arr[1] == 0 {
                    self.actions.basic += 1;
                } else {
                    self.actions.ability += 1;
                }
            }
            ObjectId::StringEncoded(_) => {
                self.actions.ability += 1;
            }
        }
        self._currently_tracked_apm += 1;
    }

    pub fn handle_0x16(&mut self, select_mode: u8, is_apm: bool) {
        if is_apm {
            self.actions.select += 1;
            self._currently_tracked_apm += 1;
        }
    }

    pub fn handle_0x51(&mut self, slot: u8, player_id: u8, player_name: String, gold: u32, lumber: u32) {
        self.resource_transfers.push(ResourceTransfer {
            slot,
            player_id,
            player_name,
            gold,
            lumber,
            ms_elapsed: self.current_time_played,
        });
    }

    pub fn handle_other(&mut self, action: &Action) {
        match action {
            Action::AssignGroupHotkey { group_number, .. } => {
                self.actions.assign_group += 1;
                self._currently_tracked_apm += 1;
                let key = (group_number + 1) % 10;
                self.group_hotkeys.entry(key).or_default().assigned += 1;
            }
            Action::SelectGroupHotkey { group_number } => {
                self.actions.select_hotkey += 1;
                self._currently_tracked_apm += 1;
                let key = (group_number + 1) % 10;
                self.group_hotkeys.entry(key).or_default().used += 1;
            }
            Action::SelectGroundItem { .. } |
            Action::CancelHeroRevival { .. } |
            Action::ChooseHeroSkillSubmenu |
            Action::EnterBuildingSubmenu => {
                self._currently_tracked_apm += 1;
            }
            Action::RemoveUnitFromQueue { .. } => {
                self.actions.remove_unit += 1;
                self._currently_tracked_apm += 1;
            }
            Action::EscPressed => {
                self.actions.esc += 1;
                self._currently_tracked_apm += 1;
            }
            _ => {}
        }
    }

    pub fn determine_hero_levels_and_handle_retrainings(&mut self) {
        let mut collector_entries: Vec<&HeroCollectorEntry> = self.hero_collector.values().collect();
        collector_entries.sort_by_key(|e| e.order);

        let mut heroes = Vec::new();
        for entry in collector_entries {
            let inferred = infer_hero_ability_levels(&entry.ability_order);
            let level = inferred.final_abilities.values().sum();
            heroes.push(HeroInfo {
                id: entry.id.clone(),
                level,
                abilities: inferred.final_abilities,
                retraining_history: inferred.retraining_history,
                ability_order: entry.ability_order.clone(),
            });
        }
        self.heroes = heroes;
    }

    pub fn cleanup(&mut self, player_action_track_interval: u32) {
        self.new_action_tracking_segment(player_action_track_interval);
        let apm_sum: u32 = self.actions.timed.iter().sum();
        self.apm = if self.current_time_played == 0 {
            0
        } else {
            let minutes = self.current_time_played as f64 / 1000.0 / 60.0;
            (apm_sum as f64 / minutes).round() as u32
        };
        self.determine_hero_levels_and_handle_retrainings();
    }
}
