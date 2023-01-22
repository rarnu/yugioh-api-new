package test

import (
	"fmt"
	. "github.com/isyscore/isc-gobase/isc"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"ygoapi/japanese"
)

func TestText(t *testing.T) {
	str := fmt.Sprintf("%%%s%%", "hello")
	t.Logf("str: %s", str)
}

func TestCharSet(t *testing.T) {
	s := `。，（）-1！@234567890abc１２３４５６７８９ａｂｃ`
	numConv := unicode.SpecialCase{
		unicode.CaseRange{
			Lo: 0x3002, // Lo 全角句号
			Hi: 0x3002, // Hi 全角句号
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x002e - 0x3002, // LowerCase 转成半角句号
				0,               // TitleCase
			},
		},
		//
		unicode.CaseRange{
			Lo: 0xFF01, // 从全角！
			Hi: 0xFF19, // 到全角 9
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0021 - 0xFF01, // LowerCase 转成半角
				0,               // TitleCase
			},
		},
	}

	fmt.Println(strings.ToLowerSpecial(numConv, s))
	fmt.Println(strings.ToUpperSpecial(numConv, s))
}

func TestInput(t *testing.T) {
	content := ISCList[rune]("abcdefghijklmnopqrstuvwxyz1234567890")
	input := string(content[:10]) + strconv.Itoa(content.Size()) + string(content[content.Size()-10:])
	t.Logf(input)
}

func TestTranslate(t *testing.T) {
	text := japanese.Translate("测试一下")
	t.Logf("text: %s", text)
}
