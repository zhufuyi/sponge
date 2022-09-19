package encoding

import (
	"bytes"
	"encoding/gob"
)

// GobEncoding gob encode
type GobEncoding struct{}

// Marshal gob encode
func (g GobEncoding) Marshal(v interface{}) ([]byte, error) {
	var (
		buffer bytes.Buffer
	)

	err := gob.NewEncoder(&buffer).Encode(v)
	return buffer.Bytes(), err
}

// Unmarshal gob encode
func (g GobEncoding) Unmarshal(data []byte, value interface{}) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(value)
	if err != nil {
		return err
	}
	return nil
}
