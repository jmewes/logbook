package utils

import "testing"

func TestToExtendedFormat(t *testing.T) {
	datetime := "20260718T1428"

	result := ToExtendedFormat(datetime)

	if result != "2026-07-18T14:28" {
		t.Fail()
	}
}
