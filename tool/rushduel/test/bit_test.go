package test

import "testing"

func TestBits(t *testing.T) {
	ct1 := 0x01 + 0x40

	b1 := ct1 & 0x01
	b2 := ct1 & 0x20
	t.Logf("b1 = %v, b2 = %v", b1, b2)
}
