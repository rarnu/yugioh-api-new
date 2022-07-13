package rushduel

import (
	"fmt"
	f0 "github.com/isyscore/isc-gobase/file"
	. "github.com/isyscore/isc-gobase/isc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

type RushDuelDB struct {
	om    *gorm.DB
	omega *gorm.DB
}

func NewRushDuelDB() *RushDuelDB {
	needCreateTable := false
	dir, _ := os.Getwd()
	fileName := filepath.Join(dir, "RushDuelJP.cdb")
	if !f0.FileExists(fileName) {
		needCreateTable = true
	}
	rush, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("RushDuel connected (%v).", rush)
	omega, err := gorm.Open(sqlite.Open(filepath.Join(dir, "OmegaDB.cdb")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("Omega connected (%v).", omega)
	ret := &RushDuelDB{om: rush, omega: omega}
	if needCreateTable {
		ret.CreateTable()
	}
	return ret
}

func (r *RushDuelDB) Close() {
	d, _ := r.om.DB()
	d.Close()
}

func (r *RushDuelDB) CreateTable() {
	sql := "create table ja_texts(id integer primary key, name text, kk text, desc text, type integer, atk integer, def integer, level integer, race integer, attribute integer)"
	r.om.Exec(sql)
	sql = "create table ja_set_text(id integer primary key AUTOINCREMENT, name text, kk text)"
	r.om.Exec(sql)
}

func (r *RushDuelDB) GetAllCardIds() ISCList[ISCInt64] {
	var ids ISCList[ISCInt64]
	r.om.Raw("select id from ja_texts").Scan(&ids)
	return ids
}

func (r *RushDuelDB) GetAllSetNames() ISCList[ISCString] {
	var names ISCList[ISCString]
	r.om.Raw("select name from ja_set_text").Scan(&names)
	return names
}

func (r *RushDuelDB) GetIdByName(name ISCString) ISCInt64 {
	var id int64 = 0
	r.omega.Raw("select id from texts where name = ? and id >= 120100000 and id <= 120999999", name).First(&id)
	return ISCInt64(id)
}

func (r *RushDuelDB) Insert(name ISCString, data RushDuelPrintOut) {
	enName := data.Name[0].ReplaceAll("[[", "").ReplaceAll("]]", "").Split("|")[1]
	jpName, jpKK := parseJapaneseName(data.JapaneseName[0])
	fmt.Printf("parsing %s\n", jpName)
	id := r.GetIdByName(enName)
	if id == 0 {
		// 没有找到卡片 id，不录入数据
		fmt.Printf("%s(%s) not found.\n", jpName, enName)
		return
	}
	if GlobalExistedCardIds.Contains(id) {
		fmt.Printf("%s(%s) already exists.\n", jpName, enName)
		return
	}

	typ0 := parsePrimaryType(data.PrimaryType)
	atk0 := data.Atk[0].ToInt64()
	def0 := data.Def[0].ToInt64()
	lvl0 := data.Level[0].ToInt64()
	race0 := parseRace(data.Type[0].FullText)
	attr0 := parseAttribute(data.Attribute[0].FullText)
	eff0, max0, cn0 := ParseHtml(name)
	GlobalSetNames.AddAll(cn0...)

	if max0 {
		typ0 += MonsterSubTypeMaximum
	}

	// 写入
	tx := r.om.Exec("insert into ja_texts(id, name, kk, desc, type, atk, def, level, race, attribute) values (?,?,?,?,?,?,?,?,?,?)",
		id, jpName, jpKK, eff0, typ0, atk0, def0, lvl0, race0, attr0)
	if tx.Error != nil {
		fmt.Printf("insert card %s failed, error is %s\n", jpName, tx.Error.Error())
	} else {
		GlobalExistedCardNames.Add(jpName)
		fmt.Printf("insert card %s success.\n", jpName)
	}
}

func (r *RushDuelDB) InsertMagicTrap(name ISCString, data RushMagicTrapPrintOut) {
	enName := data.Name[0].ReplaceAll("[[", "").ReplaceAll("]]", "").Split("|")[1]
	jpName, jpKK := parseJapaneseName(data.JapaneseName[0])
	fmt.Printf("parsing %s\n", jpName)
	id := r.GetIdByName(enName)
	if id == 0 {
		// 没有找到卡片id，不录入数据
		fmt.Printf("%s(%s) not found.\n", jpName, enName)
		return
	}
	if GlobalExistedCardIds.Contains(id) {
		fmt.Printf("%s(%s) already exists.\n", jpName, enName)
		return
	}

	typ0 := parseMagicTrapType(data.Property[0])
	eff0, _, cn0 := ParseHtml(name)
	GlobalSetNames.AddAll(cn0...)

	// 写入
	tx := r.om.Exec("insert into ja_texts(id, name, kk, desc, type, atk, def, level, race, attribute) values (?,?,?,?,?,?,?,?,?,?)",
		id, jpName, jpKK, eff0, typ0, 0, 0, 0, 0, 0)
	if tx.Error != nil {
		fmt.Printf("insert card %s failed, error is %s\n", jpName, tx.Error.Error())
	} else {
		GlobalExistedCardNames.Add(jpName)
		fmt.Printf("insert card %s success.\n", jpName)
	}
}

func (r *RushDuelDB) InsertSet(name ISCString) {
	if GlobalExistedSetNames.Contains(name) {
		fmt.Printf("set %s already exists.\n", name)
		return
	}
	if GlobalExistedCardNames.Contains(name) {
		return
	}

	// 写入
	tx := r.om.Exec("insert into ja_set_text(name, kk) values (?,?)", name, name)
	if tx.Error != nil {
		fmt.Printf("insert set %s failed, error is %s\n", name, tx.Error.Error())
	} else {
		fmt.Printf("insert set %s success.\n", name)
	}
}
