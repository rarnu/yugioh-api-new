package test

import (
	"testing"
	"tool/rushduel"
)

func TestHtmlParse(t *testing.T) {
	rushduel.ParseHtml("Blue-Eyes_White_Dragon_(Rush_Duel)")
	rushduel.ParseHtml("Wyrm_Excavator_the_Heavy_Cavalry_Draco")
	rushduel.ParseHtml("Acacia_the_Shadow_Flower_Priestess")
	rushduel.ParseHtml("A.I._Bear_Can")
}
