package util

import (
	. "github.com/isyscore/isc-gobase/isc"
)

func CalcNewIds(omegaList ISCList[int64], nameList ISCList[int64]) ISCList[int64] {
	m := make(map[int64]int64)
	var list ISCList[int64]
	for _, e := range nameList {
		m[e] = e
	}
	for _, e := range omegaList {
		if _, ok := m[e]; !ok {
			list.Add(e)
		}
	}
	return list
}
