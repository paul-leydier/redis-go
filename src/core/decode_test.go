package core

import "testing"

func TestRespDecodeSimpleString(t *testing.T) {
	decoded, err := RespDecode([]byte("+lorem ipsum\r\n")).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "lorem ipsum"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkString(t *testing.T) {
	// Decoding BulkString format
	decoded, err := RespDecode([]byte("$5\r\nhello\r\n")).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "hello"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringEmpty(t *testing.T) {
	decoded, err := RespDecode([]byte("$0\r\n\r\n")).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = ""
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringWithReturns(t *testing.T) {
	decoded, err := RespDecode([]byte("$14\r\nhello \r\nworld!\r\n")).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "hello \r\nworld!"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func BenchmarkParseBulkStringNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseBulkStringNaive([]byte("$14\r\nhello \r\nworld!\r\n"))
	}
}

func BenchmarkParseBulkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseBulkString([]byte("$14\r\nhello \r\nworld!\r\n"))
	}
}
