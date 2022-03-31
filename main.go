package main

import (
	"github.com/gin-gonic/gin"
	"ygoapi/config"
	"ygoapi/database"
	"ygoapi/japanese"
	"ygoapi/route"

	"github.com/isyscore/isc-gobase/server"
)

func main() {

	server.InitServer()
	config.LoadDatabaseConfig()
	database.NewOmega()
	database.NewYgoName()
	japanese.NewKanjiKanaData()

	// 注册路由
	(server.Engine().(*gin.Engine)).LoadHTMLGlob("./files/*.html")
	// server.Engine().LoadHTMLGlob("./files/*.html")
	server.Engine().Static("static", "./files")
	server.RegisterRoute("/", server.HmAll, route.ApiIndex)
	server.RegisterRoute("/system/status", server.HmAll, route.ApiSystemStatus)
	server.RegisterRoute("/api/common/count", server.HmAll, route.ApiCommonCount)
	server.RegisterRoute("/api/yugioh/search", server.HmPost, route.ApiSearchCards)
	server.RegisterRoute("/api/yugioh/list", server.HmAll, route.ApiGetCardList)
	server.RegisterRoute("/api/yugioh/card/:password", server.HmAll, route.ApiGetOneCard)
	server.RegisterRoute("/api/yugioh/random", server.HmAll, route.ApiGetRandomCard)
	server.RegisterRoute("/api/ydk/find", server.HmPost, route.ApiYdkFindCard)
	server.RegisterRoute("/api/ydk/names", server.HmPost, route.ApiYdkGetNamesByIds)
	server.RegisterRoute("/api/kanjikana/name", server.HmPost, route.ApiKKCardName)
	server.RegisterRoute("/api/kanjikana/effect", server.HmPost, route.ApiKKCardEffect)
	server.RegisterRoute("/api/kanjikana/text", server.HmPost, route.ApiKKNormalText)

	// 启动服务
	server.StartServer()
}
