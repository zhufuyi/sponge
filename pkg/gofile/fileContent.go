// Package gofile is file and directory management libraries.
package gofile

import (
	"bytes"
)

// FindSubBytes find first substrings, including start and end marks
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

// FindAllSubBytes find all substrings, including start and end marks
func FindAllSubBytes(data []byte, start []byte, end []byte) [][]byte {
	subBytes := [][]byte{}

	for {
		subString, endIndex := findSubByte2(data, start, end)
		if len(subString) == 0 {
			break
		}
		subBytes = append(subBytes, subString)
		data = data[endIndex:]
	}

	return subBytes
}

func findSubByte2(data []byte, start []byte, end []byte) ([]byte, int) {
	startIndex := bytes.Index(data, start)
	endIndex := bytes.Index(data, end)
	if startIndex >= endIndex {
		return []byte{}, 0
	}
	if len(data) >= endIndex+len(end) {
		endIndex += len(end)
	}
	return data[startIndex:endIndex], endIndex
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
