package dto

import (
	. "github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server/rsp"
)

type CardData struct {
	Id        int64     `json:"id"`
	Name      ISCString `json:"name"`
	Desc      ISCString `json:"desc"`
	Type      int64     `json:"type"`
	Atk       int       `json:"atk"`
	Def       int       `json:"def"`
	Level     int       `json:"level"`
	Race      int64     `json:"race"`
	Attribute int       `json:"attribute"`
	Abbr      ISCString `json:"setid"`
}

type CardName struct {
	Id   int64     `json:"id"`
	Name ISCString `json:"name"`
}

type CardNameData struct {
	rsp.ResponseBase
	Data []*CardName `json:"data"`
}

type SearchCardData struct {
	rsp.ResponseBase
	Data []*CardData `json:"data"`
}
