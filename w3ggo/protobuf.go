package w3ggo

type protoDecoder struct {
	data []byte
	pos  int
}

func newProtoDecoder(data []byte) *protoDecoder {
	return &protoDecoder{data: data}
}

func (d *protoDecoder) remaining() int {
	n := len(d.data) - d.pos
	if n < 0 {
		return 0
	}
	return n
}

func (d *protoDecoder) readVarint() (uint64, bool) {
	var result uint64
	var shift uint
	for {
		if d.pos >= len(d.data) {
			return 0, false
		}
		b := d.data[d.pos]
		d.pos++
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return 0, false
		}
	}
	return result, true
}

func (d *protoDecoder) readLengthDelimited() ([]byte, bool) {
	n, ok := d.readVarint()
	if !ok {
		return nil, false
	}
	length := int(n)
	if d.pos+length > len(d.data) {
		return nil, false
	}
	s := d.data[d.pos : d.pos+length]
	d.pos += length
	return s, true
}

func (d *protoDecoder) skipField(wireType uint64) {
	switch wireType {
	case 0:
		d.readVarint()
	case 1:
		d.pos += 8
	case 2:
		d.readLengthDelimited()
	case 5:
		d.pos += 4
	}
}

func (d *protoDecoder) decodePlayerFlat() (uint32, string, string, bool) {
	var playerID uint32
	var battleTag, clan string
	for d.remaining() > 0 {
		tag, ok := d.readVarint()
		if !ok {
			break
		}
		fieldNumber := tag >> 3
		wireType := tag & 0x7
		switch fieldNumber {
		case 1:
			v, ok := d.readVarint()
			if !ok {
				return 0, "", "", false
			}
			playerID = uint32(v)
		case 2:
			b, ok := d.readLengthDelimited()
			if !ok {
				return 0, "", "", false
			}
			battleTag = string(b)
		case 3:
			b, ok := d.readLengthDelimited()
			if !ok {
				return 0, "", "", false
			}
			clan = string(b)
		default:
			d.skipField(wireType)
		}
	}
	return playerID, battleTag, clan, true
}

func (d *protoDecoder) decodePlayerNested() []playerMeta {
	var players []playerMeta
	for d.remaining() > 0 {
		tag, ok := d.readVarint()
		if !ok {
			break
		}
		fieldNumber := tag >> 3
		wireType := tag & 0x7
		if fieldNumber == 1 && wireType == 2 {
			sub, ok := d.readLengthDelimited()
			if !ok {
				break
			}
			subDec := newProtoDecoder(sub)
			id, bt, cl, ok := subDec.decodePlayerFlat()
			if ok {
				players = append(players, playerMeta{id: id, battleTag: bt, clan: cl})
			}
		} else {
			d.skipField(wireType)
		}
	}
	return players
}

type playerMeta struct {
	id        uint32
	battleTag string
	clan      string
}

func decodePlayerMetadata(data []byte) []playerMeta {
	if len(data) == 0 {
		return nil
	}
	// peek first tag
	d := newProtoDecoder(data)
	firstTag, ok := d.readVarint()
	if !ok {
		return nil
	}
	fieldNumber := firstTag >> 3
	wireType := firstTag & 0x7

	if fieldNumber == 1 && wireType == 2 {
		// Nested format
		d2 := newProtoDecoder(data)
		return d2.decodePlayerNested()
	}
	// Flat format
	d3 := newProtoDecoder(data)
	id, bt, cl, ok := d3.decodePlayerFlat()
	if ok {
		return []playerMeta{{id: id, battleTag: bt, clan: cl}}
	}
	return nil
}
