package krand

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	l := 100

	for i := 0; i < l; i++ {
		n := Int()
		assert.True(t, n >= 0 && n <= 100)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [0, max]
		n := Int(20)
		assert.True(t, n >= 0 && n <= 20)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [min, max]
		n := Int(10, 20)
		assert.True(t, n >= 10 && n <= 20)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [max, min]
		n := Int(20, 10)
		assert.True(t, n >= 10 && n <= 20)
	}
}

func TestFloat64(t *testing.T) {
	l := 100

	for i := 0; i < l; i++ {
		// randomly generate the default random number: [0, 100]
		f := Float64(0)
		assert.True(t, f >= 0 && f <= 100)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [0, max]
		f := Float64(1, 20)
		assert.True(t, f >= 0 && f <= 20)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [min, max]
		f := Float64(2, 10, 20)
		assert.True(t, f >= 10 && f <= 20)
	}

	for i := 0; i < l; i++ {
		// randomly generated numbers: [max, min]
		f := Float64(4, 20, 10)
		assert.True(t, f >= 10 && f <= 20)
	}
}

func TestString(t *testing.T) {
	assert.Equal(t, 6, len(String(R_NUM)))
	assert.Equal(t, 32, len(Bytes(R_NUM, 32)))

	assert.Equal(t, 6, len(String(R_UPPER)))
	assert.Equal(t, 32, len(Bytes(R_UPPER, 32)))

	assert.Equal(t, 6, len(String(R_LOWER)))
	assert.Equal(t, 32, len(Bytes(R_LOWER, 32)))

	assert.Equal(t, 6, len(String(R_NUM|R_UPPER)))
	assert.Equal(t, 32, len(Bytes(R_NUM|R_UPPER, 32)))

	assert.Equal(t, 6, len(String(R_NUM|R_LOWER)))
	assert.Equal(t, 32, len(Bytes(R_NUM|R_LOWER, 32)))

	assert.Equal(t, 6, len(String(R_All)))
	assert.Equal(t, 32, len(Bytes(R_All, 32)))
}

func TestNewID(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.GreaterOrEqual(t, NewID(), time.Now().UnixMilli()*1000000)
	}
}

func TestNewStringID(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.Equal(t, 16, len(NewStringID()))
	}
}

func TestNewNewSeriesID(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.Equal(t, 26, len(NewSeriesID()))
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int()
	}
}

func BenchmarkInt_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(10000)
	}
}

func BenchmarkFloat64_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64(0)
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64(2, 10000)
	}
}

func BenchmarkString_ALL_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_All)
	}
}

func BenchmarkString_ALL_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(R_All, 16)
	}
}

func BenchmarkBytes_ALL_6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(R_All)
	}
}

func BenchmarkBytes_ALL_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(R_All, 16)
	}
}

func BenchmarkNewID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewID()
	}
}

func BenchmarkNewStringID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewStringID()
	}
}

func BenchmarkNewSeriesID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSeriesID()
	}
}
