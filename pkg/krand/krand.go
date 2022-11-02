package krand

import (
	"math/rand"
	"time"
)

// nolint
const (
	// R_NUM 纯数字
	R_NUM = 1
	// R_UPPER 大写字母
	R_UPPER = 2
	// R_LOWER 小写字母
	R_LOWER = 4
	// R_All 数字、大小写字母
	R_All = 7
)

var (
	refSlices = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	kinds     = [][]byte{refSlices[0:10], refSlices[10:36], refSlices[0:36], refSlices[36:62], refSlices[36:], refSlices[10:62], refSlices[0:62]}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// String 生成多种类型的任意长度的随机字符串，如果参数size为空，默认长度为6
// example：String(R_ALL), String(R_ALL, 16), String(R_NUM|R_LOWER, 16)
func String(kind int, size ...int) string {
	return string(Bytes(kind, size...))
}

// Bytes 生成多种类型的任意长度的随机字符串，如果参数bytesLen为空，默认长度为6
// example：Bytes(R_ALL), Bytes(R_ALL, 16), Bytes(R_NUM|R_LOWER, 16)
func Bytes(kind int, bytesLen ...int) []byte {
	if kind > 7 || kind < 1 {
		kind = R_All
	}

	length := 6 // 默认长度
	if len(bytesLen) > 0 {
		length = bytesLen[0] // 只有第0个值有效，忽略其它值
		if length < 1 {
			length = 6 // 默认长度
		}
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = kinds[kind-1][rand.Intn(len(kinds[kind-1]))]
	}

	return result
}

// Int 生成指定范围大小随机数，兼容Int()，Int(max)，Int(min, max)，Int(max, min)4种方式，min<=随机数<=max
func Int(rangeSize ...int) int {
	switch len(rangeSize) {
	case 0:
		return rand.Intn(101) // 默认0~100
	case 1:
		return rand.Intn(rangeSize[0] + 1)
	default:
		if rangeSize[0] > rangeSize[1] {
			rangeSize[0], rangeSize[1] = rangeSize[1], rangeSize[0]
		}
		return rand.Intn(rangeSize[1]-rangeSize[0]+1) + rangeSize[0]
	}
}

// Float64 生成指定范围大小随机浮点数，
// 支持4种传参方式：Float64(dpLength)，Float64(dpLength, max)，Float64(dpLength, min, max)，Float64(dpLength, max, min)，min<=随机数<=max
func Float64(dpLength int, rangeSize ...int) float64 {
	dp := 0.0
	if dpLength > 0 {
		dpmax := 1
		for i := 0; i < dpLength; i++ {
			dpmax *= 10
		}
		dp = float64(rand.Intn(dpmax)) / float64(dpmax)
	}

	switch len(rangeSize) {
	case 0:
		return float64(rand.Intn(101)) + dp // 默认0~100
	case 1:
		return float64(rand.Intn(rangeSize[0]+1)) + dp
	default:
		if rangeSize[0] > rangeSize[1] {
			rangeSize[0], rangeSize[1] = rangeSize[1], rangeSize[0]
		}
		return float64(rand.Intn(rangeSize[1]-rangeSize[0]+1)+rangeSize[0]) + dp
	}
}
