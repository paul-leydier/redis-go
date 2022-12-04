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

func TestRespElem_EncodeError(t *testing.T) {
	// Encoding a RESP error
	encoded := RespElem{
		Type:    Error,
		Content: "Error message",
	}.Encode()
	assertEqual(t, "-Error message\r\n", encoded)
}

func TestRespElem_EncodeInteger(t *testing.T) {
	// Encoding an Integer format
	encoded := RespElem{
		Type:    Integer,
		Content: 123,
	}.Encode()
	assertEqual(t, ":123\r\n", encoded)
}

func TestRespElem_EncodeIntegerNegative(t *testing.T) {
	// Encoding an Integer format
	encoded := RespElem{
		Type:    Integer,
		Content: -3,
	}.Encode()
	assertEqual(t, ":-3\r\n", encoded)
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

func TestRespElem_EncodeArrayMixedType(t *testing.T) {
	encoded := RespElem{
		Type: Array,
		Content: []RespElem{
			{Integer, 1},
			{Integer, 2},
			{Integer, 3},
			{Integer, 4},
			{BulkString, "hello"},
		},
	}.Encode()
	assertEqual(t, "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n", encoded)
}

func TestRespElem_EncodeUnknownType(t *testing.T) {
	// Trying to encode a non-existing RespType should panic
	defer func() { // assert panic
		if r := recover(); r == nil {
			t.Fatalf("did not panic")
		}
	}()
	RespElem{
		Type:    123456789, // non-existing RespType
		Content: "toto",
	}.Encode()
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
