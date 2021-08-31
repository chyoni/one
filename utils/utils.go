package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

// HandleErr is cause panic if err is not nil.
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	enc := gob.NewEncoder(&aBuffer)
	HandleErr(enc.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(dec.Decode(i))
}

func EncodeAsJSON(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func DecodeAsJSON(v interface{}, data []byte) error {
	err := json.Unmarshal(data, v)
	return err
}
