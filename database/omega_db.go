package database

import (
	"fmt"
	"github.com/isyscore/isc-gobase/time"
	"golang.org/x/text/width"
	"math/rand"
	"ygoapi/config"
	"ygoapi/dto"

	. "github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OmegaDB struct {
	om *gorm.DB
}

var Omega OmegaDB

func NewOmega() {
	omega, err := gorm.Open(sqlite.Open(config.SQLiteConfig.Host), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	logger.Info("OmegaDB connected (%v).", omega)
	Omega = OmegaDB{om: omega}
}

func (o OmegaDB) One(password int64, lang ISCString) (*dto.CardData, error) {
	var data dto.CardData
	tn := getTextTableName(lang)
	o.om.Raw(
		fmt.Sprintf(`select t.id, t.name, t.desc, d.type, d.atk, d.def, d.level, d.race, d.attribute, p.abbr from %s t
				left join datas d on t.id = d.id
				left join relations r on t.id = r.cardid
				left join packs p on r.packid = p.id
			where t.id = ? group by r.cardid`, tn), password).First(&data)

	if data.Id == 0 {
		return nil, fmt.Errorf("card %d not found", password)
	}
	data.Name = data.Name.ReplaceAll("&#64025;", "神")
	data.Desc = data.Desc.ReplaceAll("&#64025;", "神")
	data.Abbr = genPackName(lang, data.Abbr)
	if lang == "en" {
		data.Desc = data.Desc.ReplaceAll("'''", "")
	}
	if lang == "tc" {
		data.Name = data.Name.SubStringBeforeLast("【")
	}
	return &data, nil
}

func (o OmegaDB) Random(lang ISCString) (*dto.CardData, error) {
	var id int64 = 0
	tn := getTextTableName(lang)
	o.om.Raw(fmt.Sprintf(`select id from %s where id >= 10000 and id <= 99999999 order by RANDOM() limit 1`, tn)).First(&id)
	if id == 0 {
		return nil, fmt.Errorf("card not found")
	}
	return o.One(id, lang)
}

func (o OmegaDB) CardNameList(name ISCString, lang ISCString) ([]*dto.CardName, error) {
	tn := getTextTableName(lang)
	var data ISCList[*dto.CardName]
	o.om.Raw(fmt.Sprintf("select id, name from %s where name like '%%%s%%' or name like '%%%s%%' or name like '%%%s%%' limit 10", tn,
		width.Widen.String(string(name)), width.Narrow.String(string(name)), name)).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = item.Name.ReplaceAll("&#64025;", "神")
			if lang == "tc" {
				item.Name = item.Name.SubStringBeforeLast("【")
			}
		})
		return data, nil
	}
}

func (o OmegaDB) SearchCardList(req dto.ReqSearchOrigin) ([]*dto.CardData, error) {
	tn := getTextTableName(req.Lang)
	var list ISCList[*dto.CardData]
	sql := fmt.Sprintf(`select t.id, t.name, t.desc, d.type, d.atk, d.def, d.level, d.race, d.attribute, p.abbr from %s t 
					left join datas d on t.id = d.id
					left join relations r on t.id = r.cardid
					left join packs p on r.packid = p.id
				where 1 = 1`, tn)
	req.Key = req.Key.ReplaceAll("'", "''")
	if req.Key.TrimSpace() != "" {
		sql += fmt.Sprintf(" and (t.name like '%%%s%%' or t.desc like '%%%s%%')", ToDBStr(req.Key), ToDBStr(req.Key))
	}
	if req.CardType != 0 {
		sql += fmt.Sprintf(" and d.type & %d", req.CardType)
	}
	if req.Attribute != 0 {
		sql += fmt.Sprintf(" and d.attribute & %d", req.Attribute)
	}
	if req.Icon != 0 {
		sql += fmt.Sprintf(" and d.type & %d", req.Icon)
	}
	if req.SubType != 0 {
		sql += fmt.Sprintf(" and d.type & %d", req.SubType)
	}
	if req.Race != 0 {
		sql += fmt.Sprintf(" and d.race & %d", req.Race)
	}
	if req.MonsterType != 0 {
		sql += fmt.Sprintf(" and d.type & %d", req.MonsterType)
	}
	o.om.Raw(sql).Scan(&list)

	if len(list) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		list.ForEach(func(item *dto.CardData) {
			item.Name = item.Name.ReplaceAll("&#64025;", "神")
			item.Desc = item.Desc.ReplaceAll("&#64025;", "神")
			item.Abbr = genPackName(req.Lang, item.Abbr)
			if req.Lang == "en" {
				item.Desc = item.Desc.ReplaceAll("'''", "")
			}
			if req.Lang == "tc" {
				item.Name = item.Name.SubStringBeforeLast("【")
			}
		})
		return list, nil
	}
}

func (o OmegaDB) YdkFindCardNameList(req dto.ReqYdkFind) ([]*dto.CardName, error) {
	tn := getTextTableName(req.Lang)
	kf := "name"
	if req.ByEffect {
		kf = "desc"
	}
	var data ISCList[*dto.CardName]
	sql := fmt.Sprintf(`select id, name from %s where %s like '%%%s%%' limit 100`, tn, kf, ToDBStr(req.Key))
	o.om.Raw(sql).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = item.Name.ReplaceAll("&#64025;", "神")
			if req.Lang == "tc" {
				item.Name = item.Name.SubStringBeforeLast("【")
			}
		})
		return data, nil
	}
}

func (o OmegaDB) YdkNamesByIds(req dto.ReqYdkNames) ([]*dto.CardName, error) {
	tn := getTextTableName(req.Lang)
	instr := req.Ids.JoinToStringFull(",", "", "", func(item int64) string {
		return ToString(item)
	})
	var data ISCList[*dto.CardName]
	sql := fmt.Sprintf(`select id, name from %s where id in (%s)`, tn, instr)
	o.om.Raw(sql).Scan(&data)
	if len(data) == 0 {
		return nil, fmt.Errorf("card not found")
	} else {
		data.ForEach(func(item *dto.CardName) {
			item.Name = item.Name.ReplaceAll("&#64025;", "神")
			if req.Lang == "tc" {
				item.Name = item.Name.SubStringBeforeLast("【")
			}
		})
		return data, nil
	}
}

func (o OmegaDB) CardCount() int {
	var cnt int = 0
	o.om.Raw("select count(1) from datas").First(&cnt)
	return cnt
}

func getTextTableName(lang ISCString) ISCString {
	tn := ISCString("ja_texts")
	switch lang {
	case "sc":
		tn = "zhcn_texts"
	case "tc":
		tn = "zhtw_texts"
	case "en":
		tn = "texts"
	case "kr":
		tn = "ko_texts"
	case "de":
		tn = "de_texts"
	case "fr":
		tn = "fr_texts"
	case "it":
		tn = "it_texts"
	case "es":
		tn = "es_texts"
	case "th":
		tn = "th_texts"
	case "vi":
		tn = "vi_texts"
	}
	return tn
}

func genPackName(lang ISCString, abbr ISCString) ISCString {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100)
	if r == 0 {
		r = 1
	}
	rstr := ISCString(ToString(r))
	if len(rstr) < 2 {
		rstr = "0" + rstr
	}
	if lang.ToLower() == "ko" {
		lang = "kr"
	}
	if len([]rune(abbr)) != 4 {
		abbr = "LWCG"
	}
	return abbr.ToUpper() + "-" + lang.ToUpper() + "0" + rstr
}

func ToDBStr(str ISCString) ISCString {
	return str.ReplaceAll("'", "''").ReplaceAll("\n", "").ReplaceAll("\r", "")
}
