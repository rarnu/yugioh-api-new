package database

import (
	"fmt"
	. "github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"golang.org/x/text/width"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ygoapi/config"
	"ygoapi/dto"
	"ygoapi/japanese"
)

type RushDuelDB struct {
	omJp *gorm.DB
	omCn *gorm.DB
}

var Rush RushDuelDB

func NewRush() {
	rushJp, errJp := gorm.Open(sqlite.Open(config.RushDuelConfig.Jp), &gorm.Config{})
	Rushcn, errCn := gorm.Open(sqlite.Open(config.RushDuelConfig.Cn), &gorm.Config{})
	if errJp != nil {
		panic(errJp)
	}
	if errCn != nil {
		panic(errCn)
	}
	logger.Info("RushDuel connected (%v, %v).", rushJp, Rushcn)
	Rush = RushDuelDB{omJp: rushJp, omCn: Rushcn}
}

func (r RushDuelDB) RushCardNameList(name ISCString, lang ISCString) ([]*dto.CardName, error) {
	tn := r.getTextTableName(lang)
	om := IfThen(lang == "jp", r.omJp, r.omCn)
	var data ISCList[*dto.CardName]
	om.Raw(fmt.Sprintf("select id, name from %s where (name like '%%%s%%' or name like '%%%s%%' or name like '%%%s%%') and (id >= 120100000 and id <= 120999999) limit 10", tn,
		width.Widen.String(string(name)), width.Narrow.String(string(name)), name)).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = ModifyName(lang, item.Id, item.Name)
		})
		return data, nil
	}
}

func (r RushDuelDB) RushOne(password int64, lang ISCString) (*dto.CardData, error) {
	var data dto.CardData
	tn := r.getTextTableName(lang)
	if lang == "jp" {
		r.omJp.Raw(fmt.Sprintf("select id, name, desc, type, atk, def, level, race, attribute, '' as 'abbr' from %s where id = ?", tn), password).First(&data)
	} else {
		r.omCn.Raw(fmt.Sprintf(`select t.id, t.name, t.desc, d.type, d.atk, d.def, d.level, d.race, d.attribute, '' as 'abbr' from %s t left join datas d on t.id = d.id where t.id=?`, tn), password).First(&data)
	}

	if data.Id == 0 {
		return nil, fmt.Errorf("card %d not found", password)
	}
	data.Desc = data.Desc.ReplaceAll("'''", "")
	data.Desc = genRDDesc(lang, data.Desc)
	return &data, nil
}

func (r RushDuelDB) RushRandom(lang ISCString) (*dto.CardData, error) {
	var id int64 = 0
	tn := r.getTextTableName(lang)
	om := IfThen(lang == "jp", r.omJp, r.omCn)
	om.Raw(fmt.Sprintf(`select id from %s where id >= 120100000 and id <= 120999999 order by RANDOM() limit 1`, tn)).First(&id)
	if id == 0 {
		return nil, fmt.Errorf("card not found")
	}
	return r.RushOne(id, lang)
}

func (r RushDuelDB) NameKanjiKana(name ISCString) ISCString {
	var kk ISCString = ""
	r.omJp.Raw(fmt.Sprintf("select kk from ja_texts where name = '%s' or name = '%s' or name = '%s'",
		name, width.Narrow.String(string(name)), width.Widen.String(string(name)))).Scan(&kk)
	return kk
}

func (r RushDuelDB) EffectKanjiKana(name ISCString) ISCString {
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
		kk := r.NameKanjiKana(tmp)
		if kk == "" {
			tmp = tmp.ReplaceAll("トークン", "").TrimSpace()
			kk = r.SetKanjiKana(tmp)
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

func (r RushDuelDB) SetKanjiKana(name ISCString) ISCString {
	var kk ISCString = ""
	r.omJp.Raw(fmt.Sprintf("select kk from ja_set_text where name = '%s' or name = '%s' or name = '%s'",
		name, width.Narrow.String(string(name)), width.Widen.String(string(name)))).Scan(&kk)
	return kk
}

func (r RushDuelDB) RdkFindCardNameList(req dto.ReqYdkFind) ([]*dto.CardName, error) {
	tn := r.getTextTableName(req.Lang)
	om := IfThen(req.Lang == "jp", r.omJp, r.omCn)
	kf := "name"
	if req.ByEffect {
		kf = "desc"
	}
	var data ISCList[*dto.CardName]
	sql := fmt.Sprintf(`select id, name from %s where %s like '%%%s%%' limit 100`, tn, kf, ToDBStr(req.Key))
	om.Raw(sql).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = ModifyName(req.Lang, item.Id, item.Name)
		})
		return data, nil
	}
}

func (r RushDuelDB) RdkNamesByIds(req dto.ReqYdkNames) ([]*dto.CardName, error) {
	tn := r.getTextTableName(req.Lang)
	om := IfThen(req.Lang == "jp", r.omJp, r.omCn)
	instr := req.Ids.JoinToStringFull(",", "", "", func(item int64) string {
		return ToString(item)
	})
	var data ISCList[*dto.CardName]
	sql := fmt.Sprintf(`select id, name from %s where id in (%s)`, tn, instr)
	om.Raw(sql).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = ModifyName(req.Lang, item.Id, item.Name)
		})
		return data, nil
	}
}

func (r RushDuelDB) getTextTableName(lang ISCString) ISCString {
	tn := ISCString("ja_texts")
	switch lang {
	case "sc":
		tn = "texts"
	}
	//switch lang {
	//case "sc":
	//	tn = "zhcn_texts"
	//case "tc":
	//	tn = "zhtw_texts"
	//case "en":
	//	tn = "texts"
	//case "kr":
	//	tn = "ko_texts"
	//case "de":
	//	tn = "de_texts"
	//case "fr":
	//	tn = "fr_texts"
	//case "it":
	//	tn = "it_texts"
	//case "es":
	//	tn = "es_texts"
	//case "th":
	//	tn = "th_texts"
	//case "vi":
	//	tn = "vi_texts"
	//}
	return tn
}
