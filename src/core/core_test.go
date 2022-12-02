package core

import "testing"

func TestRespEncodeSimpleString(t *testing.T) {
	// Encoding a simple string format
	encoded := RespEncode(SimpleString, "lorem ipsum")
	const expected = "+lorem ipsum\r\n"
	if string(encoded) != expected {
		t.Fatalf("excpected %s, got %b", expected, encoded)
	}
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
