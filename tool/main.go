package main

import (
	"fmt"
	"github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/time"
	"log"
	"os"
	"path/filepath"
	"tool/consts"
	"tool/database"
	"tool/html"
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

func exportRemoteUpdateSQL() {
	// 连接数据库
	database.NewYgoName()
	// 获取新的卡片和SET
	updateCards := database.YgoName.GetUpdateCards()
	updateSets := database.YgoName.GetUpdateSets()

	log.Printf("update cards = %d, sets = %d", updateCards.Size(), updateSets.Size())

	if updateCards.Size() == 0 && updateSets.Size() == 0 {
		log.Printf("no update cards or sets")
		return
	}

	// SQL语句写入文件里
	updateCardSQLs := updateCards.JoinToStringFull("\n", "", "", func(s string) string { return s })
	updateSetSQLs := updateSets.JoinToStringFull("\n", "", "", func(s string) string { return s })
	file.WriteFile(fmt.Sprintf("./update_data_%s.sql", time.TimeToStringYmdHms(time.Now())), updateCardSQLs+"\n"+updateSetSQLs)

}

func downloadLastOmega() {
	err := html.DownloadFile(consts.OMEGADB_URL, filepath.Join(".", consts.OMEGADB), func(progress int64, total int64) {
		fmt.Printf("downloading %s, progress = %d/%d\r", consts.OMEGADB, progress, total)
	})
	if err != nil {
		log.Printf("download %s failed, err = %v", consts.OMEGADB, err)
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("-n		-- check for new cards and sets")
		log.Println("-e		-- export new data (donetime = 0)")
		log.Println("-u     -- export update data (donetime = 0)")
		log.Println("-d		-- download last omega")
		return
	}
	if args[1] == "-n" {
		exportLocalNewCardsSQL()
	}

	if args[1] == "-e" {
		exportRemoteUndoneSQL()
	}

	if args[1] == "-u" {
		exportRemoteUpdateSQL()
	}

	if args[1] == "-d" {
		downloadLastOmega()
	}
}
