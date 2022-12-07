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

func TestRespDecodeInteger(t *testing.T) {
	// Decoding Integer format
	encoded := NewEncodedRespElem([]byte(":1253\r\n"))
	decoded, err := RespDecode(&encoded).Int()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = 1253
	if decoded != expected {
		t.Fatalf("expected %d, got %d", expected, decoded)
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

func TestRespDecodeArrayEmpty(t *testing.T) {
	// Decoding an empty RESP Array
	encoded := NewEncodedRespElem([]byte("*0\r\n"))
	decoded, err := RespDecode(&encoded).Array()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty array, got %v", decoded)
	}
}

func TestRespDecodeArrayBulkString(t *testing.T) {
	encoded := NewEncodedRespElem([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	decoded, err := RespDecode(&encoded).Array()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	const expected = "hello"
	s, err := decoded[0].String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	if s != expected {
		t.Fatalf("expected %s, got %s", expected, s)
	}
	const expected2 = "world"
	s, err = decoded[1].String()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	if s != expected2 {
		t.Fatalf("expected %s, got %s", expected2, s)
	}
}

func TestRespDecodeArrayMixedType(t *testing.T) {
	encoded := NewEncodedRespElem([]byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n"))
	decoded, err := RespDecode(&encoded).Array()
	if err != nil {
		t.Fatalf("error decoding resp - %s", err)
	}
	expected := []RespElem{
		{Integer, 1},
		{Integer, 2},
		{Integer, 3},
		{Integer, 4},
		{BulkString, "hello"},
	}
	for i := 0; i < 5; i++ {
		if decoded[i].Type != expected[i].Type {
			t.Fatalf("expected %d, got %d", expected[i].Type, decoded[i].Type)
		}
		if decoded[i].Content != expected[i].Content {
			t.Fatalf("expected %s, got %s", expected[i].Content, decoded[i].Content)
		}
	}
}

func BenchmarkParseBulkString(b *testing.B) {
	encoded := NewEncodedRespElem([]byte("$14\r\nhello \r\nworld!\r\n"))
	for i := 0; i < b.N; i++ {
		parseBulkString(&encoded)
		encoded.cursor = 0
	}
}
