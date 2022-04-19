package database

import (
	"fmt"
	"ygoapi/config"
	"ygoapi/japanese"

	. "github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"golang.org/x/text/width"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type YgoNameDB struct {
	om *gorm.DB
}

var YgoName YgoNameDB

func NewYgoName() {
	dsnFmt := "%s:%s@tcp(%s:%d)/YugiohAPI2?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(dsnFmt, config.MySQLConfig.User, config.MySQLConfig.Password, config.MySQLConfig.Host, config.MySQLConfig.Port)
	ygo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	logger.Info("YgoNameDB connected (%v).", ygo)
	YgoName = YgoNameDB{om: ygo}
}

func (y YgoNameDB) KanaCount() int {
	var count int
	y.om.Raw("select count(1) from card_name_texts").First(&count)
	return count
}

func (y YgoNameDB) SetCount() int {
	var count int
	y.om.Raw("select count(1) from set_name_texts").First(&count)
	return count
}

func (y YgoNameDB) NameKanjiKana(name ISCString) ISCString {
	var kk ISCString = ""
	y.om.Raw(fmt.Sprintf("select kk from card_name_texts where kanji = '%s' or kanji = '%s' or kanji = '%s'",
		name, width.Narrow.String(string(name)), width.Widen.String(string(name)))).Scan(&kk)
	return kk
}

func (y YgoNameDB) SetKanjiKana(name ISCString) ISCString {
	var kk ISCString = ""
	y.om.Raw(fmt.Sprintf("select kk from set_name_texts where kanji = '%s' or kanji = '%s' or kanji = '%s'",
		name, width.Narrow.String(string(name)), width.Widen.String(string(name)))).Scan(&kk)
	return kk
}

func (y YgoNameDB) EffectKanjiKana(name ISCString) ISCString {
	cn := japanese.EffectCardNames(name)
	japanese.SortByLength(cn)
	e2 := name
	for i, e := range cn {
		e2 = e2.ReplaceAll(string(e), fmt.Sprintf("{{%d}}", i))
	}
	e2 = japanese.Kana(e2)
	for i, e := range cn {
		isToken := false
		tmp := e
		if tmp.EndsWith("トークン") {
			isToken = true
		}
		kk := y.NameKanjiKana(tmp)
		if kk == "" {
			tmp = tmp.ReplaceAll("トークン", "").TrimSpace()
			kk = y.SetKanjiKana(tmp)
			if kk == "" {
				kk = japanese.Kana(tmp)
			}
		}
		if isToken && !kk.EndsWith("トークン") {
			kk += "トークン"
		}
		e2 = e2.ReplaceAll(fmt.Sprintf("{{%d}}", i), string(kk))
	}
	return e2
}

func (y YgoNameDB) NormalKanjiKana(name ISCString) ISCString {
	cn := japanese.EffectCardNames(name)
	japanese.SortByLength(cn)
	e2 := name
	for i, e := range cn {
		e2 = e2.ReplaceAll(string(e), fmt.Sprintf("{{%d}}", i))
	}
	e2 = japanese.NormalKana(e2)
	for i, e := range cn {
		isToken := false
		tmp := e
		if tmp.EndsWith("トークン") {
			isToken = true
		}
		kk := y.NameKanjiKana(tmp)
		if kk == "" {
			tmp = tmp.ReplaceAll("トークン", "").TrimSpace()
			kk = y.SetKanjiKana(tmp)
			if kk == "" {
				kk = japanese.Kana(tmp)
			}
		}
		if isToken && !kk.EndsWith("トークン") {
			kk += "トークン"
		}
		e2 = e2.ReplaceAll(fmt.Sprintf("{{%d}}", i), string(kk))
	}
	return e2
}
