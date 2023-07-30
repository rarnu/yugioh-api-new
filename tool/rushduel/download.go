package rushduel

import (
	"encoding/json"
	"fmt"
	f0 "github.com/isyscore/isc-gobase/file"
	h0 "github.com/isyscore/isc-gobase/http"
	. "github.com/isyscore/isc-gobase/isc"
	"net/http"
	"os"
	"path/filepath"
)

const baseurl = ISCString("https://yugipedia.com/wiki/Special:Ask/mainlabel%3D-2D/format%3Djson/sort%3D/order%3Dasc/offset%3D{{offset}}/limit%3D500/-5B-5BMedium::Rush-20Duel-5D-5D-20-5B-5BCard-20type::Monster-20Card-5D-5D/-3FEnglish-20name-20(linked)%3DName/-3FJapanese-20name/-3FPrimary-20type/-3FAttribute%3D-5B-5BAttribute-5D-5D/-3FType%3D-5B-5BType-5D-5D/-3FStars-20string%3D-5B-5BLevel-5D-5D/-3FATK-20string%3D-5B-5BATK-5D-5D/-3FDEF-20string%3D-5B-5BDEF-5D-5D/-3FRush-20Duel-20status%3DStatus/prettyprint%3Dtrue/unescape%3Dtrue/searchlabel%3DJSON")
const magicurl = ISCString("https://yugipedia.com/wiki/Special:Ask/mainlabel%3D-2D/format%3Djson/sort%3D/order%3Dasc/offset%3D{{offset}}/limit%3D500/-5B-5BMedium::Rush-20Duel-5D-5D-20-5B-5BConcept:Non-2Dmonster-20cards-5D-5D/-3FEnglish-20name-20(linked)%3DName/-3FJapanese-20name/-3FCard-20type%3D-5B-5BCard-20type-5D-5D/-3FProperty%3D-5B-5BProperty-5D-5D/-3FRush-20Duel-20status%3DStatus/prettyprint%3Dtrue/unescape%3Dtrue/searchlabel%3DJSON")

func downloadJson() ISCList[ISCString] {
	var ret ISCList[ISCString]
	offset := ISCInt64(0)
	pageSize := ISCInt64(500)
	idx := 1
	dir, _ := os.Getwd()
	// 怪兽卡
	for {
		fn := ISCString(fmt.Sprintf("rushduel_%d.json", idx))
		var jsonstr any
		if !f0.FileExists(string(fn)) {
			url := baseurl.ReplaceAll("{{offset}}", string(offset.ToString()))
			_, _, jsonstr, _ = h0.Get(string(url), http.Header{"authority": {"yugipedia.com"}}, nil)
			f0.WriteFile(filepath.Join(dir, string(fn)), string(jsonstr.([]byte)))
		} else {
			jsonstr = f0.ReadFileBytes(string(fn))
		}
		ret.Add(fn)
		var data RushDuelDTO
		_ = json.Unmarshal(jsonstr.([]byte), &data)
		if data.Rows != pageSize {
			break
		}
		offset += pageSize
		idx++
	}
	// 怪兽以外
	offset = 0
	idx = 1
	for {
		fn := ISCString(fmt.Sprintf("rushduel_mt_%d.json", idx))
		var jsonstr any
		if !f0.FileExists(string(fn)) {
			url := magicurl.ReplaceAll("{{offset}}", string(offset.ToString()))
			_, _, jsonstr, _ = h0.Get(string(url), http.Header{"authority": {"yugipedia.com"}}, nil)
			f0.WriteFile(filepath.Join(dir, string(fn)), string(jsonstr.([]byte)))
		} else {
			jsonstr = f0.ReadFileBytes(string(fn))
		}
		ret.Add(fn)
		var data RushMagicTrapDTO
		_ = json.Unmarshal(jsonstr.([]byte), &data)
		if data.Rows != pageSize {
			break
		}
		offset += pageSize
		idx++
	}
	return ret
}
