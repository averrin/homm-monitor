package structs

import "time"

var VERSION = "0.6"
var PLUGIN_VERSION = 6
var URL = "161.35.209.12:28015"

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Town struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Grail bool   `json:"grail"`

	Coords

	ManaVortextUnused bool    `json:"manaVortextUnused"`
	Spells            [][]int `json:"spells"`
	BuiltThisTurn     bool    `json:"builtThisTurn"`
	GuardsValue       int     `json:"guardsValue"`

	GarrisonHero int  `json:"garrisonHero"`
	VisitingHero int  `json:"visitingHero"`
	HasFort      bool `json:"hasFort"`
}

type Item struct {
	Id   int `json:"id"`
	Type int `json:"subtype"`
}

type Hero struct {
	Coords

	Id   int    `json:"id"`
	Name string `json:"name"`

	Movement    int `json:"movement"`
	MaxMovement int `json:"maxMovement"`
	Experience  int `json:"experience"`
	Level       int `json:"level"`
	Gender      int `json:"gender"`

	Backpack        []Item `json:"backpack"`
	Weared          []Item `json:"weared"`
	LearnedSpells   []int  `json:"learned_spells"`
	AvailableSpells []int  `json:"available_spells"`

	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	SpellPower int `json:"spell_power"`
	Knowledge  int `json:"knowledge"`

	SecSkills []int `json:"sec_skills"`

	Mana      int `json:"mana"`
	ArmyValue int `json:"armyValue"`
	ArmyPower int `json:"armyPower"`

	MoraleBonus int `json:"moraleBonus"`
	LuckBonus   int `json:"luckBonus"`

	IsVisible  bool `json:"isVisible"`
	InGarrison bool `json:"inGarrison"`
}

type Combat struct {
	Hero           int `json:"hero"`
	HeroArmyValue  int `json:"heroArmyValue"`
	EnemyArmyValue int `json:"enemyArmyValue"`
}

type Resources struct {
	Gold    int `json:"gold"`
	Wood    int `json:"wood"`
	Mercury int `json:"mercury"`
	Ore     int `json:"ore"`
	Sulfur  int `json:"sulfur"`
	Crystal int `json:"crystal"`
	Gems    int `json:"gems"`

	GoldFound    int `json:"goldFound"`
	WoodFound    int `json:"woodFound"`
	MercuryFound int `json:"mercuryFound"`
	OreFound     int `json:"oreFound"`
	SulfurFound  int `json:"sulfurFound"`
	CrystalFound int `json:"crystalFound"`
	GemsFound    int `json:"gemsFound"`

	GoldSpent    int `json:"goldSpent"`
	WoodSpent    int `json:"woodSpent"`
	MercurySpent int `json:"mercurySpent"`
	OreSpent     int `json:"oreSpent"`
	SulfurSpent  int `json:"sulfurSpent"`
	CrystalSpent int `json:"crystalSpent"`
	GemsSpent    int `json:"gemsSpent"`
}

type HUD struct {
	HasWings           bool `json:"hasWings"`
	HasSpellbindersHat bool `json:"hasSpellbindersHat"`
	HasAlliance        bool `json:"hasAlliance"`
	HasShackles        bool `json:"hasShackles"`
	HasAOTD            bool `json:"hasAOTD"`

	HasTomeOfAir   bool `json:"hasTomeOfAir"`
	HasTomeOfFire  bool `json:"hasTomeOfFire"`
	HasTomeOfEarth bool `json:"hasTomeOfEarth"`

	HasIronFist bool `json:"hasIronFist"`

	HasDD          bool `json:"hasDD"`
	HasTP          bool `json:"hasTP"`
	HasFly         bool `json:"hasFly"`
	HasResurrect   bool `json:"hasResurrect"`
	HasArmageddon  bool `json:"hasArmageddon"`
	HasExpertEarth bool `json:"hasExpertEarth"`
	HasGrail       bool `json:"hasGrail"`
}

type Map struct {
	MapName           string `json:"mapName"`
	Size              int    `json:"size"`
	SubterraneanLevel int    `json:"subterraneanLevel"`
	VisionS           []int  `json:"visionS"`
	VisionU           []int  `json:"visionU"`
}

type Report struct {
	HUD HUD `json:"hud"`
	Map Map `json:"map"`

	Color int `json:"color"`

	Month int `json:"month"`
	Week  int `json:"week"`
	Day   int `json:"day"`

	Resources Resources `json:"resources"`

	TownsCount      int `json:"townsCount"`
	ObelisksVisited int `json:"obelisksVisited"`
	TotalObelisks   int `json:"totalObelisks"`

	TownsCountMax int `json:"townsCountMax"`

	TotalMPSpent int `json:"totalMPSpent"`
	TotalExp     int `json:"totalExp"`

	HeroesCount int     `json:"heroesCount"`
	Heroes      []*Hero `json:"heroes"`
	StartHero   int     `json:"startHero"`

	Towns     []*Town `json:"towns"`
	StartTown Town    `json:"startTown"`

	ArmyValue int `json:"armyValue"`

	Actions      int `json:"actions"`
	APM          int `json:"apm"`
	CleanActions int `json:"cleanActions"`
	CleanAPM     int `json:"cleanApm"`

	MaxAPM int `json:"maxApm"`
	AvgAPM int `json:"avgApm"`

	MatchId    string `json:"matchId"`
	PlayerName string `json:"playerName"`

	RecTime    time.Time `json:"recTime"`
	UpdateTime time.Time `json:"updateTime"`

	CurrentCombat Combat `json:"currentCombat"`

	Version       int    `json:"version"`
	ClientVersion string `json:"clientVersion"`
}

type HeroStore struct {
	LastMP  int
	MPSpent int
	InitExp int
	LastExp int
}

type CommenterState struct {
	Reports   map[int]Report `json:"reports"`
	MatchName string         `json:"matchName"`
	MatchDesc string         `json:"matchDesc"`
}

type Update struct {
	NewValue Report `rethinkdb:"new_val,omitempty"`
	OldValue Report `rethinkdb:"old_val,omitempty"`
}

type Match struct {
	Id      string    `rethinkdb:"id,omitempty"`
	Name    string    `rethinkdb:"name,omitempty"`
	Created time.Time `rethinkdb:"created,omitempty"`
	Winner  int       `rethinkdb:"winner,omitempty"`
}
