package w3ggo

import (
	"encoding/binary"
	"math"
)

type BufParser struct {
	data []byte
	pos  int
}

func newBufParser(data []byte) *BufParser {
	return &BufParser{data: data}
}

func (p *BufParser) remaining() int {
	n := len(p.data) - p.pos
	if n < 0 {
		return 0
	}
	return n
}

func (p *BufParser) skip(n int) {
	p.pos += n
}

func (p *BufParser) setPos(pos int) {
	p.pos = pos
}

func (p *BufParser) peekU8() (byte, bool) {
	if p.pos >= len(p.data) {
		return 0, false
	}
	return p.data[p.pos], true
}

func (p *BufParser) readU8() (byte, bool) {
	if p.pos >= len(p.data) {
		return 0, false
	}
	v := p.data[p.pos]
	p.pos++
	return v, true
}

func (p *BufParser) readU16LE() (uint16, bool) {
	if p.pos+2 > len(p.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint16(p.data[p.pos:])
	p.pos += 2
	return v, true
}

func (p *BufParser) readU32LE() (uint32, bool) {
	if p.pos+4 > len(p.data) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(p.data[p.pos:])
	p.pos += 4
	return v, true
}

func (p *BufParser) readF32LE() (float32, bool) {
	if p.pos+4 > len(p.data) {
		return 0, false
	}
	bits := binary.LittleEndian.Uint32(p.data[p.pos:])
	p.pos += 4
	return math.Float32frombits(bits), true
}

func (p *BufParser) readBytes(n int) ([]byte, bool) {
	if p.pos+n > len(p.data) {
		return nil, false
	}
	s := p.data[p.pos : p.pos+n]
	p.pos += n
	return s, true
}

func (p *BufParser) readZeroTermString() string {
	start := p.pos
	for p.pos < len(p.data) && p.data[p.pos] != 0 {
		p.pos++
	}
	s := string(p.data[start:p.pos])
	if p.pos < len(p.data) {
		p.pos++ // consume null terminator
	}
	return s
}

func (p *BufParser) readStringOfLengthAsHex(n int) string {
	b, ok := p.readBytes(n)
	if !ok {
		return ""
	}
	return encodeHex(b)
}

func (p *BufParser) readStringOfLengthUTF8(n int) string {
	b, ok := p.readBytes(n)
	if !ok {
		return ""
	}
	return string(b)
}

func (p *BufParser) readFourCC() [4]byte {
	var arr [4]byte
	for i := range arr {
		b, _ := p.readU8()
		arr[i] = b
	}
	return arr
}

func encodeHex(b []byte) string {
	const hextable = "0123456789abcdef"
	dst := make([]byte, len(b)*2)
	for i, v := range b {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0f]
	}
	return string(dst)
}
