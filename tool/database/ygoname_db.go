package database

import (
	"fmt"
	. "github.com/isyscore/isc-gobase/isc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type YgoNameDB struct {
	om *gorm.DB
}

type CardNameData struct {
	Id    int64
	Kanji string
	Kk    string
}

type SetNameData struct {
	Id    int64
	En    string
	Kanji string
	Kk    string
}

var YgoName YgoNameDB

func NewYgoName() {
	dsn := "root:root@tcp(127.0.0.1:3306)/YugiohAPI2?charset=utf8mb4&parseTime=True&loc=Local"
	ygo, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("YgoNameDB connected (%v).", ygo)
	YgoName = YgoNameDB{om: ygo}
}

func (o YgoNameDB) CardIdList() ISCList[int64] {
	var list ISCList[int64]
	o.om.Raw("select id from card_name_texts").Scan(&list)
	return list
}

func (o YgoNameDB) SetIdList() ISCSet[int64] {
	var list ISCList[int64]
	o.om.Raw("select id from set_name_texts group by id").Scan(&list)
	return list.ToSet()
}

func (o YgoNameDB) GetUndoneCards() ISCList[string] {
	var list ISCList[*CardNameData]
	o.om.Raw("select id, kanji, kk from card_name_texts where donetime = 0").Scan(&list)
	return ListToMapFrom[*CardNameData, string](list).Map(func(item *CardNameData) string {
		return fmt.Sprintf("insert into card_name_texts(id, kanji, kk, donetime) values (%d, '%s', '%s', 0);", item.Id, item.Kanji, item.Kk)
	})
}

func (o YgoNameDB) GetUndoneSets() ISCList[string] {
	var list ISCList[*SetNameData]
	o.om.Raw("select id, en, kanji, kk from set_name_texts where donetime = 0").Scan(&list)
	return ListToMapFrom[*SetNameData, string](list).Map(func(item *SetNameData) string {
		return fmt.Sprintf("insert into set_name_texts(id, en, kanji, kk, donetime) values (%d, '%s', '%s', '%s', 0);", item.Id, item.En, item.Kanji, item.Kk)
	})
}

func (o YgoNameDB) GetUpdateCards() ISCList[string] {
	var list ISCList[*CardNameData]
	o.om.Raw("select id, kanji, kk from card_name_texts where donetime = 0").Scan(&list)
	return ListToMapFrom[*CardNameData, string](list).Map(func(item *CardNameData) string {
		return fmt.Sprintf("update card_name_texts set kk='%s' where id = %d;", item.Kk, item.Id)
	})
}

func (o YgoNameDB) GetUpdateSets() ISCList[string] {
	var list ISCList[*SetNameData]
	o.om.Raw("select id, en, kanji, kk from set_name_texts where donetime = 0").Scan(&list)
	return ListToMapFrom[*SetNameData, string](list).Map(func(item *SetNameData) string {
		return fmt.Sprintf("update set_name_texts set kk='%s' where id = %d;", item.Kk, item.Id)
	})
}
