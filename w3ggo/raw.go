package w3ggo

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"
)

type RawHeader struct {
	compressedSize   uint32
	headerVersion    string
	decompressedSize uint32
	blockCount       uint32
}

type RawSubHeader struct {
	gameIdentifier  string
	version         uint32
	buildNo         uint16
	replayLengthMS  uint32
}

type dataBlock struct {
	blockSize            uint16
	blockDecompressedSize uint16
	blockContent         []byte
}

func parseRaw(input []byte) (*RawHeader, *RawSubHeader, []dataBlock) {
	magic := []byte("Warcraft III recorded game")
	start := findSubsequence(input, magic)
	if start < 0 {
		return nil, nil, nil
	}

	p := newBufParser(input)
	p.setPos(start)
	p.readZeroTermString() // consume the null-terminated magic string
	p.skip(4)              // unknown 4 bytes

	compressedSize, ok := p.readU32LE()
	if !ok {
		return nil, nil, nil
	}
	headerVersion := p.readStringOfLengthAsHex(4)
	decompressedSize, ok := p.readU32LE()
	if !ok {
		return nil, nil, nil
	}
	blockCount, ok := p.readU32LE()
	if !ok {
		return nil, nil, nil
	}

	header := &RawHeader{
		compressedSize:   compressedSize,
		headerVersion:    headerVersion,
		decompressedSize: decompressedSize,
		blockCount:       blockCount,
	}

	// Subheader
	gameIdentifier := p.readStringOfLengthUTF8(4)
	version, ok := p.readU32LE()
	if !ok {
		return nil, nil, nil
	}
	buildNo, ok := p.readU16LE()
	if !ok {
		return nil, nil, nil
	}
	p.skip(2)
	replayLengthMS, ok := p.readU32LE()
	if !ok {
		return nil, nil, nil
	}
	p.skip(4)

	subHeader := &RawSubHeader{
		gameIdentifier: gameIdentifier,
		version:        version,
		buildNo:        buildNo,
		replayLengthMS: replayLengthMS,
	}

	isReforged := buildNo >= 6089

	var blocks []dataBlock
	for p.remaining() > 0 {
		blockSize, ok := p.readU16LE()
		if !ok {
			break
		}
		if isReforged {
			p.skip(2)
		}
		blockDecompSize, ok := p.readU16LE()
		if !ok {
			break
		}
		if isReforged {
			p.skip(6)
		} else {
			p.skip(4)
		}
		content, ok := p.readBytes(int(blockSize))
		if !ok {
			break
		}
		if blockDecompSize == 8192 {
			contentCopy := make([]byte, len(content))
			copy(contentCopy, content)
			blocks = append(blocks, dataBlock{
				blockSize:            blockSize,
				blockDecompressedSize: blockDecompSize,
				blockContent:         contentCopy,
			})
		}
	}

	return header, subHeader, blocks
}

func decompressBlocks(blocks []dataBlock) []byte {
	var result []byte
	for _, block := range blocks {
		if len(block.blockContent) == 0 {
			continue
		}
		r, err := zlib.NewReader(bytes.NewReader(block.blockContent))
		if err != nil {
			continue
		}
		buf, err := io.ReadAll(r)
		r.Close()
		// Accept io.ErrUnexpectedEOF: some WC3 replay blocks omit the zlib Adler-32 checksum
		if (err == nil || isUnexpectedEOF(err)) && len(buf) > 0 {
			result = append(result, buf...)
		}
	}
	return result
}

func isUnexpectedEOF(err error) bool {
	if err == io.ErrUnexpectedEOF {
		return true
	}
	return strings.Contains(err.Error(), "unexpected EOF")
}

func findSubsequence(haystack, needle []byte) int {
	if len(needle) == 0 {
		return 0
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if bytes.Equal(haystack[i:i+len(needle)], needle) {
			return i
		}
	}
	return -1
}
