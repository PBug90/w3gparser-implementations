package w3ggo

type heroAbilities struct {
	finalAbilities    map[string]int
	retrainingHistory []RetrainingSnapshot
}

// inferHeroAbilityLevels mirrors inferHeroAbilityLevelsFromAbilityOrder.
func inferHeroAbilityLevels(abilityOrder []AbilityOrderEntry) heroAbilities {
	abilities := make(map[string]int)
	retrainingHistory := []RetrainingSnapshot{}

	for _, entry := range abilityOrder {
		switch entry.Type {
		case "ability":
			value := entry.Value
			isUltimate := ULTIMATES[value]
			current := abilities[value]
			if isUltimate && current == 1 {
				continue
			}
			if current < 3 {
				abilities[value] = current + 1
			}
		case "retraining":
			snapshot := make(map[string]int)
			for k, v := range abilities {
				snapshot[k] = v
			}
			retrainingHistory = append(retrainingHistory, RetrainingSnapshot{
				Time:      entry.Time,
				Abilities: snapshot,
			})
			abilities = make(map[string]int)
		}
	}

	return heroAbilities{
		finalAbilities:    abilities,
		retrainingHistory: retrainingHistory,
	}
}
