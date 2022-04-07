package database

import (
	"fmt"
	. "github.com/isyscore/isc-gobase/isc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"tool/html"
)

type OmegaDB struct {
	om *gorm.DB
}

type CardData struct {
	Id   int64
	Name string
}

type SetData struct {
	Officialcode int64
	Name         string
}

var Omega OmegaDB

func NewOmega() {
	omega, err := gorm.Open(sqlite.Open("./OmegaDB.cdb"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("OmegaDB connected (%v).", omega)
	Omega = OmegaDB{om: omega}
}

func (o OmegaDB) CardIdList() ISCList[int64] {
	var list ISCList[int64]
	o.om.Raw("select id from ja_texts").Scan(&list)
	return list
}

func (o OmegaDB) SetIdList() ISCSet[int64] {
	var list ISCList[int64]
	o.om.Raw("select officialcode from setcodes group by officialcode").Scan(&list)
	return list.ToSet()
}

func (o OmegaDB) GetCardSQLByIds(ids ISCList[int64]) ISCList[string] {
	var list ISCList[*CardData]
	o.om.Raw("select id, name from ja_texts where id in (?)", ids).Scan(&list)
	return ListToMapFrom[*CardData, string](list).Map(func(card *CardData) string {
		kk := html.CrawlingKaka(card.Name)
		return fmt.Sprintf("insert into card_name_texts(id, kanji, kk, donetime) values (%d, '%s', '%s', 0);", card.Id, card.Name, kk)
	})
}

func (o OmegaDB) GetSetSQLByIds(ids ISCList[int64]) ISCList[string] {
	var list ISCList[*SetData]
	o.om.Raw("select officialcode, name from setcodes where officialcode in (?) group by officialcode", ids).Scan(&list)
	return ListToMapFrom[*SetData, string](list).Map(func(card *SetData) string {
		return fmt.Sprintf("insert into set_name_texts(id, en, kanji, kk, donetime) values (%d, '%s', '', '', 0);", card.Officialcode, card.Name)
	})
}
