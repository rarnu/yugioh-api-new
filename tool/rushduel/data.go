package rushduel

import (
	. "github.com/isyscore/isc-gobase/isc"
)

type RushDuelDTO struct {
	Results ISCMap[ISCString, RushDuelPrintOuts] `json:"results"`
	Rows    ISCInt64                             `json:"rows"`
}

type RushDuelPrintOuts struct {
	Printouts RushDuelPrintOut `json:"printouts"`
}

type RushDuelPrintOut struct {
	Name         ISCList[ISCString] `json:"Name"`
	JapaneseName ISCList[ISCString] `json:"Japanese name"`
	PrimaryType  ISCList[TextAttr]  `json:"Primary type"`
	Attribute    ISCList[TextAttr]  `json:"[[Attribute]]"`
	Type         ISCList[TextAttr]  `json:"[[Type]]"`
	Level        ISCList[ISCString] `json:"[[Level]]"`
	Atk          ISCList[ISCString] `json:"[[ATK]]"`
	Def          ISCList[ISCString] `json:"[[DEF]]"`
	Status       ISCList[TextAttr]  `json:"Status"`
}

type TextAttr struct {
	FullText ISCString `json:"fulltext"`
}

type RushMagicTrapDTO struct {
	Results ISCMap[ISCString, RushMagicTrapPrintOuts] `json:"results"`
	Rows    ISCInt64                                  `json:"rows"`
}

type RushMagicTrapPrintOuts struct {
	Printouts RushMagicTrapPrintOut `json:"printouts"`
}

type RushMagicTrapPrintOut struct {
	Name         ISCList[ISCString] `json:"Name"`
	JapaneseName ISCList[ISCString] `json:"Japanese name"`
	CardType     ISCList[TextAttr]  `json:"[[Card type]]"`
	Property     ISCList[ISCString] `json:"[[Property]]"`
	Status       ISCList[TextAttr]  `json:"Status"`
}

var (
	GlobalSetNames         ISCList[ISCString]
	GlobalExistedCardIds   ISCList[ISCInt64]
	GlobalExistedSetNames  ISCList[ISCString]
	GlobalExistedCardNames ISCList[ISCString]
)

const (
	// data.type
	CardTypeMonster = 0x1
	CardTypeSpell   = 0x2
	CardTypeTrap    = 0x4

	// data.type
	MonsterTypeNormal      = 0x10
	MonsterTypeEffect      = 0x20
	MonsterTypeFusion      = 0x40
	MonsterTypeRitual      = 0x80
	MonsterSubTypeSpSummon = 0x2000000
	MonsterSubTypeMaximum  = 0x4000000

	// data.type
	IconRitual     = 0x80
	IconQuickPlay  = 0x10000
	IconContinuous = 0x20000
	IconEquip      = 0x40000
	IconField      = 0x80000
	IconCounter    = 0x100000

	// data.attribute
	AttributeEarth  = 0x1
	AttributeWater  = 0x2
	AttributeFire   = 0x4
	AttributeWind   = 0x8
	AttributeDark   = 0x10
	AttributeLight  = 0x20
	AttributeDivine = 0x40

	// data.race
	RaceGalaxy          = 0x40000000
	RaceCelestialKnight = 0x20000000
	RaceOmegaPsycho     = 0x10000000
	RaceHydragon        = 0x8000000
	RaceMagicalKnight   = 0x4000000
	RaceCyborg          = 0x2000000
	RaceCyberse         = 0x1000000
	RaceWyrm            = 0x800000
	RaceCreatorGod      = 0x400000
	RaceDivineBeast     = 0x200000
	RacePsychic         = 0x100000
	RaceReptile         = 0x80000
	RaceSeaSerpent      = 0x40000
	RaceFish            = 0x20000
	RaceDinosaur        = 0x10000
	RaceBeastWarrior    = 0x8000
	RaceBeast           = 0x4000
	RaceDragon          = 0x2000
	RaceThunder         = 0x1000
	RaceInsect          = 0x800
	RacePlant           = 0x400
	RaceWingedBeast     = 0x200
	RaceRock            = 0x100
	RacePyro            = 0x80
	RaceAqua            = 0x40
	RaceMachine         = 0x20
	RaceZombie          = 0x10
	RaceFiend           = 0x8
	RaceFairy           = 0x4
	RaceSpellcaster     = 0x2
	RaceWarrior         = 0x1
)
