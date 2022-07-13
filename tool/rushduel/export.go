package rushduel

import (
	"encoding/json"
	"fmt"
	f0 "github.com/isyscore/isc-gobase/file"
	. "github.com/isyscore/isc-gobase/isc"
	"os"
	"path/filepath"
)

func ExportRushDuelDatabase() {
	GlobalSetNames.Clear()
	GlobalExistedCardNames.Clear()
	list := downloadJson()
	db := NewRushDuelDB()
	GlobalExistedCardIds = db.GetAllCardIds()
	GlobalExistedSetNames = db.GetAllSetNames()

	dir, _ := os.Getwd()
	list.ForEach(func(item ISCString) {
		if item.Contains("_mt_") {
			var data RushMagicTrapDTO
			jsonstr := f0.ReadFileBytes(filepath.Join(dir, string(item)))
			_ = json.Unmarshal(jsonstr, &data)
			for k, v := range data.Results {
				if v.Printouts.Status[0].FullText == "Not yet released" {
					fmt.Printf("%s not yet released\n", v.Printouts.JapaneseName[0])
					continue
				}
				db.InsertMagicTrap(k, v.Printouts)
			}
		} else {
			var data RushDuelDTO
			jsonstr := f0.ReadFileBytes(filepath.Join(dir, string(item)))
			_ = json.Unmarshal(jsonstr, &data)
			for k, v := range data.Results {
				if v.Printouts.Status[0].FullText == "Not yet released" {
					fmt.Printf("%s not yet released\n", v.Printouts.JapaneseName[0])
					continue
				}
				db.Insert(k, v.Printouts)
			}
		}
	})

	// 导入 Set Name
	//GlobalSetNames = GlobalSetNames.Distinct()
	//for _, n := range GlobalSetNames {
	//	db.InsertSet(n)
	//}

	db.Close()
}
