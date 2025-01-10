package main

import (
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
	database.NewRush()
	japanese.NewKanjiKanaData()
	japanese.NewTranslateData()

	// 注册路由
	server.Engine().LoadHTMLGlob("./files/*.html")
	// server.Engine().LoadHTMLGlob("./files/*.html")
	server.Engine().Static("/assets", "./assets")
	server.RegisterRoute("/", server.HmAll, route.ApiIndex)
	server.RegisterRoute("/system/status", server.HmAll, route.ApiSystemStatus)

	server.RegisterRoute("/api/common/count", server.HmAll, route.ApiCommonCount)    // no use
	server.RegisterRoute("/api/common/translate", server.HmPost, route.ApiTranslate) // done

	server.RegisterRoute("/api/yugioh/search", server.HmPost, route.ApiSearchCards)       // no use
	server.RegisterRoute("/api/yugioh/list", server.HmAll, route.ApiGetCardList)          // done
	server.RegisterRoute("/api/yugioh/card/:password", server.HmAll, route.ApiGetOneCard) // done
	server.RegisterRoute("/api/yugioh/random", server.HmAll, route.ApiGetRandomCard)      // done
	// 新接口，取卡片效果明细
	server.RegisterRoute("/api/yugioh/effect/:password", server.HmAll, route.ApiGetCardEffectDetail) // no use

	server.RegisterRoute("/api/ydk/find", server.HmPost, route.ApiYdkFindCard)       // done
	server.RegisterRoute("/api/ydk/names", server.HmPost, route.ApiYdkGetNamesByIds) // done

	server.RegisterRoute("/api/rushduel/rdkfind", server.HmPost, route.ApiRdkFindCard)       // done
	server.RegisterRoute("/api/rushduel/rdknames", server.HmPost, route.ApiRdkGetNamesByIds) // done
	server.RegisterRoute("/api/rushduel/kkname", server.HmPost, route.ApiRushKKCardName)
	server.RegisterRoute("/api/rushduel/kkeffect", server.HmPost, route.ApiRushKKCardEffect)
	server.RegisterRoute("/api/rushduel/list", server.HmAll, route.ApiRushGetCardList)          // done
	server.RegisterRoute("/api/rushduel/card/:password", server.HmAll, route.ApiRushGetOneCard) // done
	server.RegisterRoute("/api/rushduel/random", server.HmAll, route.ApiRushRandomCard)         // done

	server.RegisterRoute("/api/kanjikana/name", server.HmPost, route.ApiKKCardName)     // done
	server.RegisterRoute("/api/kanjikana/effect", server.HmPost, route.ApiKKCardEffect) // done
	server.RegisterRoute("/api/kanjikana/text", server.HmPost, route.ApiKKNormalText)   // done

	// 启动服务
	server.StartServer()
}
