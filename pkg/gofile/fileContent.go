package gofile

import (
	"bytes"
)

// FindSubBytes 查找子字符串，包括开始和结束标记
func FindSubBytes(data []byte, start []byte, end []byte) []byte {
	startIndex := bytes.Index(data, start)
	endIndex := bytes.Index(data, end)
	if startIndex >= endIndex {
		return []byte{}
	}
	if len(data) >= endIndex+len(end) {
		endIndex += len(end)
	}
	return data[startIndex:endIndex]
}

// FindSubBytesNotIn 查找子字符串，不包括开始和结束标记
func FindSubBytesNotIn(data []byte, start []byte, end []byte) []byte {
	startIndex := bytes.Index(data, start)
	endIndex := bytes.Index(data, end)
	if startIndex+len(start) >= endIndex {
		return []byte{}
	}
	return data[startIndex+len(start) : endIndex]
}
