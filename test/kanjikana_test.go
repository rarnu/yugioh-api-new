package test

import (
	. "github.com/isyscore/isc-gobase/isc"
	"testing"
	"ygoapi/database"
	"ygoapi/japanese"
)

func TestRemoveKK(t *testing.T) {
	str1 := "[青眼の白龍(ブルーアイズ・ホワイト・ドラゴン)]"
	str2 := "[真(しん)][魔(ま)][獣(じゅう)] ガーゼット"
	str3 := "[烙(らく)][印(いん)]の[絆(きずな)]"
	str4 := "ピアニッシモ"

	e1 := japanese.RemoveKana(ISCString(str1))
	e2 := japanese.RemoveKana(ISCString(str2))
	e3 := japanese.RemoveKana(ISCString(str3))
	e4 := japanese.RemoveKana(ISCString(str4))

	t.Logf("%s -> %s", str1, e1)
	t.Logf("%s -> %s", str2, e2)
	t.Logf("%s -> %s", str3, e3)
	t.Logf("%s -> %s", str4, e4)

}

func TestKK(t *testing.T) {
	database.NewYgoName()
	ts := database.YgoName.LastSync()
	t.Logf("LastSync: %s", ts)

	str1 := database.YgoName.NameKanjiKana("青眼の白龍")
	t.Logf("%s", str1)
}

func TestEffectName(t *testing.T) {
	database.NewYgoName()
	japanese.NewKanjiKanaData()

	str1 := "①：このカードは自分フィールドの表側表示の「地霊使い」、「地霊使いアウス」１体と地属性モンスター１体を墓地へ送り、手札・デッキから特殊召喚できる。②：「地霊使い」と「地霊使いアウス」と墓地へ送り"
	sa := database.YgoName.EffectKanjiKana(ISCString(str1))

	t.Logf("%s", sa)
}

func TestKKMap(t *testing.T) {
	japanese.NewKanjiKanaData()
	t.Logf("%v", japanese.KanjiKanaMap)
}

func TestNormalKana(t *testing.T) {
	str1 := "高い攻撃力を誇る伝説のドラゴン。どんな相手でも粉砕、その破壊力は計り知れない。"
	sa := database.YgoName.NormalKanjiKana(ISCString(str1))
	t.Logf("%s", sa)
}
