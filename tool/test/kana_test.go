package test

import (
	"github.com/isyscore/isc-gobase/encoding"
	. "github.com/isyscore/isc-gobase/isc"
	"tool/html"

	"testing"
)

func TestKana(t *testing.T) {
	str1, _ := encoding.UrlEncoding("《青眼の白龍》", encoding.EUCJP)
	t.Logf("%s", str1)
	str2, _ := encoding.UrlDecoding("%A1%D4%C0%C4%B4%E3%A4%CE%C7%F2%CE%B6%A1%D5", encoding.EUCJP)
	t.Logf("%s", str2)
}

func TestKanaCrawling(t *testing.T) {
	names := ISCList[string]([]string{"青眼の白龍", "調和の宝札", "スターダスト・ドラゴン／バスター", "黒衣の大賢者"})
	names.ForEach(func(name string) {
		n0 := html.CrawlingKaka(name)
		t.Logf("%s = %s", name, n0)
	})
}
