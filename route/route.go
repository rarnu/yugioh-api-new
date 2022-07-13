package route

import (
	"encoding/json"
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
	jsonstr, _ := json.Marshal(dto.RespCommonCount{
		ResponseBase: rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: dto.CommonCount{
			CardCount: cardCount,
			KanaCount: kanaCount,
			SetCount:  setCount,
		},
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
}

func ApiSearchCards(c *gin.Context) {
	var req dto.ReqSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	reqOri := dto.ReqSearchToOrigin(req)
	logger.Info("req: %v", reqOri)
	if data, err := database.Omega.SearchCardList(reqOri); err == nil {
		jsonstr, _ := json.Marshal(dto.SearchCardData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: data,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
}

func ApiGetCardList(c *gin.Context) {
	aname := ISCString(c.Query("name"))
	lang := ISCString(c.Query("lang"))
	if aname == "" {
		jsonstr, _ := json.Marshal(dto.CardNameData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: nil,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
		return
	}

	if lang == "" {
		lang = "jp"
	}
	if data, err := database.Omega.CardNameList(aname, lang); err == nil {
		jsonstr, _ := json.Marshal(dto.CardNameData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: data,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: err.Error(),
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "card id must be an integer.",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if data, err := database.Omega.One(apassInt, lang); err == nil {
		jsonstr, _ := json.Marshal(data)
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
	}
}

func ApiGetRandomCard(c *gin.Context) {
	lang := ISCString(c.Query("lang"))
	if lang == "" {
		lang = "jp"
	}
	if data, err := database.Omega.Random(lang); err == nil {
		jsonstr, _ := json.Marshal(data)
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
	}
}

func ApiYdkFindCard(c *gin.Context) {
	var req dto.ReqYdkFind
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Key == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "no key",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if data, err := database.Omega.YdkFindCardNameList(req); err == nil {
		jsonstr, _ := json.Marshal(dto.CardNameData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: data,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
	}
}

func ApiRdkFindCard(c *gin.Context) {
	var req dto.ReqYdkFind
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Key == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "no key",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	jsonstr, _ := json.Marshal(dto.CardNameData{
		ResponseBase: rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: data,
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
}

func ApiYdkGetNamesByIds(c *gin.Context) {
	var req dto.ReqYdkNames
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if len(req.Ids) == 0 {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "no card id",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	if data, err := database.Omega.YdkNamesByIds(req); err == nil {
		jsonstr, _ := json.Marshal(dto.CardNameData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: data,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
	}
}

func ApiRdkGetNamesByIds(c *gin.Context) {
	var req dto.ReqYdkNames
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Lang == "" {
		req.Lang = "jp"
	}
	if len(req.Ids) == 0 {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "no card id",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	jsonstr, _ := json.Marshal(dto.CardNameData{
		ResponseBase: rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: data,
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
}

func ApiKKCardName(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	if req.Name == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "Name is empty",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	aname := database.YgoName.NameKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		jsonstr, _ := json.Marshal(dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "found",
			},
			Data: aname,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "not found",
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	}
}

func ApiRushKKCardName(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Name == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "Name is empty",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	aname := database.Rush.NameKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		jsonstr, _ := json.Marshal(dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "found",
			},
			Data: aname,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "not found",
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	}
}

func ApiKKCardEffect(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	if req.Name == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "Name is empty",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	aname := database.YgoName.EffectKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		jsonstr, _ := json.Marshal(dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "found",
			},
			Data: aname,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "not found",
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	}
}

func ApiRushKKCardEffect(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	if req.Name == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "Name is empty",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	aname := database.Rush.EffectKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		jsonstr, _ := json.Marshal(dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "found",
			},
			Data: aname,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "not found",
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	}
}

func ApiKKNormalText(c *gin.Context) {
	var req dto.ReqKKName
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "invalid request params",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	if req.Name == "" {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "Name is empty",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	aname := database.YgoName.NormalKanjiKana(japanese.RemoveKana(req.Name))
	if aname != "" {
		jsonstr, _ := json.Marshal(dto.RespCommonString{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "found",
			},
			Data: aname,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
	} else {
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "not found",
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: "card id must be an integer.",
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}

	jsonstr, _ := json.Marshal(rsp.DataResponse[*dto.CardData]{
		ResponseBase: rsp.ResponseBase{
			Code:    0,
			Message: "",
		},
		Data: data,
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
}

func ApiRushGetCardList(c *gin.Context) {
	aname := ISCString(c.Query("name"))
	lang := ISCString(c.Query("lang"))
	if aname == "" {
		jsonstr, _ := json.Marshal(dto.CardNameData{
			ResponseBase: rsp.ResponseBase{
				Code:    http.StatusOK,
				Message: "",
			},
			Data: nil,
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: err.Error(),
		})
		c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
		return
	}

	jsonstr, _ := json.Marshal(dto.CardNameData{
		ResponseBase: rsp.ResponseBase{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: data,
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
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
		jsonstr, _ := json.Marshal(rsp.ResponseBase{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Data(http.StatusInternalServerError, h0.ContentTypeJson, jsonstr)
		return
	}
	jsonstr, _ := json.Marshal(rsp.DataResponse[*dto.CardData]{
		ResponseBase: rsp.ResponseBase{
			Code:    0,
			Message: "",
		},
		Data: data,
	})
	c.Data(http.StatusOK, h0.ContentTypeJson, jsonstr)
}
