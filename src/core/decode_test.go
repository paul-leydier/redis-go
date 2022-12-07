package core

import "testing"

func TestRespDecodeSimpleString(t *testing.T) {
	encoded := NewEncodedRespElem([]byte("+lorem ipsum\r\n"))
	decoded, err := RespDecode(&encoded).String()
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
	encoded := NewEncodedRespElem([]byte("$5\r\nhello\r\n"))
	decoded, err := RespDecode(&encoded).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "hello"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringEmpty(t *testing.T) {
	encoded := NewEncodedRespElem([]byte("$0\r\n\r\n"))
	decoded, err := RespDecode(&encoded).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = ""
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringWithReturns(t *testing.T) {
	encoded := NewEncodedRespElem([]byte("$14\r\nhello \r\nworld!\r\n"))
	decoded, err := RespDecode(&encoded).String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "hello \r\nworld!"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func BenchmarkParseBulkString(b *testing.B) {
	encoded := NewEncodedRespElem([]byte("$14\r\nhello \r\nworld!\r\n"))
	for i := 0; i < b.N; i++ {
		parseBulkString(&encoded)
		encoded.cursor = 0
	}
}
