package encoding

import "testing"

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
