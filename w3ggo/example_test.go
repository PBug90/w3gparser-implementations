package w3ggo_test

import (
	"fmt"

	w3g "github.com/PBug90/w3gparser-implementations/w3ggo"
)

// ExampleParseFile demonstrates parsing a replay file and reading basic fields.
func ExampleParseFile() {
	out, err := w3g.ParseFile("replay.w3g")
	if err != nil {
		panic(err)
	}

	fmt.Println("Map:     ", out.Map.File)
	fmt.Println("Matchup: ", out.Matchup)
	fmt.Printf("Duration: %ds\n", out.Duration/1000)

	for _, p := range out.Players {
		fmt.Printf("  %s (%s) — APM %d\n", p.Name, p.Race, p.APM)
	}
}

// ExampleParseBytes demonstrates parsing from an in-memory byte slice.
func ExampleParseBytes() {
	data := []byte{ /* raw .w3g bytes */ }

	out := w3g.ParseBytes(data)
	if out == nil {
		panic("failed to parse replay")
	}

	fmt.Println(out.Map.File)
}

// statsHandler counts timeslots and actions during parsing.
type statsHandler struct {
	timeslots int
	actions   int
}

func (h *statsHandler) OnBasicReplayInformation(info w3g.BasicReplayInfo) {
	fmt.Printf("[event] map: %s  players: %d\n", info.Map.File, len(info.Players))
}

func (h *statsHandler) OnGameDataBlock(block w3g.GameDataBlock) {
	if ts, ok := block.(w3g.TimeslotEvent); ok {
		h.timeslots++
		for _, cmd := range ts.CommandBlocks {
			h.actions += len(cmd.Actions)
		}
	}
}

// ExampleParseFileWithHandler demonstrates the event-based handler API.
func ExampleParseFileWithHandler() {
	h := &statsHandler{}
	w3g.ParseFileWithHandler("replay.w3g", h)
	fmt.Printf("[event] %d timeslots, %d actions\n", h.timeslots, h.actions)
}
