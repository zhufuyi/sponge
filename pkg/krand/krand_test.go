package krand

import (
	"fmt"
	"testing"
)

func TestInt(t *testing.T) {
	l := 100

	fmt.Println("随机生成默认的随机数: [0, 100]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int())
	}

	fmt.Println("\n\n", "随机生成数: [0, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(20))
	}

	fmt.Println("\n\n", "随机生成数: [min, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(10, 20))
	}

	fmt.Println("\n\n", "随机生成数: [max, min]")
	for i := 0; i < l; i++ {
		fmt.Printf("%d ", Int(2000, 1000))
	}
}

func TestFloat64(t *testing.T) {
	l := 100

	fmt.Println("随机生成默认的随机数: [0, 100]")
	for i := 0; i < l; i++ {
		fmt.Printf("%.f ", Float64(0))
	}

	fmt.Println("\n\n", "随机生成数: [0, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%.1f ", Float64(1, 20))
	}

	fmt.Println("\n\n", "随机生成数: [min, max]")
	for i := 0; i < l; i++ {
		fmt.Printf("%.2f ", Float64(2, 10, 20))
	}

	fmt.Println("\n\n", "随机生成数: [max, min]")
	for i := 0; i < l; i++ {
		fmt.Printf("%.4f ", Float64(4, 2000, 1000))
	}
}

func TestString(t *testing.T) {
	// 随机纯数字
	fmt.Printf("%s\n", String(R_NUM))
	fmt.Printf("%s\n", Bytes(R_NUM, 32))

	// 随机大写字母
	fmt.Printf("%s\n", String(R_UPPER))
	fmt.Printf("%s\n", Bytes(R_UPPER, 32))

	// 随机小写字母
	fmt.Printf("%s\n", String(R_LOWER))
	fmt.Printf("%s\n", Bytes(R_LOWER, 32))

	// 随机数字、大写字母
	fmt.Printf("%s\n", String(R_NUM|R_UPPER))
	fmt.Printf("%s\n", Bytes(R_NUM|R_UPPER, 32))

	// 随机数字、小写字母
	fmt.Printf("%s\n", String(R_NUM|R_LOWER))
	fmt.Printf("%s\n", Bytes(R_NUM|R_LOWER, 32))

	// 随机数字、大写字母、小写字母
	fmt.Printf("%s\n", String(R_All))
	fmt.Printf("%s\n", Bytes(R_All, 32))
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
