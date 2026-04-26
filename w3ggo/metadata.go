package w3ggo

import "encoding/binary"

type playerRecord struct {
	playerID   uint8
	playerName string
}

type slotRecord struct {
	playerID         uint8
	downloadProgress uint8
	slotStatus       uint8
	computerFlag     uint8
	teamID           uint8
	color            uint8
	raceFlag         uint8
	aiStrength       uint8
	handicapFlag     uint8
}

type reforgedPlayerMeta struct {
	playerID uint32
	name     string
	clan     string
}

type mapMetadata struct {
	speed                 uint8
	hideTerrain           bool
	mapExplored           bool
	alwaysVisible         bool
	defaultFlag           bool
	observerMode          uint8
	teamsTogether         bool
	fixedTeams            bool
	fullSharedUnitControl bool
	randomHero            bool
	randomRaces           bool
	referees              bool
	mapChecksum           string
	mapChecksumSha1       string
	mapName               string
	creator               string
}

type metadata struct {
	gameData               []byte
	mapMeta                mapMetadata
	playerCount            uint32
	gameType               string
	localeHash             string
	playerRecords          []playerRecord
	slotRecords            []slotRecord
	reforgedPlayerMetadata []reforgedPlayerMeta
	randomSeed             uint32
	selectMode             string
	gameName               string
	startSpotCount         uint8
	isPost202              bool
}

func parseMetadata(data []byte) *metadata {
	pos := 0

	// skip 5 bytes
	pos += 5

	// Parse host record
	host, ok := parseHostRecord(data, &pos)
	if !ok {
		return nil
	}
	playerRecords := []playerRecord{host}

	// game name
	gameName := readZTSUtf8(data, &pos)

	// private string (discard)
	readZTSUtf8(data, &pos)

	// encoded string - read raw bytes
	encodedRaw := readZeroTermRaw(data, &pos)

	// Decode game meta string
	mapMetaDecoded := decodeGameMetaString(encodedRaw)
	mapMeta, ok := parseEncodedMapMetaString(mapMetaDecoded)
	if !ok {
		return nil
	}

	// player count
	playerCount, ok := readU32LE(data, &pos)
	if !ok {
		return nil
	}

	// game type (4 bytes as hex)
	gameType := readHex(data, &pos, 4)

	// locale hash (4 bytes as hex)
	localeHash := readHex(data, &pos, 4)

	// Additional player list
	additional := parsePlayerList(data, &pos)

	// JS: playerRecordsFinal = playerRecords.concat(playerRecords, parsePlayerList())
	var playerRecordsFinal []playerRecord
	playerRecordsFinal = append(playerRecordsFinal, playerRecords...)
	playerRecordsFinal = append(playerRecordsFinal, playerRecords...)
	playerRecordsFinal = append(playerRecordsFinal, additional...)

	var reforgedPlayerMetadata []reforgedPlayerMeta
	isPost202 := false

	// Reforged metadata blocks (0x38 or 0x39)
	for pos < len(data) && (data[pos] == 0x38 || data[pos] == 0x39) {
		recordType := data[pos]
		pos++
		if recordType == 0x38 {
			isPost202 = true
		}
		if pos >= len(data) {
			break
		}
		subtype := data[pos]
		pos++
		if pos+4 > len(data) {
			break
		}
		followingBytes := int(binary.LittleEndian.Uint32(data[pos:]))
		pos += 4
		blobEnd := pos + followingBytes
		if blobEnd > len(data) {
			break
		}
		blob := data[pos:blobEnd]
		if subtype == 0x03 {
			players := decodePlayerMetadata(blob)
			for _, pm := range players {
				reforgedPlayerMetadata = append(reforgedPlayerMetadata, reforgedPlayerMeta{
					playerID: pm.id,
					name:     pm.battleTag,
					clan:     pm.clan,
				})
			}
		}
		pos += followingBytes
	}

	// expect 0x19
	pos++ // _check
	if pos+2 > len(data) {
		return nil
	}
	pos += 2 // _remaining_bytes
	if pos >= len(data) {
		return nil
	}
	slotRecordCount := int(data[pos])
	pos++

	slotRecords := parseSlotRecords(data, &pos, slotRecordCount)

	randomSeed, ok := readU32LE(data, &pos)
	if !ok {
		return nil
	}
	selectMode := readHex(data, &pos, 1)
	if pos >= len(data) {
		return nil
	}
	startSpotCount := data[pos]
	pos++

	gameData := make([]byte, len(data)-pos)
	copy(gameData, data[pos:])

	return &metadata{
		gameData:               gameData,
		mapMeta:                mapMeta,
		playerCount:            playerCount,
		gameType:               gameType,
		localeHash:             localeHash,
		playerRecords:          playerRecordsFinal,
		slotRecords:            slotRecords,
		reforgedPlayerMetadata: reforgedPlayerMetadata,
		randomSeed:             randomSeed,
		selectMode:             selectMode,
		gameName:               gameName,
		startSpotCount:         startSpotCount,
		isPost202:              isPost202,
	}
}

func parseHostRecord(data []byte, pos *int) (playerRecord, bool) {
	if *pos >= len(data) {
		return playerRecord{}, false
	}
	playerID := data[*pos]
	*pos++
	playerName := readZTSUtf8(data, pos)
	if *pos >= len(data) {
		return playerRecord{}, false
	}
	addData := int(data[*pos])
	*pos++
	*pos += addData
	return playerRecord{playerID: playerID, playerName: playerName}, true
}

func parsePlayerList(data []byte, pos *int) []playerRecord {
	var list []playerRecord
	for *pos < len(data) && data[*pos] == 22 {
		*pos++ // consume the 22
		record, ok := parseHostRecord(data, pos)
		if ok {
			list = append(list, record)
		}
		*pos += 4
	}
	return list
}

func parseSlotRecords(data []byte, pos *int, count int) []slotRecord {
	slots := make([]slotRecord, 0, count)
	for i := 0; i < count; i++ {
		s := slotRecord{}
		if *pos < len(data) { s.playerID = data[*pos]; *pos++ }
		if *pos < len(data) { s.downloadProgress = data[*pos]; *pos++ }
		if *pos < len(data) { s.slotStatus = data[*pos]; *pos++ }
		if *pos < len(data) { s.computerFlag = data[*pos]; *pos++ }
		if *pos < len(data) { s.teamID = data[*pos]; *pos++ }
		if *pos < len(data) { s.color = data[*pos]; *pos++ }
		if *pos < len(data) { s.raceFlag = data[*pos]; *pos++ }
		if *pos < len(data) { s.aiStrength = data[*pos]; *pos++ }
		if *pos < len(data) { s.handicapFlag = data[*pos]; *pos++ }
		slots = append(slots, s)
	}
	return slots
}

func decodeGameMetaString(b []byte) []byte {
	decoded := make([]byte, 0, len(b))
	var mask uint8
	for i, bv := range b {
		if i%8 == 0 {
			mask = bv
		} else {
			bitPos := uint(i % 8)
			if (mask & (1 << bitPos)) == 0 {
				decoded = append(decoded, bv-1)
			} else {
				decoded = append(decoded, bv)
			}
		}
	}
	return decoded
}

func parseEncodedMapMetaString(data []byte) (mapMetadata, bool) {
	pos := 0
	if pos >= len(data) {
		return mapMetadata{}, false
	}
	speed := data[pos]; pos++
	if pos >= len(data) {
		return mapMetadata{}, false
	}
	secondByte := data[pos]; pos++
	if pos >= len(data) {
		return mapMetadata{}, false
	}
	thirdByte := data[pos]; pos++
	if pos >= len(data) {
		return mapMetadata{}, false
	}
	fourthByte := data[pos]; pos++
	pos += 5 // skip 5
	checksum := readHex(data, &pos, 4)
	mapName := readZTSUtf8(data, &pos)
	creator := readZTSUtf8(data, &pos)
	pos++ // skip 1
	checksumSha1 := readHex(data, &pos, 20)

	return mapMetadata{
		speed:                 speed,
		hideTerrain:           (secondByte & 0b00000001) != 0,
		mapExplored:           (secondByte & 0b00000010) != 0,
		alwaysVisible:         (secondByte & 0b00000100) != 0,
		defaultFlag:           (secondByte & 0b00001000) != 0,
		observerMode:          (secondByte & 0b00110000) >> 4,
		teamsTogether:         (secondByte & 0b01000000) != 0,
		fixedTeams:            (thirdByte & 0b00000110) != 0,
		fullSharedUnitControl: (fourthByte & 0b00000001) != 0,
		randomHero:            (fourthByte & 0b00000010) != 0,
		randomRaces:           (fourthByte & 0b00000100) != 0,
		referees:              (fourthByte & 0b01000000) != 0,
		mapChecksum:           checksum,
		mapChecksumSha1:       checksumSha1,
		mapName:               mapName,
		creator:               creator,
	}, true
}

// --- helpers ---

func readZeroTermRaw(data []byte, pos *int) []byte {
	start := *pos
	for *pos < len(data) && data[*pos] != 0 {
		*pos++
	}
	result := make([]byte, *pos-start)
	copy(result, data[start:*pos])
	if *pos < len(data) {
		*pos++ // consume null
	}
	return result
}

func readZTSUtf8(data []byte, pos *int) string {
	start := *pos
	for *pos < len(data) && data[*pos] != 0 {
		*pos++
	}
	s := string(data[start:*pos])
	if *pos < len(data) {
		*pos++
	}
	return s
}

func readU32LE(data []byte, pos *int) (uint32, bool) {
	if *pos+4 > len(data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(data[*pos:])
	*pos += 4
	return v, true
}

func readHex(data []byte, pos *int, n int) string {
	if *pos+n > len(data) {
		return ""
	}
	s := encodeHex(data[*pos : *pos+n])
	*pos += n
	return s
}
