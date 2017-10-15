package plumber

import (
	"reflect"
	"testing"

	"github.com/linkedin/goavro"
)

type StringRecord struct {
	Field1 string `avro:"field1"`
	Field2 string `avro:"field2"`
}

type IntRecord struct {
	Field1 int32 `avro:"field1"`
}

func TestEncodeDecode_AllStringFields_DecodedMatchesOriginal(t *testing.T) {
	codec, err := goavro.NewCodec(`
		{
			"namespace": "plumber",
			"type": "record",
			"name": "StringTest",
			"fields": [
			  {"name": "field1",  "type": "string"},
			  {"name": "field2",  "type": "string"}
			]
		}`)
	if err != nil {
		t.Errorf("Codec is nil: %v", err)
		return
	}

	encoding := RecordEncoding{codec, reflect.TypeOf(StringRecord{})}

	original := StringRecord{"value1", "value2"}

	encoded, _ := encoding.Encode(original)

	var decoded StringRecord
	encoding.Decode(encoded, &decoded)

	if original != decoded {
		t.Errorf("Original %v did not match decodede %v", original, decoded)
	}
}

func TestEncodeDecode_IntField_DecodedMatchesOriginal(t *testing.T) {
	codec, err := goavro.NewCodec(`
		{
			"namespace": "plumber",
			"type": "record",
			"name": "IntTest",
			"fields": [
			  {"name": "field1",  "type": "int"}
			]
		}`)
	if err != nil {
		t.Errorf("Codec is nil: %v", err)
		return
	}

	encoding := RecordEncoding{codec, reflect.TypeOf(IntRecord{})}

	original := IntRecord{42}

	encoded, _ := encoding.Encode(original)

	var decoded IntRecord
	encoding.Decode(encoded, &decoded)

	if original != decoded {
		t.Errorf("Original %v did not match decodede %v", original, decoded)
	}
}
