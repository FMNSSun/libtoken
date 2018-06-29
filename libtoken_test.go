package libtoken

import "testing"

func TestAPI(t *testing.T) {
	tg1, _ := NewTokenGenerator("dummy", 4)
	tg2, _ := NewTokenGenerator("dummy", 3)

	token := Join("-", tg1, tg2)

	if token != "AAAA-AAA" {
		t.Errorf("Unexpected: %s", token)
	}
}

func BenchmarkFallback(b *testing.B) {
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		ReadBytesFallback(buf)
	}
}

func BenchmarkIt(b *testing.B) {
	buf := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		ReadBytesNoFallback(buf)
	}
}
