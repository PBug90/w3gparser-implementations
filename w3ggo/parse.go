package w3ggo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"time"
)

type noopEventHandler struct{}

func (noopEventHandler) OnBasicReplayInformation(_ BasicReplayInfo) {}
func (noopEventHandler) OnGameDataBlock(_ GameDataBlock)            {}

// ParseFile parses a WC3 replay from a file path.
func ParseFile(path string) (*ParserOutput, error) {
	return ParseFileWithHandler(path, noopEventHandler{})
}

// ParseFileWithHandler parses a WC3 replay from a file path, calling h for
// each event as the replay is processed.
func ParseFileWithHandler(path string, h EventHandler) (*ParserOutput, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := ParseBytesWithHandler(data, h)
	if out == nil {
		return nil, fmt.Errorf("parse failed")
	}
	return out, nil
}

// ParseBytes parses a WC3 replay from raw bytes.
func ParseBytes(input []byte) *ParserOutput {
	return ParseBytesWithHandler(input, noopEventHandler{})
}

// ParseBytesWithHandler parses a WC3 replay from raw bytes, calling h for
// each event as the replay is processed.
func ParseBytesWithHandler(input []byte, h EventHandler) *ParserOutput {
	startTime := time.Now()

	header, subHeader, blocks := parseRaw(input)
	if header == nil {
		return nil
	}
	uncompressed := decompressBlocks(blocks)
	meta := parseMetadata(uncompressed)
	if meta == nil {
		return nil
	}

	gameDataBlocks := parseGameData(meta.gameData, meta.isPost202)

	// Build player records map
	type tempPlayerRecord struct {
		playerID   uint8
		playerName string
	}
	tempPlayers := make(map[uint8]*tempPlayerRecord)
	for _, pr := range meta.playerRecords {
		tempPlayers[pr.playerID] = &tempPlayerRecord{
			playerID:   pr.playerID,
			playerName: pr.playerName,
		}
	}

	// Apply reforged player names
	for _, extra := range meta.reforgedPlayerMetadata {
		if rec, ok := tempPlayers[uint8(extra.playerID)]; ok {
			rec.playerName = extra.name
		}
	}

	players := make(map[uint8]*player)
	slotToPlayerID := make(map[int]uint8)

	for slotIndex, slot := range meta.slotRecords {
		if slot.slotStatus > 1 {
			slotToPlayerID[slotIndex] = slot.playerID

			name := "Computer"
			if rec, ok := tempPlayers[slot.playerID]; ok {
				name = rec.playerName
			}

			race := raceFlagToString(slot.raceFlag)
			players[slot.playerID] = newPlayer(slot.playerID, name, slot.teamID, slot.color, race)
		}
	}

	knownPlayerIDs := make(map[uint8]bool)
	for id := range players {
		knownPlayerIDs[id] = true
	}

	// Emit OnBasicReplayInformation before game data processing
	basicPlayers := make([]BasicPlayerInfo, 0, len(meta.slotRecords))
	for _, slot := range meta.slotRecords {
		if slot.slotStatus > 1 {
			name := "Computer"
			if rec, ok := tempPlayers[slot.playerID]; ok {
				name = rec.playerName
			}
			basicPlayers = append(basicPlayers, BasicPlayerInfo{
				PlayerID: slot.playerID,
				Name:     name,
				TeamID:   slot.teamID,
				Color:    slot.color,
				Race:     raceFlagToString(slot.raceFlag),
			})
		}
	}
	h.OnBasicReplayInformation(BasicReplayInfo{
		BuildNumber: uint32(subHeader.buildNo),
		Version:     gameVersion(int(subHeader.version)),
		GameName:    meta.gameName,
		RandomSeed:  meta.randomSeed,
		StartSpots:  meta.startSpotCount,
		Map: MapInfo{
			Path:         meta.mapMeta.mapName,
			File:         mapFilename(meta.mapMeta.mapName),
			Checksum:     meta.mapMeta.mapChecksum,
			ChecksumSha1: meta.mapMeta.mapChecksumSha1,
		},
		Players:   basicPlayers,
		Expansion: subHeader.gameIdentifier == "PX3W",
	})

	const playerActionTrackInterval = 60000

	totalTimeTracker := 0
	timeSegmentTracker := 0
	chatLog := []ChatMessage{}
	var leaveEvents []leaveGameBlock

	for _, block := range gameDataBlocks {
		// Emit OnGameDataBlock before internal processing
		switch block.typ {
		case gdTimeslot:
			ts := block.timeslot
			cmdBlocks := make([]CommandBlock, len(ts.commandBlocks))
			for i, cmd := range ts.commandBlocks {
				cmdBlocks[i] = CommandBlock{PlayerID: cmd.playerID, Actions: cmd.actions}
			}
			h.OnGameDataBlock(TimeslotEvent{TimeIncrement: ts.timeIncrement, CommandBlocks: cmdBlocks})
		case gdChatMessage:
			chat := block.chat
			h.OnGameDataBlock(ChatEvent{PlayerID: chat.playerID, Mode: chat.mode, Message: chat.message})
		case gdLeaveGame:
			lg := block.leaveGame
			h.OnGameDataBlock(LeaveGameEvent{PlayerID: lg.playerID, Reason: lg.reason, Result: lg.result})
		}

		switch block.typ {
		case gdTimeslot:
			ts := block.timeslot
			totalTimeTracker += int(ts.timeIncrement)
			timeSegmentTracker += int(ts.timeIncrement)

			if timeSegmentTracker > playerActionTrackInterval {
				for _, p := range players {
					p.newActionTrackingSegment(playerActionTrackInterval)
				}
				timeSegmentTracker = 0
			}

			for _, cmd := range ts.commandBlocks {
				if !knownPlayerIDs[cmd.playerID] {
					continue
				}
				p := players[cmd.playerID]
				p.currentTimePlayed = totalTimeTracker
				p.lastActionWasDeselect = false

				for _, action := range cmd.actions {
					processAction(action, p, totalTimeTracker, slotToPlayerID)
				}
			}

		case gdChatMessage:
			chat := block.chat
			if p, ok := players[chat.playerID]; ok {
				mode := chatModeFromInt(chat.mode)
				chatLog = append(chatLog, ChatMessage{
					PlayerName: p.name,
					PlayerID:   int(chat.playerID),
					Mode:       mode,
					TimeMS:     totalTimeTracker,
					Message:    chat.message,
				})
			}

		case gdLeaveGame:
			leaveEvents = append(leaveEvents, block.leaveGame)
		}
	}

	// Build player names map for resource transfer second pass
	playerNames := make(map[uint8]string)
	for id, p := range players {
		playerNames[id] = p.name
	}
	for _, p := range players {
		for i := range p.resourceTransfers {
			pid := uint8(p.resourceTransfers[i].PlayerID)
			if name, ok := playerNames[pid]; ok {
				p.resourceTransfers[i].PlayerName = name
			}
		}
	}

	versionNum := subHeader.version

	// Separate observers from players (iterate in sorted ID order for determinism)
	observers := []string{}
	finalPlayers := make(map[uint8]*player)

	allPlayerIDs := make([]uint8, 0, len(players))
	for id := range players {
		allPlayerIDs = append(allPlayerIDs, id)
	}
	sort.Slice(allPlayerIDs, func(i, j int) bool { return allPlayerIDs[i] < allPlayerIDs[j] })

	for _, id := range allPlayerIDs {
		p := players[id]
		isObs := isObserver(p, versionNum)
		p.cleanup(playerActionTrackInterval)
		if isObs {
			observers = append(observers, p.name)
		} else {
			finalPlayers[id] = p
		}
	}

	// Determine matchup and gametype
	gametype, matchup := determineMatchup(finalPlayers)

	// Determine winning team
	winningTeamID := determineWinningTeam(gametype, leaveEvents, finalPlayers)

	// Generate ID
	id := generateID(meta.randomSeed, finalPlayers, meta.gameName)

	// Sort players by teamid then id
	sortedIDs := make([]uint8, 0, len(finalPlayers))
	for id := range finalPlayers {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Slice(sortedIDs, func(i, j int) bool {
		a := finalPlayers[sortedIDs[i]]
		b := finalPlayers[sortedIDs[j]]
		if a.teamID != b.teamID {
			return a.teamID < b.teamID
		}
		return a.id < b.id
	})

	playerOutputs := make([]PlayerOutput, 0, len(sortedIDs))
	for _, pid := range sortedIDs {
		playerOutputs = append(playerOutputs, finalPlayers[pid].toOutput())
	}

	settings := buildSettings(meta.mapMeta, versionNum)

	parseTimeMS := time.Since(startTime).Milliseconds()

	return &ParserOutput{
		ID:            id,
		Gamename:      meta.gameName,
		RandomSeed:    int(meta.randomSeed),
		StartSpots:    int(meta.startSpotCount),
		Observers:     observers,
		Players:       playerOutputs,
		Matchup:       matchup,
		Creator:       meta.mapMeta.creator,
		Type:          gametype,
		Chat:          chatLog,
		APM:           APMConfig{TrackingInterval: playerActionTrackInterval},
		Map: MapInfo{
			Path:         meta.mapMeta.mapName,
			File:         mapFilename(meta.mapMeta.mapName),
			Checksum:     meta.mapMeta.mapChecksum,
			ChecksumSha1: meta.mapMeta.mapChecksumSha1,
		},
		BuildNumber:   int(subHeader.buildNo),
		Version:       gameVersion(int(subHeader.version)),
		Duration:      int(subHeader.replayLengthMS),
		Expansion:     subHeader.gameIdentifier == "PX3W",
		ParseTime:     parseTimeMS,
		WinningTeamID: winningTeamID,
		Settings:      settings,
	}
}

func processAction(action Action, p *player, totalTime int, slotToPlayerID map[int]uint8) {
	switch action.Type {
	case ActUnitAbilityNoParams:
		fmtID := FormatObjectID(action.OrderID)
		if fmtID.IsStringEncoded() {
			s := fmtID.StrVal
			if s == "tert" || s == "tret" {
				p.handleRetraining(totalTime)
			}
		}
		p.handle0x10(fmtID, totalTime)
	case ActUnitAbilityTargetPos:
		p.handle0x11(FormatObjectID(action.OrderID), totalTime)
	case ActUnitAbilityTargetObj:
		p.handle0x12(FormatObjectID(action.OrderID), totalTime)
	case ActGiveItemToUnit:
		p.handle0x13()
	case ActUnitAbilityTwoTargets:
		p.handle0x14(FormatObjectID(action.OrderID))
	case ActUnitAbilityTwoTargetsItem:
		p.handle0x14(FormatObjectID(action.OrderID))
	case ActChangeSelection:
		if action.SelectMode == 0x02 {
			p.lastActionWasDeselect = true
			p.handle0x16(action.SelectMode, true)
		} else {
			if !p.lastActionWasDeselect {
				p.handle0x16(action.SelectMode, true)
			}
			p.lastActionWasDeselect = false
		}
	case ActAssignGroupHotkey, ActSelectGroupHotkey, ActSelectGroundItem, ActCancelHeroRevival,
		ActRemoveUnitFromQueue, ActEscPressed, ActChooseHeroSkillSubmenu, ActEnterBuildingSubmenu:
		p.handleOther(action)
	case ActTransferResources:
		if pid, ok := slotToPlayerID[int(action.Slot)]; ok {
			p.handle0x51(action.Slot, pid, "", action.Gold, action.Lumber)
		}
	}
}

func isObserver(p *player, version uint32) bool {
	if version >= 29 {
		return p.teamID == 24
	}
	return p.teamID == 12
}

func determineMatchup(players map[uint8]*player) (string, string) {
	teamRaces := make(map[uint8][]string)
	for _, p := range players {
		race := p.raceDetected
		if race == "" {
			race = p.race
		}
		teamRaces[p.teamID] = append(teamRaces[p.teamID], race)
	}

	sizes := make([]int, 0, len(teamRaces))
	for _, races := range teamRaces {
		sizes = append(sizes, len(races))
	}
	sort.Ints(sizes)

	sizeStrs := make([]string, len(sizes))
	for i, s := range sizes {
		sizeStrs[i] = fmt.Sprintf("%d", s)
	}
	gametype := joinStrings(sizeStrs, "on")

	raceGroups := make([]string, 0, len(teamRaces))
	for _, races := range teamRaces {
		sorted := make([]string, len(races))
		copy(sorted, races)
		sort.Strings(sorted)
		raceGroups = append(raceGroups, joinStrings(sorted, ""))
	}
	sort.Strings(raceGroups)
	matchup := joinStrings(raceGroups, "v")

	return gametype, matchup
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for _, s := range ss[1:] {
		result += sep + s
	}
	return result
}

func determineWinningTeam(gametype string, leaveEvents []leaveGameBlock, players map[uint8]*player) int {
	if gametype != "1on1" {
		return -1
	}

	winningTeamID := -1
	for index, event := range leaveEvents {
		p, ok := players[event.playerID]
		if !ok {
			continue
		}
		if winningTeamID != -1 {
			continue
		}
		if event.result == "09000000" {
			winningTeamID = int(p.teamID)
			continue
		}
		if event.reason == "0c000000" {
			winningTeamID = int(p.teamID)
			continue
		}
		if index == len(leaveEvents)-1 {
			winningTeamID = int(p.teamID)
		}
	}
	return winningTeamID
}

func generateID(randomSeed uint32, players map[uint8]*player, gameName string) string {
	sortedIDs := make([]uint8, 0, len(players))
	for id := range players {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Slice(sortedIDs, func(i, j int) bool { return sortedIDs[i] < sortedIDs[j] })

	names := ""
	for _, id := range sortedIDs {
		names += players[id].name
	}
	idBase := fmt.Sprintf("%d%s%s", randomSeed, names, gameName)
	h := sha256.Sum256([]byte(idBase))
	return encodeHex(h[:])
}

func raceFlagToString(flag uint8) string {
	switch flag {
	case 0x01, 0x41:
		return "H"
	case 0x02, 0x42:
		return "O"
	case 0x04, 0x44:
		return "N"
	case 0x08, 0x48:
		return "U"
	default:
		return "R"
	}
}

func chatModeFromInt(mode uint32) ChatMode {
	switch mode {
	case 0x00:
		return ChatModeAll
	case 0x01:
		return ChatModeTeam
	case 0x02:
		return ChatModeObservers
	default:
		return ChatModePrivate
	}
}

func buildSettings(mm mapMetadata, version uint32) Settings {
	observerMode := getObserverMode(mm.referees, mm.observerMode)
	return Settings{
		ObserverMode:          observerMode,
		Referees:              mm.referees,
		FixedTeams:            mm.fixedTeams,
		FullSharedUnitControl: mm.fullSharedUnitControl,
		AlwaysVisible:         mm.alwaysVisible,
		HideTerrain:           mm.hideTerrain,
		MapExplored:           mm.mapExplored,
		TeamsTogether:         mm.teamsTogether,
		RandomHero:            mm.randomHero,
		RandomRaces:           mm.randomRaces,
		Speed:                 int(mm.speed),
	}
}

func getObserverMode(refereeFlag bool, observerMode uint8) ObserverMode {
	if (observerMode == 3 || observerMode == 0) && refereeFlag {
		return ObserverModeReferees
	} else if observerMode == 2 {
		return ObserverModeOnDefeat
	} else if observerMode == 3 {
		return ObserverModeFull
	}
	return ObserverModeNone
}
