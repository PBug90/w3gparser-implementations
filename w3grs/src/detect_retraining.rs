use crate::player::AbilityOrderEntry;

const DETECTION_TIME_RANGE: u32 = 60 * 1000;

/// Mirrors getRetrainingIndex from detectRetraining.ts.
/// Returns the index in ability_order where a retraining marker should be inserted,
/// or -1 if not detected.
pub fn get_retraining_index(
    ability_order: &[AbilityOrderEntry],
    time_of_tome_purchase: u32,
) -> i32 {
    if ability_order.len() < 3 {
        return -1;
    }

    let mut candidate = &ability_order[0];
    let mut candidate_index: i32 = 0;
    let mut abilities_in_range: usize = 0;

    for i in 1..ability_order.len() {
        let entry = &ability_order[i];
        if entry.time().saturating_sub(candidate.time()) < DETECTION_TIME_RANGE {
            abilities_in_range += 1;
        } else {
            abilities_in_range = 0;
            candidate = entry;
            candidate_index = i as i32;
        }

        if abilities_in_range == 2
            && candidate.time().saturating_sub(time_of_tome_purchase) <= DETECTION_TIME_RANGE
        {
            return candidate_index;
        }
    }
    -1
}
