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

func TestEncodingError(t *testing.T) {
	gobE := GobEncoding{}
	// gob error test
	err := gobE.Unmarshal([]byte("foo"), nil)
	assert.Error(t, err)

	jsonE := JSONEncoding{}
	// json error test
	err = jsonE.Unmarshal([]byte("foo"), nil)
	assert.Error(t, err)

	jsonDE := JSONGzipEncoding{}
	// gzip error test
	_, err = jsonDE.Marshal(make(chan string))
	assert.Error(t, err)
	err = jsonDE.Unmarshal([]byte("foo"), nil)
	assert.Error(t, err)
	data, _ := GzipEncode([]byte("foo"))
	err = jsonDE.Unmarshal(data, make(chan string))
	assert.Error(t, err)
	_, err = GzipDecode(nil)
	assert.Error(t, err)

	jsonSE := JSONSnappyEncoding{}
	// snappy error test
	_, err = jsonSE.Marshal(make(chan string))
	assert.Error(t, err)
	err = jsonSE.Unmarshal([]byte("foo"), nil)
	assert.Error(t, err)

	msgE := MsgPackEncoding{}
	// pack error test
	err = msgE.Unmarshal([]byte("foo"), nil)
	assert.Error(t, err)
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
	RegisterCodec(&codec{})
	c := GetCodec("json")
	assert.NotNil(t, c)

	defer func() { recover() }()
	RegisterCodec(nil)
}

type codec2 struct{ codec }

func (c codec2) Name() string {
	return ""
}

func TestRegisterCodec2(t *testing.T) {
	defer func() { recover() }()
	RegisterCodec(&codec2{})
}

type encoder struct{}

func (e encoder) Marshal(v interface{}) ([]byte, error) {
	return nil, errors.New("mock Marshal error")
}

func (e encoder) Unmarshal(data []byte, v interface{}) error {
	return errors.New("mock Unmarshal error")
}

func (e encoder) MarshalBinary() (data []byte, err error) {
	return nil, nil
}

func (e encoder) UnmarshalBinary(data []byte) error {
	return nil
}

func TestMarshal(t *testing.T) {
	_, err := Marshal(encoder{}, nil)
	assert.Error(t, err)

	_, err = Marshal(nil, &encoder{})
	assert.NoError(t, err)

	_, err = Marshal(encoder{}, &encoder{})
	assert.NoError(t, err)
}

func TestUnmarshall(t *testing.T) {
	err := Unmarshal(encoder{}, nil, nil)
	assert.Error(t, err)

	err = Unmarshal(nil, []byte("foo"), &encoder{})
	assert.NoError(t, err)

	err = Unmarshal(encoder{}, []byte("foo"), &encoder{})
	assert.NoError(t, err)
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
