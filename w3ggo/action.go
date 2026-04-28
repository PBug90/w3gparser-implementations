package w3ggo

import (
	"encoding/binary"
	"math"
)

// ObjectID represents a formatted 4-byte object id.
type ObjectID struct {
	IsString bool
	StrVal   string
	ArrVal   [4]byte
}

// IsStringEncoded reports whether this ObjectID is a human-readable string.
func (o *ObjectID) IsStringEncoded() bool { return o.IsString }

// FirstChar returns the first character of a string-encoded ObjectID.
func (o *ObjectID) FirstChar() (byte, bool) {
	if o.IsString && len(o.StrVal) > 0 {
		return o.StrVal[0], true
	}
	return 0, false
}

// FormatObjectID converts a raw 4-byte order id to an ObjectID, mirroring
// objectIdFormatter in w3gjs.
func FormatObjectID(arr [4]byte) ObjectID {
	if arr[3] >= 0x41 && arr[3] <= 0x7a {
		s := string([]byte{arr[3], arr[2], arr[1], arr[0]})
		return ObjectID{IsString: true, StrVal: s}
	}
	return ObjectID{IsString: false, ArrVal: arr}
}

// ActionType discriminates Action variants.
type ActionType int

const (
	ActUnitAbilityNoParams ActionType = iota
	ActUnitAbilityTargetPos
	ActUnitAbilityTargetObj
	ActGiveItemToUnit
	ActUnitAbilityTwoTargets
	ActUnitAbilityTwoTargetsItem
	ActChangeSelection
	ActAssignGroupHotkey
	ActSelectGroupHotkey
	ActSelectSubgroup
	ActPreSubselection
	ActSelectUnit
	ActSelectGroundItem
	ActCancelHeroRevival
	ActRemoveUnitFromQueue
	ActTransferResources
	ActEscPressed
	ActChooseHeroSkillSubmenu
	ActEnterBuildingSubmenu
	ActW3MMDStoreInt
	ActW3MMDStoreReal
	ActW3MMDStoreBool
	ActW3MMDClearInt
	ActW3MMDClearReal
	ActW3MMDClearBool
	ActW3MMDClearUnit
	ActAllyPing
	ActArrowKey
	ActSetGameSpeed
	ActTrackableHit
	ActBlzSync
	ActCommandFrame
	ActMouseAction
	ActW3Api
)

// CacheData holds W3MMD cache key information used in W3MMD* actions.
type CacheData struct {
	Filename   string
	MissionKey string
	Key        string
}

// Action represents a single player action inside a command block.
// The Type field discriminates which variant this is; the remaining fields
// carry variant-specific data.
type Action struct {
	Type         ActionType
	AbilityFlags uint16
	OrderID      [4]byte
	OrderID2     [4]byte
	SelectMode   uint8
	GroupNumber  uint8
	NumberUnits  uint16
	SlotNumber   uint8
	Slot         uint8
	Gold         uint32
	Lumber       uint32
	ArrowKey     uint8
	GameSpeed    uint8
	EventID      uint32
	EventIDB     uint8
	Button       uint8
	ValF         float32
	Text         string
	Identifier   string
	Value        string
	CmdID        uint32
	CmdData      uint32
	Cache        CacheData
	W3MMDValue   uint32
	W3MMDReal    float32
	W3MMDBool    uint8
	Object       [2]uint32
	Item         [2]uint32
	Unit         [2]uint32
	Hero         [2]uint32
	TargetPos    [2]float32
	TargetB      [2]float32
	Flags        uint32
	Category     uint32
	Owner        uint8
	ItemID4      [4]byte
	Duration     float32
	PingPos      [2]float32
}

func parseActions(data []byte, isPost202 bool) []Action {
	var actions []Action
	pos := 0

	for pos < len(data) {
		if pos >= len(data) {
			break
		}
		actionIDRaw := data[pos]
		pos++

		actionID := actionIDRaw
		if isPost202 && actionIDRaw > 0x77 {
			actionID = actionIDRaw + 1
		}

		a, ok := parseSingleAction(actionID, data, &pos)
		if ok {
			actions = append(actions, a)
		}
	}
	return actions
}

func parseSingleAction(id uint8, data []byte, pos *int) (Action, bool) {
	switch id {
	case 0x01:
		advance(data, pos, 1)
		return Action{}, false
	case 0x02, 0x04, 0x05:
		return Action{}, false
	case 0x03:
		gs, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActSetGameSpeed, GameSpeed: gs}, true
	case 0x06:
		rzts(data, pos)
		rzts(data, pos)
		advance(data, pos, 1)
		return Action{}, false
	case 0x07:
		advance(data, pos, 4)
		return Action{}, false
	case 0x10:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		return Action{Type: ActUnitAbilityNoParams, AbilityFlags: af, OrderID: oid}, true
	case 0x11:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		t, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActUnitAbilityTargetPos, AbilityFlags: af, OrderID: oid, TargetPos: t}, true
	case 0x12:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		t, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		obj, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActUnitAbilityTargetObj, AbilityFlags: af, OrderID: oid, TargetPos: t, Object: obj}, true
	case 0x13:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		t, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		unitTag, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		itemTag, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActGiveItemToUnit, AbilityFlags: af, OrderID: oid, TargetPos: t, Unit: unitTag, Item: itemTag}, true
	case 0x14:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid1, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		ta, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		oid2, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		fl, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		cat, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		own, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		tb, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActUnitAbilityTwoTargets, AbilityFlags: af, OrderID: oid1, OrderID2: oid2, TargetPos: ta, TargetB: tb, Flags: fl, Category: cat, Owner: own}, true
	case 0x15:
		af, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		oid1, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 8)
		ta, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		oid2, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		fl, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		cat, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		own, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		tb, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		obj, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActUnitAbilityTwoTargetsItem, AbilityFlags: af, OrderID: oid1, OrderID2: oid2, TargetPos: ta, TargetB: tb, Flags: fl, Category: cat, Owner: own, Object: obj}, true
	case 0x16:
		sm, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		nu, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, int(nu)*8)
		return Action{Type: ActChangeSelection, SelectMode: sm, NumberUnits: nu}, true
	case 0x17:
		gn, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		nu, ok := ru16(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, int(nu)*8)
		return Action{Type: ActAssignGroupHotkey, GroupNumber: gn, NumberUnits: nu}, true
	case 0x18:
		gn, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, 1)
		return Action{Type: ActSelectGroupHotkey, GroupNumber: gn}, true
	case 0x19:
		iid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		obj, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActSelectSubgroup, ItemID4: iid, Object: obj}, true
	case 0x1a:
		return Action{Type: ActPreSubselection}, true
	case 0x1b:
		advance(data, pos, 1)
		obj, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActSelectUnit, Object: obj}, true
	case 0x1c:
		advance(data, pos, 1)
		itm, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActSelectGroundItem, Item: itm}, true
	case 0x1d:
		hr, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActCancelHeroRevival, Hero: hr}, true
	case 0x1e, 0x1f:
		sn, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		iid, ok := rfourcc(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActRemoveUnitFromQueue, SlotNumber: sn, ItemID4: iid}, true
	case 0x20:
		return Action{}, false
	case 0x21:
		advance(data, pos, 8)
		return Action{}, false
	case 0x22, 0x23, 0x24, 0x25, 0x26:
		return Action{}, false
	case 0x27, 0x28:
		advance(data, pos, 5)
		return Action{}, false
	case 0x29, 0x2a, 0x2b, 0x2c:
		return Action{}, false
	case 0x2d:
		advance(data, pos, 5)
		return Action{}, false
	case 0x2e:
		advance(data, pos, 4)
		return Action{}, false
	case 0x2f:
		return Action{}, false
	case 0x50:
		advance(data, pos, 1)
		advance(data, pos, 4)
		return Action{}, false
	case 0x51:
		sl, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		g, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		l, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActTransferResources, Slot: sl, Gold: g, Lumber: l}, true
	case 0x60:
		advance(data, pos, 8)
		rzts(data, pos)
		return Action{}, false
	case 0x61:
		return Action{Type: ActEscPressed}, true
	case 0x62:
		advance(data, pos, 12)
		return Action{}, false
	case 0x63:
		advance(data, pos, 8)
		return Action{}, false
	case 0x64, 0x65:
		obj, ok := rnettag(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActTrackableHit, Object: obj}, true
	case 0x66:
		return Action{Type: ActChooseHeroSkillSubmenu}, true
	case 0x67:
		return Action{Type: ActEnterBuildingSubmenu}, true
	case 0x68:
		pp, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		dur, ok := rf32(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActAllyPing, PingPos: pp, Duration: dur}, true
	case 0x69, 0x6a:
		advance(data, pos, 16)
		return Action{}, false
	case 0x6b:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		v, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDStoreInt, Cache: c, W3MMDValue: v}, true
	case 0x6c:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		v, ok := rf32(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDStoreReal, Cache: c, W3MMDReal: v}, true
	case 0x6d:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		v, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDStoreBool, Cache: c, W3MMDBool: v}, true
	case 0x6e:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		readCacheUnit(data, pos)
		return Action{Type: ActW3MMDClearUnit, Cache: c}, true
	case 0x70:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDClearInt, Cache: c}, true
	case 0x71:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDClearReal, Cache: c}, true
	case 0x72:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDClearBool, Cache: c}, true
	case 0x73:
		c, ok := rcache(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActW3MMDClearUnit, Cache: c}, true
	case 0x75:
		ak, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActArrowKey, ArrowKey: ak}, true
	case 0x76:
		eid, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		pp, ok := rvec2(data, pos)
		if !ok {
			return Action{}, false
		}
		btn, ok := ru8(data, pos)
		if !ok {
			return Action{}, false
		}
		return Action{Type: ActMouseAction, EventIDB: eid, PingPos: pp, Button: btn}, true
	case 0x77:
		cid, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		dv, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		bl, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		advance(data, pos, int(bl))
		return Action{Type: ActW3Api, CmdID: cid, CmdData: dv}, true
	case 0x78:
		ident := rzts(data, pos)
		val := rzts(data, pos)
		advance(data, pos, 4)
		return Action{Type: ActBlzSync, Identifier: ident, Value: val}, true
	case 0x79:
		advance(data, pos, 8)
		eid, ok := ru32(data, pos)
		if !ok {
			return Action{}, false
		}
		v, ok := rf32(data, pos)
		if !ok {
			return Action{}, false
		}
		txt := rzts(data, pos)
		return Action{Type: ActCommandFrame, EventID: eid, ValF: v, Text: txt}, true
	case 0x7a:
		advance(data, pos, 20)
		return Action{}, false
	case 0x7b:
		advance(data, pos, 16)
		return Action{}, false
	case 0xa0:
		advance(data, pos, 14)
		return Action{}, false
	case 0xa1:
		advance(data, pos, 9)
		return Action{}, false
	default:
		return Action{}, false
	}
}

func readCacheUnit(data []byte, pos *int) {
	advance(data, pos, 4) // unitId
	if *pos+4 > len(data) {
		return
	}
	itemsCount := int(binary.LittleEndian.Uint32(data[*pos:]))
	*pos += 4
	advance(data, pos, itemsCount*12)
	advance(data, pos, 4*4+2*4)
	advance(data, pos, 4+4+4+4)
	advance(data, pos, 4+4)
	if *pos+4 > len(data) {
		return
	}
	heroAbilCount := int(binary.LittleEndian.Uint32(data[*pos:]))
	*pos += 4
	advance(data, pos, heroAbilCount*8)
	advance(data, pos, 12)
	if *pos+4 > len(data) {
		return
	}
	damageCount := int(binary.LittleEndian.Uint32(data[*pos:]))
	*pos += 4
	advance(data, pos, damageCount*4)
	advance(data, pos, 4+2)
}

// --- tiny helpers ---

func advance(data []byte, pos *int, n int) {
	newPos := *pos + n
	if newPos > len(data) {
		newPos = len(data)
	}
	*pos = newPos
}

func ru8(data []byte, pos *int) (uint8, bool) {
	if *pos >= len(data) {
		return 0, false
	}
	v := data[*pos]
	*pos++
	return v, true
}

func ru16(data []byte, pos *int) (uint16, bool) {
	if *pos+2 > len(data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint16(data[*pos:])
	*pos += 2
	return v, true
}

func ru32(data []byte, pos *int) (uint32, bool) {
	if *pos+4 > len(data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(data[*pos:])
	*pos += 4
	return v, true
}

func rf32(data []byte, pos *int) (float32, bool) {
	if *pos+4 > len(data) {
		return 0, false
	}
	bits := binary.LittleEndian.Uint32(data[*pos:])
	*pos += 4
	return math.Float32frombits(bits), true
}

func rfourcc(data []byte, pos *int) ([4]byte, bool) {
	if *pos+4 > len(data) {
		return [4]byte{}, false
	}
	arr := [4]byte{data[*pos], data[*pos+1], data[*pos+2], data[*pos+3]}
	*pos += 4
	return arr, true
}

func rnettag(data []byte, pos *int) ([2]uint32, bool) {
	a, ok := ru32(data, pos)
	if !ok {
		return [2]uint32{}, false
	}
	b, ok := ru32(data, pos)
	if !ok {
		return [2]uint32{}, false
	}
	return [2]uint32{a, b}, true
}

func rvec2(data []byte, pos *int) ([2]float32, bool) {
	x, ok := rf32(data, pos)
	if !ok {
		return [2]float32{}, false
	}
	y, ok := rf32(data, pos)
	if !ok {
		return [2]float32{}, false
	}
	return [2]float32{x, y}, true
}

func rzts(data []byte, pos *int) string {
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

func rcache(data []byte, pos *int) (CacheData, bool) {
	fn := rzts(data, pos)
	mk := rzts(data, pos)
	k := rzts(data, pos)
	return CacheData{Filename: fn, MissionKey: mk, Key: k}, true
}
