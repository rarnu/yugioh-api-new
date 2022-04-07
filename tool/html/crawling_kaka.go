package html

import (
	"fmt"
	"github.com/isyscore/isc-gobase/encoding"
	"github.com/isyscore/isc-gobase/http"
	. "github.com/isyscore/isc-gobase/isc"
	"log"
	"tool/consts"
)

func CrawlingKaka(name string) string {
	u0, _ := encoding.UrlEncoding(name, encoding.EUCJP)
	b, _ := http.GetSimple(fmt.Sprintf(consts.WIKI_URL, u0))
	// 这里拿到的是日语编码的字符串，要改成utf8的，否则下面无法操作
	tmp, _ := encoding.StringToUTF8(string(b), encoding.EUCJP)
	nameStr := ISCString(tmp).SubStringAfter("<h2 id=\"content_1_0\">《").SubStringBefore("》")
	nameStr = nameStr.ReplaceAll("<rb>", "").ReplaceAll("</rb>", "").
		ReplaceAll("<rp>", "").ReplaceAll("</rp>", "").
		ReplaceAll("<rt>", "").ReplaceAll("</rt>", "").
		ReplaceAll("<ruby>", "[").ReplaceAll("</ruby>", "]")
	if nameStr.Contains("/") {
		nameStr = nameStr.SubStringBefore("/")
	}
	log.Printf("name = %s\n", nameStr)
	if nameStr.Contains("<?xml") {
		nameStr = ISCString(name)
	}
	return string(nameStr)
}
