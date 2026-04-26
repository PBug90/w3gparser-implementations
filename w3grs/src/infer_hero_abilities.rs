use std::collections::HashMap;
use crate::player::AbilityOrderEntry;
use crate::mappings::ULTIMATES;
use serde::ser::{Serialize, Serializer, SerializeMap};

#[derive(Debug, Clone)]
pub struct HeroAbilities {
    pub final_abilities: HashMap<String, u32>,
    pub retraining_history: Vec<RetrainingSnapshot>,
}

#[derive(Debug, Clone)]
pub struct RetrainingSnapshot {
    pub time: u32,
    pub abilities: HashMap<String, u32>,
}

impl Serialize for RetrainingSnapshot {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        let mut m = s.serialize_map(Some(2))?;
        m.serialize_entry("time", &self.time)?;
        m.serialize_entry("abilities", &self.abilities)?;
        m.end()
    }
}

/// Mirrors inferHeroAbilityLevelsFromAbilityOrder.
pub fn infer_hero_ability_levels(ability_order: &[AbilityOrderEntry]) -> HeroAbilities {
    let mut abilities: HashMap<String, u32> = HashMap::new();
    let mut retraining_history = Vec::new();

    for entry in ability_order {
        match entry {
            AbilityOrderEntry::Ability { value, .. } => {
                let is_ultimate = ULTIMATES.contains(value.as_str());
                let current = abilities.entry(value.clone()).or_insert(0);
                if is_ultimate && *current == 1 {
                    continue;
                }
                if *current < 3 {
                    *current += 1;
                }
            }
            AbilityOrderEntry::Retraining { time } => {
                retraining_history.push(RetrainingSnapshot {
                    time: *time,
                    abilities: abilities.clone(),
                });
                abilities = HashMap::new();
            }
        }
    }

    HeroAbilities { final_abilities: abilities, retraining_history }
}
