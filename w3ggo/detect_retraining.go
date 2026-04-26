package w3ggo

const detectionTimeRange = 60 * 1000

// getRetrainingIndex mirrors getRetrainingIndex from detectRetraining.ts.
// Returns the index in abilityOrder where a retraining marker should be inserted,
// or -1 if not detected.
func getRetrainingIndex(abilityOrder []AbilityOrderEntry, timeOfTomePurchase int) int {
	if len(abilityOrder) < 3 {
		return -1
	}

	candidate := &abilityOrder[0]
	candidateIndex := 0
	abilitiesInRange := 0

	for i := 1; i < len(abilityOrder); i++ {
		entry := &abilityOrder[i]
		diff := entry.Time - candidate.Time
		if diff < 0 {
			diff = -diff
		}
		if diff < detectionTimeRange {
			abilitiesInRange++
		} else {
			abilitiesInRange = 0
			candidate = entry
			candidateIndex = i
		}

		candidateDiff := candidate.Time - timeOfTomePurchase
		if candidateDiff < 0 {
			candidateDiff = -candidateDiff
		}
		if abilitiesInRange == 2 && candidateDiff <= detectionTimeRange {
			return candidateIndex
		}
	}
	return -1
}
