package dto

import (
	. "github.com/isyscore/isc-gobase/isc"
)

func StrToCardType(astr ISCString) int64 {
	var t int64 = 0
	switch astr {
	case "pendulum":
		t = 0x1000000
	case "trap":
		t = 0x04
	case "spell":
		t = 0x02
	}
	return t
}

func StrToAttribute(astr ISCString) int {
	var a = 0
	switch astr {
	case "divine":
		a = 0x40
	case "dark":
		a = 0x20
	case "light":
		a = 0x10
	case "wind":
		a = 0x08
	case "fire":
		a = 0x04
	case "water":
		a = 0x02
	case "earth":
		a = 0x01
	}
	return a
}

func StrToIcon(astr ISCString) int64 {
	var i int64 = 0
	switch astr {
	case "counter":
		i = 0x100000
	case "field":
		i = 0x80000
	case "equip":
		i = 0x40000
	case "continuous":
		i = 0x20000
	case "quick-play":
		i = 0x10000
	case "ritual":
		i = 0x80
	}
	return i
}

func StrToSubType(astr ISCString) int64 {
	var s int64 = 0
	switch astr {
	case "spsummon":
		s = 0x2000000
	case "link":
		s = 0x4000000
	case "pendulum":
		s = 0x1000000
	case "xyz":
		s = 0x800000
	case "synchro":
		s = 0x2000
	case "ritual":
		s = 0x80
	case "fusion":
		s = 0x40
	case "toon":
		s = 0x400000
	case "flip":
		s = 0x200000
	case "tuner":
		s = 0x1000
	case "gemini":
		s = 0x800
	case "union":
		s = 0x400
	case "spirit":
		s = 0x200
	case "effect":
		s = 0x20
	case "normal":
		s = 0x10
	}
	return s
}

func StrToRace(astr ISCString) int64 {
	var r int64 = 0
	switch astr {
	case "cyberse":
		r = 0x1000000
	case "wyrm":
		r = 0x800000
	case "creatorGod":
		r = 0x400000
	case "divineBeast":
		r = 0x200000
	case "psychic":
		r = 0x100000
	case "reptile":
		r = 0x80000
	case "seaSerpent":
		r = 0x40000
	case "fish":
		r = 0x20000
	case "dinosaur":
		r = 0x10000
	case "beastWarrior":
		r = 0x8000
	case "beast":
		r = 0x4000
	case "dragon":
		r = 0x2000
	case "thunder":
		r = 0x1000
	case "insect":
		r = 0x800
	case "plant":
		r = 0x400
	case "wingedBeast":
		r = 0x200
	case "rock":
		r = 0x100
	case "pyro":
		r = 0x80
	case "aqua":
		r = 0x40
	case "machine":
		r = 0x20
	case "zombie":
		r = 0x10
	case "fiend":
		r = 0x8
	case "fairy":
		r = 0x4
	case "spellcaster":
		r = 0x2
	case "warrior":
		r = 0x1
	}
	return r
}

func StrToMonsterType(astr ISCString) int64 {
	var m int64 = 0
	switch astr {
	case "link":
		m = 0x4000000
	case "xyz":
		m = 0x800000
	case "token":
		m = 0x4000
	case "synchro":
		m = 0x2000
	case "ritual":
		m = 0x80
	case "fusion":
		m = 0x40
	case "effect":
		m = 0x20
	case "normal":
		m = 0x10
	}
	return m
}
