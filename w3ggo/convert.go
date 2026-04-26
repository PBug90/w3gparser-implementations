package w3ggo

import (
	"fmt"
	"strings"
)

var playerColors = [...]string{
	"#ff0303", "#0042ff", "#1ce6b9", "#540081", "#fffc00",
	"#fe8a0e", "#20c000", "#e55bb0", "#959697", "#7ebff1",
	"#106246", "#4a2a04", "#9b0000", "#0000c3", "#00eaff",
	"#be00fe", "#ebcd87", "#f8a48b", "#bfff80", "#dcb9eb",
	"#282828", "#ebf0ff", "#00781e", "#a46f33",
}

func playerColor(color int) string {
	if color >= 0 && color < len(playerColors) {
		return playerColors[color]
	}
	return "000000"
}

func gameVersion(version int) string {
	if version == 10030 {
		return "1.30.2+"
	} else if version > 10030 && version < 10100 {
		s := fmt.Sprintf("%d", version)
		return "1." + s[len(s)-2:]
	} else if version >= 10100 {
		s := fmt.Sprintf("%d", version)
		return "2." + s[len(s)-2:]
	}
	return fmt.Sprintf("1.%d", version)
}

func mapFilename(mapPath string) string {
	// Normalise separators
	p := strings.ReplaceAll(mapPath, "\\", "/")
	idx := strings.LastIndex(p, "/")
	var candidate string
	if idx >= 0 {
		candidate = p[idx+1:]
	} else {
		candidate = p
	}
	if strings.HasSuffix(candidate, ".w3x") || strings.HasSuffix(candidate, ".w3m") {
		return candidate
	}
	return ""
}
