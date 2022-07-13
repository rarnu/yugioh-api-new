package rushduel

import (
	"fmt"
	h0 "github.com/isyscore/isc-gobase/http"
	. "github.com/isyscore/isc-gobase/isc"
	"math/rand"
	"time"
)

const pageurl = "https://yugipedia.com/wiki/"

func ParseHtml(name ISCString) (ISCString, bool, ISCList[ISCString]) {
	urlName := name.ReplaceAll(" ", "_")
	url := pageurl + urlName
	bHtml, _ := h0.GetSimple(string(url))
	htmlStr := ISCString(bHtml.([]byte))
	isLegend := htmlStr.Contains("<dd><a href=\"/wiki/Legend_Card\" title=\"Legend Card\">Legend Card</a></dd>")
	isMaximum := htmlStr.Contains("<dd><a href=\"/wiki/Requires_Maximum_Mode\" title=\"Requires Maximum Mode\">Requires Maximum Mode</a></dd>")
	htmlLines := NewListWithList(htmlStr.Split("\n"))
	pkgName := parsePackage(htmlLines)
	var maxAtk ISCInt64 = 0
	if isMaximum {
		maxAtk = parseMaximumAtk(htmlLines)
	}
	effect := parseEffect(htmlLines)
	// 组织效果文本
	retStr := pkgName
	if isLegend {
		retStr += "\t《传说卡》"
	}
	retStr += "\n"
	if isMaximum && maxAtk != 0 {
		retStr += ISCString(fmt.Sprintf("极大攻击 %d\n", maxAtk))
	}
	retStr += effect
	retStr = removeKana(retStr)
	cNames := effectCardNames(retStr)
	return retStr, isMaximum, cNames
}

/*
RD/KP04-JP022	《传说卡》
极大攻击 3500
可以和「幻龙重骑 超斗轮挖掘鳞虫［L］」「幻龙重骑 超斗轮挖掘鳞虫［R］」集齐作极大召唤。

【条件】
极大模式

【永续效果】
这张卡不会被对方的陷阱卡的效果破坏。
*/

func parsePackage(lines ISCList[ISCString]) ISCString {
	pkgLines := lines.Filter(func(item ISCString) bool {
		return item.TrimSpace().StartsWith("<style data-mw-deduplicate=\"TemplateStyles:r4444804\">")
	})
	if pkgLines.IsEmpty() {
		// 如果找不到卡包，生成 LWCG 卡包
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(100)
		if r == 0 {
			r = 1
		}
		rstr := ISCString(ToString(r))
		if len(rstr) < 2 {
			rstr = "0" + rstr
		}
		return "RD/LWCG-JP0" + rstr
	}

	pkg := pkgLines[0].SubStringAfter("title=\"").SubStringBefore("\"")
	return pkg
}

func parseMaximumAtk(lines ISCList[ISCString]) ISCInt64 {
	lIdx := 0
	for idx, line := range lines {
		if line.Contains("<a href=\"/wiki/MAXIMUM_ATK\" title=\"MAXIMUM ATK\">MAXIMUM ATK</a>") {
			lIdx = idx + 2
			break
		}
	}
	line := lines[lIdx]
	atk := line.SubStringBefore("<").TrimSpace()
	if atk == "" {
		return 0
	}
	return ISCInt64(atk.ToInt64())
}

func parseEffect(lines ISCList[ISCString]) ISCString {
	lIdx := 0
	for idx, line := range lines {
		if line.Contains("<th scope=\"row\" rowspan=") && line.Contains(">Japanese</th>") {
			lIdx = idx + 2
			break
		}
	}
	var effectLines ISCList[ISCString]
	for {
		if lines[lIdx].TrimSpace().StartsWith("</td>") {
			break
		}
		effectLines.Add(lines[lIdx])
		lIdx++
	}
	mLine := ListToMapFrom[ISCString, ISCString](effectLines).Map(func(item ISCString) ISCString {
		tmp := item.ReplaceAll("<span lang=\"ja\">", "").ReplaceAll("</span>", "")
		tmp = tmp.ReplaceAll("<ruby lang=\"ja\">", "").ReplaceAll("</ruby>", "")
		tmp = tmp.ReplaceAll("<rp>（</rp>", "(").ReplaceAll("<rp>）</rp>", ")]")
		tmp = tmp.ReplaceAll("<rt>", "").ReplaceAll("</rt>", "")
		tmp = tmp.ReplaceAll("<rb>", "[").ReplaceAll("</rb>", "")
		tmp = tmp.ReplaceAll("<dt>", "").ReplaceAll("</dt>", "")
		tmp = tmp.ReplaceAll("<p>", "").ReplaceAll("</p>", "\n\n")
		tmp = tmp.ReplaceAll("<dl>", "").ReplaceAll("</dl>", "")
		tmp = tmp.ReplaceAll("<dd>", "").ReplaceAll("</dd>", "\n\n")
		tmp = tmp.ReplaceAll("】", "】\n")
		return tmp
	}).JoinToStringFull("", "", "", func(item ISCString) string {
		return string(item)
	})

	mLine = ListToMapFrom[ISCString, ISCString](ISCString(mLine).Split("\n")).Map(func(item ISCString) ISCString {
		return item.TrimSpace()
	}).JoinToStringFull("\n", "", "", func(item ISCString) string {
		return string(item)
	})
	mLine = string(ISCString(mLine).TrimSpace())
	return ISCString(mLine)
}
