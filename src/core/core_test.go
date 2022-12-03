package core

import "testing"

func assertEqual(t *testing.T, a string, b []byte) {
	if a != string(b) {
		t.Fatalf("expected %s, got %s", a, b)
	}
}

func TestRespEncodeSimpleString(t *testing.T) {
	// Encoding a simple string format
	encoded := RespEncode(SimpleString, "lorem ipsum")
	assertEqual(t, "+lorem ipsum\r\n", encoded)
}

func TestRespEncodeBulkString(t *testing.T) {
	// Encoding BulkString format
	encoded := RespEncode(BulkString, "hello")
	assertEqual(t, "$5\r\nhello\r\n", encoded)
}

func TestRespEncodeBulkStringEmpty(t *testing.T) {
	encoded := RespEncode(BulkString, "")
	assertEqual(t, "$0\r\n\r\n", encoded)
}

func TestRespEncodeBulkStringWithReturns(t *testing.T) {
	encoded := RespEncode(BulkString, "hello \r\nworld!")
	assertEqual(t, "$14\r\nhello \r\nworld!\r\n", encoded)
}

func TestRespDecodeSimpleString(t *testing.T) {
	respType, decoded := RespDecode([]byte("+lorem ipsum\r\n"))
	if respType != SimpleString {
		t.Fatalf("invalid response type: expected %d, got %d", SimpleString, respType)
	}
	const expected = "lorem ipsum"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkString(t *testing.T) {
	// Decoding BulkString format
	respType, decoded := RespDecode([]byte("$5\r\nhello\r\n"))
	if respType != BulkString {
		t.Fatalf("invalid response type: expected %d, got %d", BulkString, respType)
	}
	const expected = "hello"
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringEmpty(t *testing.T) {
	respType, decoded := RespDecode([]byte("$0\r\n\r\n"))
	if respType != BulkString {
		t.Fatalf("invalid response type: expected %d, got %d", BulkString, respType)
	}
	const expected = ""
	if decoded != expected {
		t.Fatalf("expected %s, got %s", expected, decoded)
	}
}

func TestRespDecodeBulkStringWithReturns(t *testing.T) {
	respType, decoded := RespDecode([]byte("$14\r\nhello \r\nworld!\r\n"))
	if respType != BulkString {
		t.Fatalf("invalid response type: expected %d, got %d", BulkString, respType)
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
