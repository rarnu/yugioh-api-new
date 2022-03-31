package main

import (
	"fmt"
	"github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/time"
	"log"
	"os"
	"tool/database"
	"tool/util"
)

func exportLocalNewCardsSQL() {
	// 连接数据库
	database.NewOmega()
	database.NewYgoName()

	// 查最新的卡片和SET总集合
	omegaCardList := database.Omega.CardIdList()
	omegaSetList := database.Omega.SetIdList()
	log.Printf("omega card count = %d, set count = %d", omegaCardList.Size(), omegaSetList.Size())

	// 查现有的卡片和SET总集合
	nameCardList := database.YgoName.CardIdList()
	nameSetList := database.YgoName.SetIdList()
	log.Printf("name card count = %d, set count = %d", nameCardList.Size(), nameSetList.Size())

	// 计算最新和现有的差异
	diffCards := util.CalcNewIds(omegaCardList, nameCardList)
	diffSets := util.CalcNewIds(omegaSetList.ToList(), nameSetList.ToList())
	diffSets.ForEach(func(id int64) {
		log.Printf("new set id = %d", id)
	})
	log.Printf("new cards = %d, new sets = %d", diffCards.Size(), diffSets.Size())

	if diffCards.Size() == 0 && diffSets.Size() == 0 {
		log.Printf("no new cards or sets")
		return
	}

	// 根据差异获取新数据的SQL语句
	newCards := database.Omega.GetCardSQLByIds(diffCards)
	newSets := database.Omega.GetSetSQLByIds(diffSets)

	// SQL语句存到文件里
	newCardSQLs := newCards.JoinToStringFull("\n", "", "", func(s string) string { return s })
	newSetSQLs := newSets.JoinToStringFull("\n", "", "", func(s string) string { return s })
	file.WriteFile(fmt.Sprintf("./local_data_%s.sql", time.TimeToStringYmdHms(time.Now())), newCardSQLs+"\n"+newSetSQLs)

}

func exportRemoteUndoneSQL() {
	// 连接数据库
	database.NewYgoName()

	// 获取新的卡片和SET
	undoneCards := database.YgoName.GetUndoneCards()
	undoneSets := database.YgoName.GetUndoneSets()
	log.Printf("undone cards = %d, sets = %d", undoneCards.Size(), undoneSets.Size())

	if undoneCards.Size() == 0 && undoneSets.Size() == 0 {
		log.Printf("no undone cards or sets")
		return
	}

	// SQL语句写入文件里
	undoneCardSQLs := undoneCards.JoinToStringFull("\n", "", "", func(s string) string { return s })
	undoneSetSQLs := undoneSets.JoinToStringFull("\n", "", "", func(s string) string { return s })
	file.WriteFile(fmt.Sprintf("./remote_data_%s.sql", time.TimeToStringYmdHms(time.Now())), undoneCardSQLs+"\n"+undoneSetSQLs)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("-n   -- check for new cards and sets")
		log.Println("-e   -- export new data (donetime = 0)")
		return
	}
	if args[1] == "-n" {
		exportLocalNewCardsSQL()
	}

	if args[1] == "-e" {
		exportRemoteUndoneSQL()
	}
}
