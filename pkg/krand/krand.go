// Package krand is a library for generating random strings, integers, floating point numbers.
package krand

import (
	"math/rand"
	"time"
)

// nolint
const (
	R_NUM   = 1 // only number
	R_UPPER = 2 // only capital letters
	R_LOWER = 4 // only lowercase letters
	R_All   = 7 // numbers, upper and lower case letters
)

var (
	refSlices = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	kinds     = [][]byte{refSlices[0:10], refSlices[10:36], refSlices[0:36], refSlices[36:62], refSlices[36:], refSlices[10:62], refSlices[0:62]}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// String generate random strings of any length of multiple types, default length is 6 if size is empty
// example: String(R_ALL), String(R_ALL, 16), String(R_NUM|R_LOWER, 16)
func String(kind int, size ...int) string {
	return string(Bytes(kind, size...))
}

// Bytes generate random strings of any length of multiple types, default length is 6 if bytesLen is empty
// example: Bytes(R_ALL), Bytes(R_ALL, 16), Bytes(R_NUM|R_LOWER, 16)
func Bytes(kind int, bytesLen ...int) []byte {
	if kind > 7 || kind < 1 {
		kind = R_All
	}

	length := 6 // default length 6
	if len(bytesLen) > 0 {
		length = bytesLen[0]
		if length < 1 {
			length = 6
		}
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = kinds[kind-1][rand.Intn(len(kinds[kind-1]))]
	}

	return result
}

// Int generate random numbers of specified range size,
// compatible with Int(), Int(max), Int(min, max), Int(max, min) 4 ways, min<=random number<=max
func Int(rangeSize ...int) int {
	switch len(rangeSize) {
	case 0:
		return rand.Intn(101) // default 0~100
	case 1:
		return rand.Intn(rangeSize[0] + 1)
	default:
		if rangeSize[0] > rangeSize[1] {
			rangeSize[0], rangeSize[1] = rangeSize[1], rangeSize[0]
		}
		return rand.Intn(rangeSize[1]-rangeSize[0]+1) + rangeSize[0]
	}
}

// Float64 generates a random floating point number of the specified range size,
// Four types of passing references are supported, example: Float64(dpLength), Float64(dpLength, max),
// Float64(dpLength, min, max), Float64(dpLength, max, min), min<=random numbers<=max
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
		return float64(rand.Intn(101)) + dp // default 0~100
	case 1:
		return float64(rand.Intn(rangeSize[0]+1)) + dp
	default:
		if rangeSize[0] > rangeSize[1] {
			rangeSize[0], rangeSize[1] = rangeSize[1], rangeSize[0]
		}
		return float64(rand.Intn(rangeSize[1]-rangeSize[0]+1)+rangeSize[0]) + dp
	}
}
