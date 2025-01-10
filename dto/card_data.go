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

type CardEffectData struct {
	Id                 int64                   `json:"id"`                 // 卡片 ID
	CardName           string                  `json:"cardName"`           // 卡片名称
	EffectText         string                  `json:"effectText"`         // 完整的效果文本
	EffectCount        int                     `json:"effectCount"`        // 效果条数
	Summon             string                  `json:"summon"`             // 召唤条件
	Limitation         string                  `json:"limitation"`         // 限制文本
	PendulumLimitation string                  `json:"pendulumLimitation"` // 灵摆限制文本
	Detail             []*CardEffectDetailData `json:"detail" gorm:"-"`    // 卡片效果明细
}

type CardEffectDetailData struct {
	IsPendulumEffect bool     `json:"isPendulumEffect"`       // 是否灵摆效果
	EffectId         int      `json:"effectId"`               // 效果序号
	EffectCondition  string   `json:"effectCondition"`        // 发动条件
	EffectTreatment  string   `json:"effectTreatment"`        // 效果处理
	EffectLimitation string   `json:"effectLimitation"`       // 发动后的限制或后续动作
	EffectOptions    []string `json:"effectOptions" gorm:"-"` // 效果中的可选项
	EffectOpt        string   `json:"-" gorm:"column:effect_options"`
}
