package w3ggo

// GameDataBlock is implemented by TimeslotEvent, ChatEvent, and LeaveGameEvent.
// Use a type switch to determine which variant you received.
type GameDataBlock interface {
	gameDataBlock()
}

// TimeslotEvent is emitted for each timeslot block during parsing.
type TimeslotEvent struct {
	TimeIncrement uint16
	CommandBlocks []CommandBlock
}

func (TimeslotEvent) gameDataBlock() {}

// CommandBlock carries the player ID and actions for one command block within
// a timeslot.
type CommandBlock struct {
	PlayerID uint8
	Actions  []Action
}

// ChatEvent is emitted for each chat message block.
type ChatEvent struct {
	PlayerID uint8
	Mode     uint32
	Message  string
}

func (ChatEvent) gameDataBlock() {}

// LeaveGameEvent is emitted when a player leaves.
type LeaveGameEvent struct {
	PlayerID uint8
	Reason   string
	Result   string
}

func (LeaveGameEvent) gameDataBlock() {}

// BasicReplayInfo holds the information available after metadata parsing but
// before game data is processed. It is passed to EventHandler.OnBasicReplayInformation.
type BasicReplayInfo struct {
	BuildNumber uint32
	Version     string
	GameName    string
	RandomSeed  uint32
	StartSpots  uint8
	Map         MapInfo
	Players     []BasicPlayerInfo
	Expansion   bool
}

// BasicPlayerInfo is a player record derived from the replay slot data.
type BasicPlayerInfo struct {
	PlayerID uint8
	Name     string
	TeamID   uint8
	Color    uint8
	Race     string
}

// EventHandler receives events during replay parsing. Implement this interface
// and pass it to ParseFileWithHandler or ParseBytesWithHandler.
type EventHandler interface {
	// OnBasicReplayInformation is called once after metadata is parsed, before
	// game data blocks are processed.
	OnBasicReplayInformation(info BasicReplayInfo)
	// OnGameDataBlock is called for each game data block before the parser
	// processes it internally.
	OnGameDataBlock(block GameDataBlock)
}
