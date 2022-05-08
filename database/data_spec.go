package database

import (
	. "github.com/isyscore/isc-gobase/isc"
)

func ModifyName(lang ISCString, id int64, name ISCString) ISCString {
	name = name.ReplaceAll("&#64025;", "神")
	// here place the name replace
	if lang == "tc" {
		name = name.SubStringBeforeLast("【")
	}

	return name
}

// ModifyDesc 特殊变更描述文字
func ModifyDesc(lang ISCString, id int64, typ int64, desc ISCString) ISCString {
	desc = desc.ReplaceAll("&#64025;", "神")
	if lang == "jp" && id == 67616300 {
		// 日语版试胆竞速
		return desc.ReplaceAll("\n", "")
	}
	if lang == "en" {
		desc = desc.ReplaceAll("'''", "")
	}
	if lang == "tc" {
		isPendulum := (typ & 0x1000000) != 0
		if isPendulum {
			desc = ModifyPendulumDesc(lang, id, desc)
		}
	}
	return desc
}

// ModifyPendulumDesc 针对灵摆卡的特殊变更描述文字
func ModifyPendulumDesc(lang ISCString, id int64, desc ISCString) ISCString {
	if lang == "tc" {
		sl := NewListWithList(desc.Split("\n")).Filter(func(s ISCString) bool {
			return s.TrimSpace() != ""
		})
		sl.Delete(1)
		sl.Delete(1)
		newDesc := ISCString(sl.JoinToStringFull("\n", "", "", func(s ISCString) string {
			return string(s.ReplaceAll("\r", ""))
		}))
		return newDesc.ReplaceAll("【怪獸效果】", "【Monster Effect】")
	}
	return desc
}
