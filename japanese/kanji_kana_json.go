package japanese

import (
	"encoding/json"
	"github.com/isyscore/isc-gobase/file"
	. "github.com/isyscore/isc-gobase/isc"
	"regexp"
	"sort"
)

type kkSlice ISCList[ISCString]

func (k kkSlice) Len() int {
	return len(k)
}
func (k kkSlice) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}
func (k kkSlice) Less(i, j int) bool {
	return len(k[i]) > len(k[j])
}

func SortByLength(list ISCList[ISCString]) {
	sort.Sort(kkSlice(list))
}

var KanjiKanaMap = NewMap[ISCString, ISCString]()
var KanjiKanaReg *regexp.Regexp

func NewKanjiKanaData() {
	jsonPath := "./files/kanji-kana.json"
	jsonData := file.ReadFile(jsonPath)
	if err := json.Unmarshal([]byte(jsonData), &KanjiKanaMap); err != nil {
		panic(err)
	}
	keys := KanjiKanaMap.Keys()
	SortByLength(keys)
	regstr := keys.JoinToStringFull("|", "", "", func(item ISCString) string {
		return string(item)
	})
	KanjiKanaReg, _ = regexp.Compile(regstr)
}
