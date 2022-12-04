package core

import "testing"

func assertEqual(t *testing.T, a string, b []byte) {
	if a != string(b) {
		t.Fatalf("expected %s, got %s", a, b)
	}
}

func TestRespElem_EncodeSimpleString(t *testing.T) {
	// Encoding a simple string format
	encoded := RespElem{
		Type:    SimpleString,
		Content: "lorem ipsum",
	}.Encode()
	assertEqual(t, "+lorem ipsum\r\n", encoded)
}

func TestRespElem_EncodeBulkString(t *testing.T) {
	// Encoding BulkString format
	encoded := RespElem{
		Type:    BulkString,
		Content: "hello",
	}.Encode()
	assertEqual(t, "$5\r\nhello\r\n", encoded)
}

func TestRespElem_EncodeBulkStringEmpty(t *testing.T) {
	encoded := RespElem{
		Type:    BulkString,
		Content: "",
	}.Encode()
	assertEqual(t, "$0\r\n\r\n", encoded)
}

func TestRespElem_EncodeBulkStringWithReturns(t *testing.T) {
	encoded := RespElem{
		Type:    BulkString,
		Content: "hello \r\nworld!",
	}.Encode()
	assertEqual(t, "$14\r\nhello \r\nworld!\r\n", encoded)
}

func TestRespElem_EncodeArrayEmpty(t *testing.T) {
	// Encoding an empty RESP Array
	encoded := RespElem{
		Type:    Array,
		Content: []RespElem{},
	}.Encode()
	assertEqual(t, "*0\r\n", encoded)
}

func TestRespElem_EncodeArrayBulkStrings(t *testing.T) {
	encoded := RespElem{
		Type: Array,
		Content: []RespElem{
			{
				Type:    BulkString,
				Content: "hello",
			},
			{
				Type:    BulkString,
				Content: "world",
			},
		},
	}.Encode()
	assertEqual(t, "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", encoded)
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

func BenchmarkRespElem_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Encoding a simple string format
		RespElem{
			Type:    SimpleString,
			Content: "lorem ipsum",
		}.Encode()
	}
}
