package hello_world

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

func SerializeSayHelloRequest(name string) ([]byte, error) {
	utf16Name := utf16.Encode([]rune(name))
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, utf16Name)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeSayHelloRequest(data []byte) (string, error) {
	if len(data)%2 != 0 {
		return "", errors.New("invalid data length")
	}

	utf16Name := make([]uint16, len(data)/2)
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &utf16Name)
	if err != nil {
		return "", err
	}

	return string(utf16.Decode(utf16Name)), nil
}

func SerializeSayHelloResponse(name string) ([]byte, error) {
	return SerializeSayHelloRequest(name)
}

func DeserializeSayHelloResponse(data []byte) (string, error) {
	return DeserializeSayHelloRequest(data)
}
