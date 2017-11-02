package plumber

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/linkedin/goavro"
)

// RecordEncoding encodes and decodes Avro records.
type RecordEncoding struct {
	codec *goavro.Codec
}

// NewRecordEncoding creates a record encoding
func NewRecordEncoding(codec *goavro.Codec) *RecordEncoding {
	return &RecordEncoding{codec}
}

// Encode a record as an Avro binary value.
func (e *RecordEncoding) Encode(v interface{}) ([]byte, error) {
	native := make(map[string]interface{})

	value := reflect.ValueOf(v)
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		valueField := value.FieldByName(typeField.Name)
		avroFieldName := avroFieldName(typeField)
		native[avroFieldName] = encodeField(valueField)
	}

	binary, err := e.codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, fmt.Errorf("Could encode message %v", native)
	}

	return binary, nil
}

func encodeField(valueField reflect.Value) interface{} {
	return valueField.Interface()
}

// Decode a record into an interface.
func (e *RecordEncoding) Decode(message []byte, v interface{}) error {
	value := reflect.ValueOf(v).Elem()
	vType := value.Type()

	native, _, err := e.codec.NativeFromBinary(message)
	if err != nil {
		fmt.Printf("Could not decode Avro message as %s: %s", vType.Name(), err)
		return err
	}

	record, ok := native.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid message, as it is not a record")
		return err
	}

	for i := 0; i < vType.NumField(); i++ {
		structField := vType.Field(i)
		avroName := avroFieldName(structField)

		recordValue, ok := record[avroName]
		if !ok {
			return fmt.Errorf("Missing field %s in record %v", avroName, record)
		}

		valueField := value.Field(i)
		valueField.Set(reflect.ValueOf(recordValue))
	}

	return nil
}

// DecodeFromType decodes a record, given only a reflect.Type and the bytes.
func (e *RecordEncoding) DecodeFromType(message []byte, messageType reflect.Type) (interface{}, error) {
	decodedValue := reflect.New(messageType)
	decodedInterface := decodedValue.Interface()
	err := e.Decode(message, decodedInterface)
	if err != nil {
		return nil, err
	}
	log.Printf("Decoded interface: %v", decodedInterface)
	return decodedInterface, nil
}

func avroFieldName(structField reflect.StructField) string {
	avroName, ok := structField.Tag.Lookup("avro")
	if ok {
		return strings.Trim(avroName, "\"")
	}
	return structField.Name
}
