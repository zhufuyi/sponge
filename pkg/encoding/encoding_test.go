package encoding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type obj struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func xEncoding(e Encoding) error {
	o1 := &obj{ID: 1, Name: "foo"}
	data, err := Marshal(e, o1)
	if err != nil {
		return err
	}

	o2 := &obj{}
	err = Unmarshal(e, data, o2)
	if err != nil {
		return err
	}

	if o1.ID != o2.ID {
		return errors.New("Unmarshal failed")
	}

	return nil
}

func TestEncoding(t *testing.T) {
	err := xEncoding(GobEncoding{})
	assert.NoError(t, err)

	err = xEncoding(JSONEncoding{})
	assert.NoError(t, err)

	err = xEncoding(JSONGzipEncoding{})
	assert.NoError(t, err)

	err = xEncoding(JSONSnappyEncoding{})
	assert.NoError(t, err)

	err = xEncoding(MsgPackEncoding{})
	assert.NoError(t, err)
}

type codec struct{}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	return []byte{}, nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	return nil
}

func (c codec) Name() string {
	return "json"
}

func TestRegisterCodec(t *testing.T) {
	defer func() { recover() }()

	RegisterCodec(&codec{})
	c := GetCodec("json")
	assert.NotNil(t, c)

	RegisterCodec(nil)
}

func BenchmarkJsonMarshal(b *testing.B) {
	a := make([]int, 0, 400)
	for i := 0; i < 400; i++ {
		a = append(a, i)
	}
	jsonEncoding := JSONEncoding{}
	for n := 0; n < b.N; n++ {
		_, err := jsonEncoding.Marshal(a)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	a := make([]int, 0, 400)
	for i := 0; i < 400; i++ {
		a = append(a, i)
	}
	jsonEncoding := JSONEncoding{}
	data, err := jsonEncoding.Marshal(a)
	if err != nil {
		b.Error(err)
	}
	var result []int
	for n := 0; n < b.N; n++ {
		err = jsonEncoding.Unmarshal(data, &result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMsgpack(b *testing.B) {
	// run the Fib function b.N times
	a := make([]int, 400)
	for i := 0; i < 400; i++ {
		a = append(a, i)
	}
	msgPackEncoding := MsgPackEncoding{}
	data, err := msgPackEncoding.Marshal(a)
	if err != nil {
		b.Error(err)
	}
	var result []int
	for n := 0; n < b.N; n++ {
		err = msgPackEncoding.Unmarshal(data, &result)
		if err != nil {
			b.Error(err)
		}
	}
}
