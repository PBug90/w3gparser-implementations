package w3ggo

// ObserverMode mirrors the TS ObserverMode enum.
type ObserverMode string

const (
	ObserverModeOnDefeat ObserverMode = "ON_DEFEAT"
	ObserverModeFull     ObserverMode = "FULL"
	ObserverModeReferees ObserverMode = "REFEREES"
	ObserverModeNone     ObserverMode = "NONE"
)

// ChatMode mirrors the TS ChatMessageMode.
type ChatMode string

const (
	ChatModeAll       ChatMode = "All"
	ChatModeTeam      ChatMode = "Team"
	ChatModeObservers ChatMode = "Observers"
	ChatModePrivate   ChatMode = "Private"
)

// MapInfo holds map metadata.
type MapInfo struct {
	Path         string `json:"path"`
	File         string `json:"file"`
	Checksum     string `json:"checksum"`
	ChecksumSha1 string `json:"checksumSha1"`
}

// APMConfig holds APM tracking config.
type APMConfig struct {
	TrackingInterval int `json:"trackingInterval"`
}

// Settings mirrors the TS settings object.
type Settings struct {
	ObserverMode          ObserverMode `json:"observerMode"`
	Referees              bool         `json:"referees"`
	FixedTeams            bool         `json:"fixedTeams"`
	FullSharedUnitControl bool         `json:"fullSharedUnitControl"`
	AlwaysVisible         bool         `json:"alwaysVisible"`
	HideTerrain           bool         `json:"hideTerrain"`
	MapExplored           bool         `json:"mapExplored"`
	TeamsTogether         bool         `json:"teamsTogether"`
	RandomHero            bool         `json:"randomHero"`
	RandomRaces           bool         `json:"randomRaces"`
	Speed                 int          `json:"speed"`
}

// ChatMessage is a parsed in-game chat message.
type ChatMessage struct {
	PlayerName string   `json:"playerName"`
	PlayerID   int      `json:"playerId"`
	Mode       ChatMode `json:"mode"`
	TimeMS     int      `json:"timeMS"`
	Message    string   `json:"message"`
}

// OrderEntry is a single train/build/upgrade/item event.
type OrderEntry struct {
	ID string `json:"id"`
	MS int    `json:"ms"`
}

// Summary holds both a frequency map and an ordered list.
type Summary struct {
	SummaryMap map[string]int `json:"summary"`
	Order      []OrderEntry   `json:"order"`
}

func newSummary() Summary {
	return Summary{SummaryMap: make(map[string]int), Order: []OrderEntry{}}
}

func (s *Summary) add(id string, ms int) {
	s.SummaryMap[id]++
	s.Order = append(s.Order, OrderEntry{ID: id, MS: ms})
}

// AbilityOrderEntry is one entry in a hero's ability learning sequence.
type AbilityOrderEntry struct {
	Type  string `json:"type"`
	Time  int    `json:"time"`
	Value string `json:"value,omitempty"` // omitted for retraining entries
}

// RetrainingSnapshot captures the state of a hero's abilities at the time of retraining.
type RetrainingSnapshot struct {
	Time      int            `json:"time"`
	Abilities map[string]int `json:"abilities"`
}

// HeroInfo holds tracked hero data.
type HeroInfo struct {
	ID                string               `json:"id"`
	Level             int                  `json:"level"`
	Abilities         map[string]int       `json:"abilities"`
	RetrainingHistory []RetrainingSnapshot `json:"retrainingHistory"`
	AbilityOrder      []AbilityOrderEntry  `json:"abilityOrder"`
}

// GroupHotkey tracks assign/use count for a control group key.
type GroupHotkey struct {
	Assigned int `json:"assigned"`
	Used     int `json:"used"`
}

// Actions holds categorized action counts.
type Actions struct {
	Timed        []int `json:"timed"`
	AssignGroup  int   `json:"assigngroup"`
	RightClick   int   `json:"rightclick"`
	Basic        int   `json:"basic"`
	BuildTrain   int   `json:"buildtrain"`
	Ability      int   `json:"ability"`
	Item         int   `json:"item"`
	Select       int   `json:"select"`
	RemoveUnit   int   `json:"removeunit"`
	Subgroup     int   `json:"subgroup"`
	SelectHotkey int   `json:"selecthotkey"`
	ESC          int   `json:"esc"`
}

// ResourceTransfer records a gold/lumber transfer between players.
type ResourceTransfer struct {
	Slot       int    `json:"slot"`
	PlayerID   int    `json:"playerId"`
	PlayerName string `json:"playerName"`
	Gold       int    `json:"gold"`
	Lumber     int    `json:"lumber"`
	MSElapsed  int    `json:"msElapsed"`
}

// PlayerOutput is the final exported player structure (mirrors Player.toJSON).
type PlayerOutput struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	TeamID            int                `json:"teamid"`
	Color             string             `json:"color"`
	Race              string             `json:"race"`
	RaceDetected      string             `json:"raceDetected"`
	Units             Summary            `json:"units"`
	Buildings         Summary            `json:"buildings"`
	Items             Summary            `json:"items"`
	Upgrades          Summary            `json:"upgrades"`
	Heroes            []HeroInfo         `json:"heroes"`
	Actions           Actions            `json:"actions"`
	GroupHotkeys      map[int]GroupHotkey `json:"groupHotkeys"`
	ResourceTransfers []ResourceTransfer  `json:"resourceTransfers"`
	APM               int                `json:"apm"`
	CurrentTimePlayed int                `json:"currentTimePlayed"`
}

// ParserOutput is the final result of parsing a replay. Mirrors the TS ParserOutput.
type ParserOutput struct {
	ID            string        `json:"id"`
	Gamename      string        `json:"gamename"`
	RandomSeed    int           `json:"randomseed"`
	StartSpots    int           `json:"startSpots"`
	Observers     []string      `json:"observers"`
	Players       []PlayerOutput `json:"players"`
	Matchup       string        `json:"matchup"`
	Creator       string        `json:"creator"`
	Type          string        `json:"type"`
	Chat          []ChatMessage `json:"chat"`
	APM           APMConfig     `json:"apm"`
	Map           MapInfo       `json:"map"`
	BuildNumber   int           `json:"buildNumber"`
	Version       string        `json:"version"`
	Duration      int           `json:"duration"`
	Expansion     bool          `json:"expansion"`
	ParseTime     int64         `json:"parseTime"`
	WinningTeamID int           `json:"winningTeamId"`
	Settings      Settings      `json:"settings"`
}
