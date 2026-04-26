package w3ggo

import (
	"encoding/binary"
	"math"
)

// ObjectId represents a formatted 4-byte object id.
type ObjectId struct {
	isString bool
	strVal   string
	arrVal   [4]byte
}

func (o *ObjectId) isStringEncoded() bool { return o.isString }

func (o *ObjectId) firstChar() (byte, bool) {
	if o.isString && len(o.strVal) > 0 {
		return o.strVal[0], true
	}
	return 0, false
}

// formatObjectId mirrors objectIdFormatter in JS/Rust.
func formatObjectId(arr [4]byte) ObjectId {
	if arr[3] >= 0x41 && arr[3] <= 0x7a {
		// string encoded: reverse array, convert to chars
		s := string([]byte{arr[3], arr[2], arr[1], arr[0]})
		return ObjectId{isString: true, strVal: s}
	}
	return ObjectId{isString: false, arrVal: arr}
}

// Action types
type actionType int

const (
	actUnitAbilityNoParams actionType = iota
	actUnitAbilityTargetPos
	actUnitAbilityTargetObj
	actGiveItemToUnit
	actUnitAbilityTwoTargets
	actUnitAbilityTwoTargetsItem
	actChangeSelection
	actAssignGroupHotkey
	actSelectGroupHotkey
	actSelectSubgroup
	actPreSubselection
	actSelectUnit
	actSelectGroundItem
	actCancelHeroRevival
	actRemoveUnitFromQueue
	actTransferResources
	actEscPressed
	actChooseHeroSkillSubmenu
	actEnterBuildingSubmenu
	actW3MMDStoreInt
	actW3MMDStoreReal
	actW3MMDStoreBool
	actW3MMDClearInt
	actW3MMDClearReal
	actW3MMDClearBool
	actW3MMDClearUnit
	actAllyPing
	actArrowKey
	actSetGameSpeed
	actTrackableHit
	actBlzSync
	actCommandFrame
	actMouseAction
	actW3Api
)

type cacheData struct {
	filename   string
	missionKey string
	key        string
}

type Action struct {
	typ         actionType
	abilityFlags uint16
	orderId     [4]byte
	orderId2    [4]byte
	selectMode  uint8
	groupNumber uint8
	numberUnits uint16
	slotNumber  uint8
	slot        uint8
	gold        uint32
	lumber      uint32
	arrowKey    uint8
	gameSpeed   uint8
	eventID     uint32
	eventIDB    uint8
	button      uint8
	valF        float32
	text        string
	identifier  string
	value       string
	cmdID       uint32
	cmdData     uint32
	cache       cacheData
	w3mmdValue  uint32
	w3mmdReal   float32
	w3mmdBool   uint8
	object      [2]uint32
	item        [2]uint32
	unit        [2]uint32
	hero        [2]uint32
	targetPos   [2]float32
	targetB     [2]float32
	flags       uint32
	category    uint32
	owner       uint8
	itemId4     [4]byte
	duration    float32
	pingPos     [2]float32
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
		if !ok { return Action{}, false }
		return Action{typ: actSetGameSpeed, gameSpeed: gs}, true
	case 0x06:
		rzts(data, pos); rzts(data, pos); advance(data, pos, 1)
		return Action{}, false
	case 0x07:
		advance(data, pos, 4)
		return Action{}, false
	case 0x10:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		return Action{typ: actUnitAbilityNoParams, abilityFlags: af, orderId: oid}, true
	case 0x11:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		t, ok := rvec2(data, pos); if !ok { return Action{}, false }
		return Action{typ: actUnitAbilityTargetPos, abilityFlags: af, orderId: oid, targetPos: t}, true
	case 0x12:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		t, ok := rvec2(data, pos); if !ok { return Action{}, false }
		obj, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actUnitAbilityTargetObj, abilityFlags: af, orderId: oid, targetPos: t, object: obj}, true
	case 0x13:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		t, ok := rvec2(data, pos); if !ok { return Action{}, false }
		unitTag, ok := rnettag(data, pos); if !ok { return Action{}, false }
		itemTag, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actGiveItemToUnit, abilityFlags: af, orderId: oid, targetPos: t, unit: unitTag, item: itemTag}, true
	case 0x14:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid1, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		ta, ok := rvec2(data, pos); if !ok { return Action{}, false }
		oid2, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		fl, ok := ru32(data, pos); if !ok { return Action{}, false }
		cat, ok := ru32(data, pos); if !ok { return Action{}, false }
		own, ok := ru8(data, pos); if !ok { return Action{}, false }
		tb, ok := rvec2(data, pos); if !ok { return Action{}, false }
		return Action{typ: actUnitAbilityTwoTargets, abilityFlags: af, orderId: oid1, orderId2: oid2, targetPos: ta, targetB: tb, flags: fl, category: cat, owner: own}, true
	case 0x15:
		af, ok := ru16(data, pos); if !ok { return Action{}, false }
		oid1, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 8)
		ta, ok := rvec2(data, pos); if !ok { return Action{}, false }
		oid2, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		fl, ok := ru32(data, pos); if !ok { return Action{}, false }
		cat, ok := ru32(data, pos); if !ok { return Action{}, false }
		own, ok := ru8(data, pos); if !ok { return Action{}, false }
		tb, ok := rvec2(data, pos); if !ok { return Action{}, false }
		obj, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actUnitAbilityTwoTargetsItem, abilityFlags: af, orderId: oid1, orderId2: oid2, targetPos: ta, targetB: tb, flags: fl, category: cat, owner: own, object: obj}, true
	case 0x16:
		sm, ok := ru8(data, pos); if !ok { return Action{}, false }
		nu, ok := ru16(data, pos); if !ok { return Action{}, false }
		advance(data, pos, int(nu)*8)
		return Action{typ: actChangeSelection, selectMode: sm, numberUnits: nu}, true
	case 0x17:
		gn, ok := ru8(data, pos); if !ok { return Action{}, false }
		nu, ok := ru16(data, pos); if !ok { return Action{}, false }
		advance(data, pos, int(nu)*8)
		return Action{typ: actAssignGroupHotkey, groupNumber: gn, numberUnits: nu}, true
	case 0x18:
		gn, ok := ru8(data, pos); if !ok { return Action{}, false }
		advance(data, pos, 1)
		return Action{typ: actSelectGroupHotkey, groupNumber: gn}, true
	case 0x19:
		iid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		obj, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actSelectSubgroup, itemId4: iid, object: obj}, true
	case 0x1a:
		return Action{typ: actPreSubselection}, true
	case 0x1b:
		advance(data, pos, 1)
		obj, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actSelectUnit, object: obj}, true
	case 0x1c:
		advance(data, pos, 1)
		itm, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actSelectGroundItem, item: itm}, true
	case 0x1d:
		hr, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actCancelHeroRevival, hero: hr}, true
	case 0x1e, 0x1f:
		sn, ok := ru8(data, pos); if !ok { return Action{}, false }
		iid, ok := rfourcc(data, pos); if !ok { return Action{}, false }
		return Action{typ: actRemoveUnitFromQueue, slotNumber: sn, itemId4: iid}, true
	case 0x20:
		return Action{}, false
	case 0x21:
		advance(data, pos, 8); return Action{}, false
	case 0x22, 0x23, 0x24, 0x25, 0x26:
		return Action{}, false
	case 0x27, 0x28:
		advance(data, pos, 5); return Action{}, false
	case 0x29, 0x2a, 0x2b, 0x2c:
		return Action{}, false
	case 0x2d:
		advance(data, pos, 5); return Action{}, false
	case 0x2e:
		advance(data, pos, 4); return Action{}, false
	case 0x2f:
		return Action{}, false
	case 0x50:
		advance(data, pos, 1); advance(data, pos, 4); return Action{}, false
	case 0x51:
		sl, ok := ru8(data, pos); if !ok { return Action{}, false }
		g, ok := ru32(data, pos); if !ok { return Action{}, false }
		l, ok := ru32(data, pos); if !ok { return Action{}, false }
		return Action{typ: actTransferResources, slot: sl, gold: g, lumber: l}, true
	case 0x60:
		advance(data, pos, 8); rzts(data, pos); return Action{}, false
	case 0x61:
		return Action{typ: actEscPressed}, true
	case 0x62:
		advance(data, pos, 12); return Action{}, false
	case 0x63:
		advance(data, pos, 8); return Action{}, false
	case 0x64, 0x65:
		obj, ok := rnettag(data, pos); if !ok { return Action{}, false }
		return Action{typ: actTrackableHit, object: obj}, true
	case 0x66:
		return Action{typ: actChooseHeroSkillSubmenu}, true
	case 0x67:
		return Action{typ: actEnterBuildingSubmenu}, true
	case 0x68:
		pp, ok := rvec2(data, pos); if !ok { return Action{}, false }
		dur, ok := rf32(data, pos); if !ok { return Action{}, false }
		return Action{typ: actAllyPing, pingPos: pp, duration: dur}, true
	case 0x69, 0x6a:
		advance(data, pos, 16); return Action{}, false
	case 0x6b:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		v, ok := ru32(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDStoreInt, cache: c, w3mmdValue: v}, true
	case 0x6c:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		v, ok := rf32(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDStoreReal, cache: c, w3mmdReal: v}, true
	case 0x6d:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		v, ok := ru8(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDStoreBool, cache: c, w3mmdBool: v}, true
	case 0x6e:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		readCacheUnit(data, pos)
		return Action{typ: actW3MMDClearUnit, cache: c}, true
	case 0x70:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDClearInt, cache: c}, true
	case 0x71:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDClearReal, cache: c}, true
	case 0x72:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDClearBool, cache: c}, true
	case 0x73:
		c, ok := rcache(data, pos); if !ok { return Action{}, false }
		return Action{typ: actW3MMDClearUnit, cache: c}, true
	case 0x75:
		ak, ok := ru8(data, pos); if !ok { return Action{}, false }
		return Action{typ: actArrowKey, arrowKey: ak}, true
	case 0x76:
		eid, ok := ru8(data, pos); if !ok { return Action{}, false }
		pp, ok := rvec2(data, pos); if !ok { return Action{}, false }
		btn, ok := ru8(data, pos); if !ok { return Action{}, false }
		return Action{typ: actMouseAction, eventIDB: eid, pingPos: pp, button: btn}, true
	case 0x77:
		cid, ok := ru32(data, pos); if !ok { return Action{}, false }
		dv, ok := ru32(data, pos); if !ok { return Action{}, false }
		bl, ok := ru32(data, pos); if !ok { return Action{}, false }
		advance(data, pos, int(bl))
		return Action{typ: actW3Api, cmdID: cid, cmdData: dv}, true
	case 0x78:
		ident := rzts(data, pos)
		val := rzts(data, pos)
		advance(data, pos, 4)
		return Action{typ: actBlzSync, identifier: ident, value: val}, true
	case 0x79:
		advance(data, pos, 8)
		eid, ok := ru32(data, pos); if !ok { return Action{}, false }
		v, ok := rf32(data, pos); if !ok { return Action{}, false }
		txt := rzts(data, pos)
		return Action{typ: actCommandFrame, eventID: eid, valF: v, text: txt}, true
	case 0x7a:
		advance(data, pos, 20); return Action{}, false
	case 0x7b:
		advance(data, pos, 16); return Action{}, false
	case 0xa0:
		advance(data, pos, 14); return Action{}, false
	case 0xa1:
		advance(data, pos, 9); return Action{}, false
	default:
		return Action{}, false
	}
}

func readCacheUnit(data []byte, pos *int) {
	advance(data, pos, 4) // unitId
	if *pos+4 > len(data) { return }
	itemsCount := int(binary.LittleEndian.Uint32(data[*pos:]))
	*pos += 4
	advance(data, pos, itemsCount*12)
	advance(data, pos, 4*4+2*4)
	advance(data, pos, 4+4+4+4)
	advance(data, pos, 4+4)
	if *pos+4 > len(data) { return }
	heroAbilCount := int(binary.LittleEndian.Uint32(data[*pos:]))
	*pos += 4
	advance(data, pos, heroAbilCount*8)
	advance(data, pos, 12)
	if *pos+4 > len(data) { return }
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
	if !ok { return [2]uint32{}, false }
	b, ok := ru32(data, pos)
	if !ok { return [2]uint32{}, false }
	return [2]uint32{a, b}, true
}

func rvec2(data []byte, pos *int) ([2]float32, bool) {
	x, ok := rf32(data, pos)
	if !ok { return [2]float32{}, false }
	y, ok := rf32(data, pos)
	if !ok { return [2]float32{}, false }
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

func rcache(data []byte, pos *int) (cacheData, bool) {
	fn := rzts(data, pos)
	mk := rzts(data, pos)
	k := rzts(data, pos)
	return cacheData{filename: fn, missionKey: mk, key: k}, true
}
