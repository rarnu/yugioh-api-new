package dto

import (
	. "github.com/isyscore/isc-gobase/isc"
)

type ReqSearch struct {
	Key         ISCString `json:"key"`
	CardType    ISCString `json:"cardtype"`
	Attribute   ISCString `json:"attribute"`
	Icon        ISCString `json:"icon"`
	SubType     ISCString `json:"subtype"`
	Race        ISCString `json:"race"`
	MonsterType ISCString `json:"monstertype"`
	Lang        ISCString `json:"lang"`
}

type ReqSearchOrigin struct {
	Key         ISCString
	CardType    int64
	Attribute   int
	Icon        int64
	SubType     int64
	Race        int64
	MonsterType int64
	Lang        ISCString
}

type ReqYdkFind struct {
	ByEffect bool      `json:"byEffect"`
	Key      ISCString `json:"key"`
	Lang     ISCString `json:"lang"`
}

type ReqYdkNames struct {
	Lang ISCString      `json:"lang"`
	Ids  ISCList[int64] `json:"ids"`
}

type ReqKKName struct {
	Name ISCString `json:"name"`
}

func ReqSearchToOrigin(r ReqSearch) ReqSearchOrigin {
	return ReqSearchOrigin{
		Key:         r.Key,
		CardType:    StrToCardType(r.CardType),
		Attribute:   StrToAttribute(r.Attribute),
		Icon:        StrToIcon(r.Icon),
		SubType:     StrToSubType(r.SubType),
		Race:        StrToRace(r.Race),
		MonsterType: StrToMonsterType(r.MonsterType),
		Lang:        r.Lang,
	}
}
