package route

import (
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server/rsp"
	"net/http"
	"time"
	"ygoapi/database"
	"ygoapi/dto"
	"ygoapi/japanese"

	"github.com/gin-gonic/gin"
	h0 "github.com/isyscore/isc-gobase/http"
	. "github.com/isyscore/isc-gobase/isc"
	t0 "github.com/isyscore/isc-gobase/time"
)

func ApiIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ApiSystemStatus(c *gin.Context) {
	c.Data(http.StatusOK, h0.ContentTypeText, []byte(t0.TimeToStringYmdHms(time.Now())))
}

func ApiCommonCount(c *gin.Context) {
	cardCount := database.Omega.CardCount()
	kanaCount := database.YgoName.KanaCount()
	setCount := database.YgoName.SetCount()
	c.JSON(http.StatusOK, dto.RespCommonCount{
		ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
		Data: dto.CommonCount{
			CardCount: cardCount,
			KanaCount: kanaCount,
			SetCount:  setCount,
		},
	})
}

func ApiSearchCards(c *gin.Context) {
	var req dto.ReqSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	reqOri := dto.ReqSearchToOrigin(req)
	logger.Info("req: %v", reqOri)
	if data, err := database.Omega.SearchCardList(reqOri); err == nil {
		c.JSON(http.StatusOK, dto.SearchCardData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         data,
		})
	} else {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
}

func ApiGetCardList(c *gin.Context) {
	aname := ISCString(c.Query("name"))
	lang := ISCString(c.Query("lang"))
	if aname == "" {
		c.JSON(http.StatusOK, dto.CardNameData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         nil,
		})
		return
	}

	if lang == "" {
		lang = "jp"
	}
	if data, err := database.Omega.CardNameList(aname, lang); err == nil {
		c.JSON(http.StatusOK, dto.CardNameData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         data,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: err.Error()})
	}
}

func ApiGetOneCard(c *gin.Context) {
	apass := ISCString(c.Param("password"))
	lang := ISCString(c.Query("lang"))
	if lang == "" {
		lang = "jp"
	}
	apassInt := ToInt64(apass)
	if apassInt == 0 {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "card id must be an integer."})
		return
	}
	if data, err := database.Omega.One(apassInt, lang); err == nil {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}

func ApiGetRandomCard(c *gin.Context) {
	lang := ISCString(c.Query("lang"))
	if lang == "" {
		lang = "jp"
	}
	if data, err := database.Omega.Random(lang); err == nil {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}

func ApiYdkFindCard(c *gin.Context) {
	var req dto.ReqYdkFind
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Key == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "no key"})
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if data, err := database.Omega.YdkFindCardNameList(req); err == nil {
		c.JSON(http.StatusOK, dto.CardNameData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         data,
		})
	} else {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}

func ApiRdkFindCard(c *gin.Context) {
	var req dto.ReqYdkFind
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Key == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "no key"})
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	var data []*dto.CardName
	var err error = nil
	if req.Lang == "en" {
		data, err = database.Omega.YdkFindCardNameList(req)
	} else {
		data, err = database.Rush.RdkFindCardNameList(req)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.CardNameData{
		ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
		Data:         data,
	})
}

func ApiYdkGetNamesByIds(c *gin.Context) {
	var req dto.ReqYdkNames
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if len(req.Ids) == 0 {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "no card id"})
		return
	}

	if data, err := database.Omega.YdkNamesByIds(req); err == nil {
		c.JSON(http.StatusOK, dto.CardNameData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         data,
		})
	} else {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}

func ApiRdkGetNamesByIds(c *gin.Context) {
	var req dto.ReqYdkNames
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if len(req.Ids) == 0 {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "no card id"})
		return
	}
	var data []*dto.CardName
	var err error = nil
	if req.Lang == "en" {
		data, err = database.Omega.YdkNamesByIds(req)
	} else {
		data, err = database.Rush.RdkNamesByIds(req)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.CardNameData{
		ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
		Data:         data,
	})
}

func ApiKKCardName(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "Name is empty"})
		return
	}

	aname := database.YgoName.NameKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		c.JSON(http.StatusOK, dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: "found"},
			Data:         aname,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: "not found"})
	}
}

func ApiRushKKCardName(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "Name is empty"})
		return
	}
	aname := database.Rush.NameKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		c.JSON(http.StatusOK, dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: "found"},
			Data:         aname,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: "not found"})
	}
}

func ApiKKCardEffect(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "Name is empty"})
		return
	}

	aname := database.YgoName.EffectKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		c.JSON(http.StatusOK, dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: "found"},
			Data:         aname,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: "not found"})
	}
}

func ApiRushKKCardEffect(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "Name is empty"})
		return
	}
	aname := database.Rush.EffectKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		c.JSON(http.StatusOK, dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: "found"},
			Data:         aname,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: "not found"})
	}
}

func ApiKKNormalText(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "invalid request params"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "Name is empty"})
		return
	}

	aname := database.YgoName.NormalKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		c.JSON(http.StatusOK, dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: "found"},
			Data:         aname,
		})
	} else {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: "not found"})
	}
}

func ApiRushGetOneCard(c *gin.Context) {
	apass := ISCString(c.Param("password"))
	lang := ISCString(c.Query("lang"))
	if lang == "" {
		lang = "jp"
	}
	apassInt := ToInt64(apass)
	if apassInt == 0 {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: "card id must be an integer."})
		return
	}
	var data *dto.CardData
	var err error = nil
	if lang == "sc" || lang == "jp" {
		data, err = database.Rush.RushOne(apassInt, lang)
	} else {
		data, err = database.Omega.RushOne(apassInt, lang)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp.DataResponse[*dto.CardData]{
		ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
		Data:         data,
	})
}

func ApiRushGetCardList(c *gin.Context) {
	aname := ISCString(c.Query("name"))
	lang := ISCString(c.Query("lang"))
	if aname == "" {
		c.JSON(http.StatusOK, dto.CardNameData{
			ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
			Data:         nil,
		})
		return
	}

	if lang == "" {
		lang = "jp"
	}
	var data []*dto.CardName
	var err error = nil
	if lang == "sc" || lang == "jp" {
		data, err = database.Rush.RushCardNameList(aname, lang)
	} else {
		data, err = database.Omega.RushCardNameList(aname, lang)
	}
	if err != nil {
		c.JSON(http.StatusOK, rsp.ResponseBase{Code: http.StatusOK, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.CardNameData{
		ResponseBase: rsp.ResponseBase{Code: http.StatusOK, Message: ""},
		Data:         data,
	})
}

func ApiRushRandomCard(c *gin.Context) {
	lang := ISCString(c.Query("lang"))
	if lang == "" {
		lang = "jp"
	}
	var data *dto.CardData
	var err error = nil
	if lang == "sc" || lang == "jp" {
		data, err = database.Rush.RushRandom(lang)
	} else {
		data, err = database.Omega.RushRandom(lang)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, rsp.ResponseBase{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp.DataResponse[*dto.CardData]{
		ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
		Data:         data,
	})
}

func ApiTranslate(c *gin.Context) {
	var req dto.ReqTranslate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, rsp.ResponseBase{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	// 如果没有要翻译的内容，直接返回，不再占用有道API的请求次数
	if req.Query == "" {
		c.JSON(http.StatusOK, dto.RespTranslate{
			ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
			Data:         "",
		})
		return
	}
	retText := japanese.Translate(req.Query)
	// 翻译得到空结果
	if retText == "" {
		c.JSON(http.StatusOK, dto.RespTranslate{
			ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
			Data:         "",
		})
		return
	}
	if !req.KK {
		c.JSON(http.StatusOK, dto.RespTranslate{
			ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
			Data:         retText,
		})
		return
	}
	retTextKK := ""
	if req.KKMode == "normal" {
		retTextKK = string(database.YgoName.NormalKanjiKana(retText))
	} else {
		retTextKK = string(database.YgoName.EffectKanjiKana(retText))
	}
	if retTextKK != "" {
		c.JSON(http.StatusOK, dto.RespTranslate{
			ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
			Data:         ISCString(retTextKK),
		})
	} else {
		c.JSON(http.StatusOK, dto.RespTranslate{
			ResponseBase: rsp.ResponseBase{Code: 0, Message: ""},
			Data:         retText,
		})
	}
}
