package utils

import "testing"

func TestGetSimpleFileName(t *testing.T) {
	path := "/path/to/2026/07/18/test/20260718T1428.md"

	result := SimpleFileName(path)

	if result != "20260718T1428" {
		t.Fail()
	}
}
