package dto

import (
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server/rsp"
)

type CommonCount struct {
	CardCount int `json:"cardCount"`
	KanaCount int `json:"kanaCount"`
	SetCount  int `json:"setCount"`
}

type RespCommonCount struct {
	rsp.ResponseBase
	Data CommonCount `json:"data"`
}

type RespCommonString struct {
	rsp.ResponseBase
	Data isc.ISCString `json:"data"`
}
