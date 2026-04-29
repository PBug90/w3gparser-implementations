package w3ggo

import "encoding/binary"

type leaveGameBlock struct {
	playerID uint8
	reason   string
	result   string
}

type commandBlock struct {
	playerID uint8
	actions  []Action
}

type timeslotBlock struct {
	timeIncrement uint16
	commandBlocks []commandBlock
}

type chatMessageBlock struct {
	playerID uint8
	mode     uint32
	message  string
}

type gameDataBlockType int

const (
	gdLeaveGame gameDataBlockType = iota
	gdTimeslot
	gdChatMessage
)

type gameDataBlock struct {
	typ       gameDataBlockType
	leaveGame leaveGameBlock
	timeslot  timeslotBlock
	chat      chatMessageBlock
}

func parseGameData(data []byte, isPost202 bool) []gameDataBlock {
	var blocks []gameDataBlock
	pos := 0

	for pos < len(data) {
		if pos >= len(data) {
			break
		}
		id := data[pos]
		pos++

		switch id {
		case 0x17:
			b, ok := parseLeaveGame(data, &pos)
			if ok {
				blocks = append(blocks, gameDataBlock{typ: gdLeaveGame, leaveGame: b})
			}
		case 0x1a, 0x1b, 0x1c:
			if pos+4 > len(data) {
				pos = len(data)
			} else {
				pos += 4
			}
		case 0x1e, 0x1f:
			b, ok := parseTimeslot(data, &pos, isPost202)
			if ok {
				blocks = append(blocks, gameDataBlock{typ: gdTimeslot, timeslot: b})
			}
		case 0x20:
			b, ok := parseChatMessage(data, &pos)
			if ok {
				blocks = append(blocks, gameDataBlock{typ: gdChatMessage, chat: b})
			}
		case 0x22:
			if pos >= len(data) {
				break
			}
			ln := int(data[pos])
			pos++
			if pos+ln > len(data) {
				pos = len(data)
			} else {
				pos += ln
			}
		case 0x23:
			if pos+10 > len(data) {
				pos = len(data)
			} else {
				pos += 10
			}
		case 0x2f:
			if pos+8 > len(data) {
				pos = len(data)
			} else {
				pos += 8
			}
		default:
			// unknown block id - skip and continue (mirrors w3gjs behavior)
			continue
		}
	}
	return blocks
}

func parseLeaveGame(data []byte, pos *int) (leaveGameBlock, bool) {
	reason := gdReadHex(data, pos, 4)
	if *pos >= len(data) {
		return leaveGameBlock{}, false
	}
	playerID := data[*pos]
	*pos++
	result := gdReadHex(data, pos, 4)
	*pos += 4
	return leaveGameBlock{playerID: playerID, reason: reason, result: result}, true
}

func parseTimeslot(data []byte, pos *int, isPost202 bool) (timeslotBlock, bool) {
	if *pos+2 > len(data) {
		return timeslotBlock{}, false
	}
	byteCount := int(binary.LittleEndian.Uint16(data[*pos:]))
	*pos += 2
	if *pos+2 > len(data) {
		return timeslotBlock{}, false
	}
	timeIncrement := binary.LittleEndian.Uint16(data[*pos:])
	*pos += 2

	actionBlockLastOffset := *pos + byteCount - 2
	if actionBlockLastOffset > len(data) {
		actionBlockLastOffset = len(data)
	}

	var commandBlocks []commandBlock
	for *pos < actionBlockLastOffset && *pos < len(data) {
		if *pos >= len(data) {
			break
		}
		playerID := data[*pos]
		*pos++
		if *pos+2 > len(data) {
			break
		}
		actionBlockLength := int(binary.LittleEndian.Uint16(data[*pos:]))
		*pos += 2
		end := *pos + actionBlockLength
		if end > len(data) {
			end = len(data)
		}
		actionData := data[*pos:end]
		actions := parseActions(actionData, isPost202)
		*pos = end
		commandBlocks = append(commandBlocks, commandBlock{playerID: playerID, actions: actions})
	}

	return timeslotBlock{timeIncrement: timeIncrement, commandBlocks: commandBlocks}, true
}

func parseChatMessage(data []byte, pos *int) (chatMessageBlock, bool) {
	if *pos >= len(data) {
		return chatMessageBlock{}, false
	}
	playerID := data[*pos]
	*pos++
	*pos += 2 // byteCount
	if *pos >= len(data) {
		return chatMessageBlock{}, false
	}
	flags := data[*pos]
	*pos++
	var mode uint32
	if flags == 0x20 {
		if *pos+4 > len(data) {
			return chatMessageBlock{}, false
		}
		mode = binary.LittleEndian.Uint32(data[*pos:])
		*pos += 4
	}
	message := gdReadZTS(data, pos)
	return chatMessageBlock{playerID: playerID, mode: mode, message: message}, true
}

func gdReadHex(data []byte, pos *int, n int) string {
	if *pos+n > len(data) {
		return ""
	}
	s := encodeHex(data[*pos : *pos+n])
	*pos += n
	return s
}

func gdReadZTS(data []byte, pos *int) string {
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
