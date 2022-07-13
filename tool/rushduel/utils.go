package rushduel

import (
	. "github.com/isyscore/isc-gobase/isc"
	"regexp"
)

func parseJapaneseName(name ISCString) (kanji ISCString, kk ISCString) {
	// "[暴(あば)]れ[牛(うし)][鬼(おに)]"
	tmp := name.ReplaceAll("<ruby lang=\"ja\">", "").ReplaceAll("</ruby>", "")
	tmp = tmp.ReplaceAll("<rp>（</rp>", "(").ReplaceAll("<rp>）</rp>", ")]")
	tmp = tmp.ReplaceAll("<rt>", "").ReplaceAll("</rt>", "")
	tmp = tmp.ReplaceAll("<rb>", "[").ReplaceAll("</rb>", "")
	kj := removeKana(tmp)
	return kj, tmp
}

func removeKana(str ISCString) ISCString {
	re, _ := regexp.Compile("\\[.*?\\(.*?\\)]")
	dest := re.ReplaceAllStringFunc(string(str), func(s string) string {
		ss := ISCString(s)
		ss = ss.ReplaceAll("[", "")
		ss = ss.SubStringBefore("(")
		return string(ss)
	})
	return ISCString(dest)
}

func effectCardNames(str ISCString) ISCList[ISCString] {
	re, _ := regexp.Compile("「.*?」")
	sa := ListToMapFrom[string, ISCString](re.FindAllString(string(str), -1))
	return sa.Map(func(item string) ISCString {
		return ISCString(item).ReplaceAll("「", "").ReplaceAll("」", "")
	}).Distinct()
}

func parsePrimaryType(ts []TextAttr) ISCInt64 {
	r := CardTypeMonster
	for _, t := range ts {
		if t.FullText.Contains("Normal") {
			r += MonsterTypeNormal
		}
		if t.FullText.Contains("Effect") {
			r += MonsterTypeEffect
		}
		if t.FullText.Contains("Fusion") {
			r += MonsterTypeFusion
		}
		if t.FullText.Contains("Ritual") {
			r += MonsterTypeRitual
		}
	}
	return ISCInt64(r)
}

func parseMagicTrapType(prop ISCString) ISCInt64 {
	arr := ListToMapFrom[ISCString, ISCString](prop.Split(" ")).Map(func(item ISCString) ISCString {
		return item.ToLower()
	})
	r := 0
	if arr[1] == "spell" {
		r += CardTypeSpell
	} else if arr[1] == "trap" {
		r += CardTypeTrap
	}
	switch arr[0] {
	case "counter":
		r += IconCounter
	case "filed":
		r += IconField
	case "equip":
		r += IconEquip
	case "continuous":
		r += IconContinuous
	case "quick-play":
		r += IconQuickPlay
	case "ritual":
		r += IconRitual
	}
	return ISCInt64(r)
}

func parseAttribute(a ISCString) ISCInt64 {
	switch a {
	case "DIVINE":
		return 0x40
	case "DARK":
		return 0x20
	case "LIGHT":
		return 0x10
	case "WIND":
		return 0x8
	case "FIRE":
		return 0x4
	case "WATER":
		return 0x2
	case "EARTH":
		return 0x1
	}
	return 0x0
}

func parseRace(t ISCString) ISCInt64 {
	t0 := t.ReplaceAll("-", "").ToLower()
	switch t0 {
	case "galaxy":
		return 0x40000000
	case "celestialknight":
		return 0x20000000
	case "omegapsycho":
		return 0x10000000
	case "hydragon":
		return 0x8000000
	case "magicalknight":
		return 0x4000000
	case "cyborg":
		return 0x2000000
	case "cyberse":
		return 0x1000000
	case "wyrm":
		return 0x800000
	case "creatorgod":
		return 0x400000
	case "divinebeast":
		return 0x200000
	case "psychic":
		return 0x100000
	case "reptile":
		return 0x80000
	case "seaserpent":
		return 0x40000
	case "fish":
		return 0x20000
	case "dinosaur":
		return 0x10000
	case "beastwarrior":
		return 0x8000
	case "beast":
		return 0x4000
	case "dragon":
		return 0x2000
	case "thunder":
		return 0x1000
	case "insect":
		return 0x800
	case "plant":
		return 0x400
	case "wingedbeast":
		return 0x200
	case "rock":
		return 0x100
	case "pyro":
		return 0x80
	case "aqua":
		return 0x40
	case "machine":
		return 0x20
	case "zombie":
		return 0x10
	case "fiend":
		return 0x8
	case "fairy":
		return 0x4
	case "spellcaster":
		return 0x2
	case "warrior":
		return 0x1
	}
	return 0x0
}
