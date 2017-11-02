package plumber

import (
	"log"
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

	encoding := RecordEncoding{codec}

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

	encoding := RecordEncoding{codec}

	original := IntRecord{42}

	encoded, _ := encoding.Encode(original)

	var decoded IntRecord
	encoding.Decode(encoded, &decoded)

	if original != decoded {
		t.Errorf("Original %v did not match decoded %v", original, decoded)
	}
}

func TestDecodeFromType_DecodeIntField_DecodedMatchesOriginal(t *testing.T) {
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

	encoding := RecordEncoding{codec}

	original := IntRecord{42}

	encoded, _ := encoding.Encode(original)

	decodedInterface, err := encoding.DecodeFromType(encoded, reflect.TypeOf(IntRecord{}))
	if err != nil {
		t.Fatalf("Could not decode IntRecord")
	}

	log.Println("Type of decodedInterface:", reflect.TypeOf(decodedInterface))

	decoded, ok := decodedInterface.(*IntRecord)
	if !ok {
		t.Fatalf("Could not assert decoded to *IntRecord")
	}

	log.Println("Type of decoded:", reflect.TypeOf(decoded))

	if original != *decoded {
		t.Errorf("Original %v did not match decoded %v", original, *decoded)
	}
}
