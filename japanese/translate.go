package japanese

import (
	"encoding/json"
	"fmt"
	"github.com/isyscore/isc-gobase/coder"
	"github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/goid"
	h0 "github.com/isyscore/isc-gobase/http"
	. "github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/time"
	"strconv"
)

const (
	YoudaoUrl = "https://openapi.youdao.com/api"
	AppKey    = "4a163f5a6c14ce0e"
	AppSecret = "cv5jn949S9nekKm2Ozbz3tqVh2yiYjsT"
)

type TransQuery struct {
	Q        string `json:"q"`
	From     string `json:"from"`
	To       string `json:"to"`
	AppKey   string `json:"appKey"`
	Salt     string `json:"salt"`
	Sign     string `json:"sign"`
	SignType string `json:"signType"`
	Curtime  string `json:"curtime"`
}

type TransResult struct {
	ErrorCode   string             `json:"errorCode"`
	Query       string             `json:"query"`
	Translation ISCList[ISCString] `json:"translation"`
}

func Translate(content ISCString) ISCString {
	salt := goid.GenerateUUID()
	m := map[string]string{
		"q":        string(content),
		"from":     "zh-CHS",
		"to":       "ja",
		"appKey":   AppKey,
		"salt":     salt,
		"curtime":  fmt.Sprintf("%d", time.Now().Unix()),
		"signType": "v3",
		"sign":     sign(ISCList[rune](content), salt),
	}
	ret, err := h0.PostForm(YoudaoUrl, nil, m)
	if err != nil {
		return ""
	}
	var data TransResult
	err = json.Unmarshal(ret.([]byte), &data)
	if err != nil {
		return ""
	}
	if data.ErrorCode != "0" {
		return ""
	} else {
		str := ISCString(data.Translation.JoinToStringFull("", "", "", func(it ISCString) string {
			return string(it)
		}))
		TranslateMap.ForEach(func(key ISCString, value ISCString) {
			str = str.ReplaceAll(string(key), string(value))
		})
		return str
	}
}

func sign(content ISCList[rune], salt string) string {
	input := ""
	if content.Size() > 20 {
		// q前10个字符 + q长度 + q后10个字符
		input = string(content[:10]) + strconv.Itoa(content.Size()) + string(content[content.Size()-10:])
	} else {
		input = string(content)
	}
	// sha256(应用ID+input+salt+curtime+应用密钥)
	str := fmt.Sprintf("%s%s%s%d%s", AppKey, input, salt, time.Now().Unix(), AppSecret)
	return coder.Sha256String(str)
}

var TranslateMap = NewMap[ISCString, ISCString]()

func NewTranslateData() {
	jsonPath := "./files/trans-text.json"
	jsonData := file.ReadFile(jsonPath)
	if err := json.Unmarshal([]byte(jsonData), &TranslateMap); err != nil {
		panic(err)
	}
}
