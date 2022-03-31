package japanese

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/isyscore/isc-gobase/coder"
	. "github.com/isyscore/isc-gobase/isc"
	"os/exec"
	"regexp"
)

func RemoveKana(str ISCString) ISCString {
	re, _ := regexp.Compile("\\[.*?\\(.*?\\)]")
	dest := re.ReplaceAllStringFunc(string(str), func(s string) string {
		ss := ISCString(s)
		ss = ss.ReplaceAll("[", "")
		ss = ss.SubStringBefore("(")
		return string(ss)
	})
	return ISCString(dest)
}

func EffectCardNames(str ISCString) ISCList[ISCString] {
	re, _ := regexp.Compile("「.*?」")
	sa := ListToMapFrom[string, ISCString](re.FindAllString(string(str), -1))
	return sa.Map(func(item string) ISCString {
		return ISCString(item).ReplaceAll("「", "").ReplaceAll("」", "")
	}).Distinct()
}

func Kana(str ISCString) ISCString {
	re, _ := regexp.Compile("\\[.*?\\(.*?\\)]")
	sret := re.ReplaceAllStringFunc(string(str), func(s string) string {
		return fmt.Sprintf("|%s|", s)
	})
	ret := ISCList[ISCString](ISCString(sret).Split("|")).JoinToStringFull("", "", "", func(item ISCString) string {
		reg, _ := regexp.Compile("\\[.*?\\(.*?\\)]")
		if !reg.MatchString(string(item)) {
			return KanjiKanaReg.ReplaceAllStringFunc(string(item), func(ss string) string {
				if v, ok := KanjiKanaMap[ISCString(ss)]; ok {
					return string(v)
				} else {
					return ""
				}
			})
		} else {
			return string(item)
		}
	})
	return ISCString(ret)
}

type kkResult struct {
	Kana string `json:"kana"`
}

func NormalKana(str ISCString) ISCString {
	jarPath := "./files/DIYKana.jar"
	b64 := coder.Base64Encrypt(string(str))
	cmd := exec.Command("java", "-jar", jarPath, b64)
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	var kk kkResult
	_ = json.Unmarshal(out.Bytes(), &kk)
	return ISCString(coder.Base64Decrypt(kk.Kana))
}

/*
func CharacterToHalf(str ISCString) ISCString {
	ret := ""
	for _, c := range str {
		if c == '　' {
			ret += " "
		} else if c == '﹒' {
			ret += "·"
		} else if (c == '＠' || c == '．' || c == '＆' || c == '？' || c == '！') || (c >= 65313 && c <= 65338) || (c >= 65338 && c <= 65370) {
			ret += string(c - 65248)
		} else {
			ret += string(c)
		}
	}
	re, _ := regexp.Compile("「.*?」")
	ret = re.ReplaceAllStringFunc(ret, func(s string) string {
		return string(NumberToHalf(ISCString(s)))
	})
	return ISCString(ret)
}

func NumberToHalf(str ISCString) ISCString {
	ret := ""
	for _, c := range str {
		if c >= 65296 && c <= 65305 {
			ret += string(c - 65248)
		} else {
			ret += string(c)
		}
	}
	return ISCString(ret)
}

func ToDBC(str ISCString) ISCString {
	s := CharacterToHalf(str)
	s = NumberToHalf(s)
	s = s.ReplaceAll("Ɐ", "∀").ReplaceAll("´", "’")
	return s
}

*/
