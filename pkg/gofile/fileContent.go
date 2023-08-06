// Package gofile is file and directory management libraries.
package gofile

import (
	"bytes"
)

// FindSubBytes find substrings, including start and end marks
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

// FindSubBytesNotIn find substrings, excluding start and end tags
func FindSubBytesNotIn(data []byte, start []byte, end []byte) []byte {
	startIndex := bytes.Index(data, start)
	endIndex := bytes.Index(data, end)
	if startIndex+len(start) >= endIndex {
		return []byte{}
	}
	return data[startIndex+len(start) : endIndex]
}
