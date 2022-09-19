package encoding

import "github.com/vmihailenco/msgpack"

// MsgPackEncoding msgpack 格式
type MsgPackEncoding struct{}

// Marshal msgpack encode
func (mp MsgPackEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := msgpack.Marshal(v)
	return buf, err
}

// Unmarshal msgpack decode
func (mp MsgPackEncoding) Unmarshal(data []byte, value interface{}) error {
	err := msgpack.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}
