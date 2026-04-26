package w3ggo

import "math"

type heroCollectorEntry struct {
	id           string
	order        int
	abilityOrder []AbilityOrderEntry
}

type player struct {
	id                  uint8
	name                string
	teamID              uint8
	color               string
	race                string
	raceDetected        string
	units               Summary
	upgrades            Summary
	items               Summary
	buildings           Summary
	heroes              []HeroInfo
	heroCollector       map[string]*heroCollectorEntry
	heroCount           int
	actions             Actions
	groupHotkeys        map[int]GroupHotkey
	resourceTransfers   []ResourceTransfer
	apm                 int
	currentTimePlayed   int

	// internal
	currentlyTrackedAPM     int
	lastActionWasDeselect   bool
	lastRetrainingTime      int
}

func newPlayer(id uint8, name string, teamID uint8, color uint8, race string) *player {
	groupHotkeys := make(map[int]GroupHotkey)
	for i := 0; i < 10; i++ {
		groupHotkeys[i] = GroupHotkey{}
	}
	return &player{
		id:            id,
		name:          name,
		teamID:        teamID,
		color:         playerColor(int(color)),
		race:          race,
		units:             newSummary(),
		upgrades:          newSummary(),
		items:             newSummary(),
		buildings:         newSummary(),
		heroCollector:     make(map[string]*heroCollectorEntry),
		groupHotkeys:      groupHotkeys,
		resourceTransfers: []ResourceTransfer{},
		actions:           Actions{Timed: []int{}},
	}
}

func (p *player) newActionTrackingSegment(trackingInterval int) {
	scaled := int(math.Round(float64(p.currentlyTrackedAPM) * (60000.0 / float64(trackingInterval))))
	p.actions.Timed = append(p.actions.Timed, scaled)
	p.currentlyTrackedAPM = 0
}

func (p *player) detectRaceByActionID(actionID string) {
	if p.raceDetected == "" && len(actionID) > 0 {
		switch actionID[0] {
		case 'e':
			p.raceDetected = "N"
		case 'o':
			p.raceDetected = "O"
		case 'h':
			p.raceDetected = "H"
		case 'u':
			p.raceDetected = "U"
		}
	}
}

func (p *player) handleStringEncodedItemID(actionID string, gametime int) {
	if _, ok := UNITS[actionID]; ok {
		p.units.add(actionID, gametime)
	} else if _, ok := ITEMS[actionID]; ok {
		p.items.add(actionID, gametime)
	} else if _, ok := BUILDINGS[actionID]; ok {
		p.buildings.add(actionID, gametime)
	} else if _, ok := UPGRADES[actionID]; ok {
		p.upgrades.add(actionID, gametime)
	}
}

func (p *player) handleHeroSkill(actionID string, gametime int) {
	heroID, ok := ABILITY_TO_HERO[actionID]
	if !ok {
		return
	}

	if _, exists := p.heroCollector[heroID]; !exists {
		p.heroCount++
		p.heroCollector[heroID] = &heroCollectorEntry{
			id:           heroID,
			order:        p.heroCount,
			abilityOrder: []AbilityOrderEntry{},
		}
	}

	entry := p.heroCollector[heroID]
	entry.abilityOrder = append(entry.abilityOrder, AbilityOrderEntry{
		Type:  "ability",
		Time:  gametime,
		Value: actionID,
	})

	if p.lastRetrainingTime > 0 {
		lrt := p.lastRetrainingTime
		idx := getRetrainingIndex(entry.abilityOrder, lrt)
		if idx >= 0 {
			// insert retraining marker at idx
			entry.abilityOrder = append(entry.abilityOrder, AbilityOrderEntry{})
			copy(entry.abilityOrder[idx+1:], entry.abilityOrder[idx:])
			entry.abilityOrder[idx] = AbilityOrderEntry{Type: "retraining", Time: lrt}
			p.lastRetrainingTime = 0
		}
	}
}

func (p *player) handleRetraining(gametime int) {
	p.lastRetrainingTime = gametime
}

func (p *player) handle0x10(itemID ObjectId, gametime int) {
	if itemID.isStringEncoded() {
		s := itemID.strVal
		if len(s) > 0 {
			switch s[0] {
			case 'A':
				p.handleHeroSkill(s, gametime)
			case 'R':
				p.handleStringEncodedItemID(s, gametime)
			case 'u', 'e', 'h', 'o':
				if p.raceDetected == "" {
					p.detectRaceByActionID(s)
				}
				p.handleStringEncodedItemID(s, gametime)
			default:
				p.handleStringEncodedItemID(s, gametime)
			}
		}
		// build/train check: s[0] != '0'
		if len(s) > 0 && s[0] != '0' {
			p.actions.BuildTrain++
		} else {
			p.actions.Ability++
		}
	} else {
		// alphanumeric: always build_train
		p.actions.BuildTrain++
	}
	p.currentlyTrackedAPM++
}

func (p *player) handle0x11(itemID ObjectId, gametime int) {
	p.currentlyTrackedAPM++
	if itemID.isStringEncoded() {
		p.handleStringEncodedItemID(itemID.strVal, gametime)
	} else {
		// alphanumeric
		arr := itemID.arrVal
		if arr[0] <= 0x19 && arr[1] == 0 {
			p.actions.Basic++
		} else {
			p.actions.Ability++
		}
	}
}

func (p *player) handle0x12(itemID ObjectId, gametime int) {
	if itemID.isStringEncoded() {
		p.handleStringEncodedItemID(itemID.strVal, gametime)
		p.actions.Ability++
	} else {
		arr := itemID.arrVal
		if arr[0] == 0x03 && arr[1] == 0 {
			p.actions.RightClick++
		} else if arr[0] <= 0x19 && arr[1] == 0 {
			p.actions.Basic++
		} else {
			p.actions.Ability++
		}
	}
	p.currentlyTrackedAPM++
}

func (p *player) handle0x13() {
	p.actions.Item++
	p.currentlyTrackedAPM++
}

func (p *player) handle0x14(itemID ObjectId) {
	if itemID.isStringEncoded() {
		p.actions.Ability++
	} else {
		arr := itemID.arrVal
		if arr[0] == 0x03 && arr[1] == 0 {
			p.actions.RightClick++
		} else if arr[0] <= 0x19 && arr[1] == 0 {
			p.actions.Basic++
		} else {
			p.actions.Ability++
		}
	}
	p.currentlyTrackedAPM++
}

func (p *player) handle0x16(selectMode uint8, isAPM bool) {
	if isAPM {
		p.actions.Select++
		p.currentlyTrackedAPM++
	}
}

func (p *player) handle0x51(playerID uint8, playerName string, gold uint32, lumber uint32) {
	p.resourceTransfers = append(p.resourceTransfers, ResourceTransfer{
		PlayerID:   int(playerID),
		PlayerName: playerName,
		Gold:       int(gold),
		Lumber:     int(lumber),
		MSElapsed:  p.currentTimePlayed,
	})
}

func (p *player) handleOther(action Action) {
	switch action.typ {
	case actAssignGroupHotkey:
		p.actions.AssignGroup++
		p.currentlyTrackedAPM++
		key := (int(action.groupNumber) + 1) % 10
		hk := p.groupHotkeys[key]
		hk.Assigned++
		p.groupHotkeys[key] = hk
	case actSelectGroupHotkey:
		p.actions.SelectHotkey++
		p.currentlyTrackedAPM++
		key := (int(action.groupNumber) + 1) % 10
		hk := p.groupHotkeys[key]
		hk.Used++
		p.groupHotkeys[key] = hk
	case actSelectGroundItem, actCancelHeroRevival, actChooseHeroSkillSubmenu, actEnterBuildingSubmenu:
		p.currentlyTrackedAPM++
	case actRemoveUnitFromQueue:
		p.actions.RemoveUnit++
		p.currentlyTrackedAPM++
	case actEscPressed:
		p.actions.ESC++
		p.currentlyTrackedAPM++
	}
}

func (p *player) determineHeroLevelsAndHandleRetrainings() {
	// sort heroes by order
	type entry struct {
		key string
		e   *heroCollectorEntry
	}
	entries := make([]entry, 0, len(p.heroCollector))
	for k, v := range p.heroCollector {
		entries = append(entries, entry{k, v})
	}
	// sort by order
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].e.order < entries[i].e.order {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	heroes := []HeroInfo{}
	for _, e := range entries {
		inferred := inferHeroAbilityLevels(e.e.abilityOrder)
		level := 0
		for _, v := range inferred.finalAbilities {
			level += v
		}
		heroes = append(heroes, HeroInfo{
			ID:               e.e.id,
			Level:            level,
			Abilities:        inferred.finalAbilities,
			RetrainingHistory: inferred.retrainingHistory,
			AbilityOrder:     e.e.abilityOrder,
		})
	}
	p.heroes = heroes
}

func (p *player) cleanup(playerActionTrackInterval int) {
	p.newActionTrackingSegment(playerActionTrackInterval)
	sum := 0
	for _, v := range p.actions.Timed {
		sum += v
	}
	if p.currentTimePlayed == 0 {
		p.apm = 0
	} else {
		minutes := float64(p.currentTimePlayed) / 1000.0 / 60.0
		p.apm = int(math.Round(float64(sum) / minutes))
	}
	p.determineHeroLevelsAndHandleRetrainings()
}

func (p *player) toOutput() PlayerOutput {
	out := PlayerOutput{
		ID:                int(p.id),
		Name:              p.name,
		TeamID:            int(p.teamID),
		Color:             p.color,
		Race:              p.race,
		RaceDetected:      p.raceDetected,
		Units:             p.units,
		Buildings:         p.buildings,
		Items:             p.items,
		Upgrades:          p.upgrades,
		Heroes:            p.heroes,
		Actions:           p.actions,
		GroupHotkeys:      p.groupHotkeys,
		ResourceTransfers: p.resourceTransfers,
		APM:               p.apm,
		CurrentTimePlayed: p.currentTimePlayed,
	}
	return out
}
