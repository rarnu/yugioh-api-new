package test

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
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
