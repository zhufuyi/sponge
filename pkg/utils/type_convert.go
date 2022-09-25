package utils

import "strconv"

// StrToInt string to int
func StrToInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

// StrToIntE string to int with error
func StrToIntE(str string) (int, error) {
	return strconv.Atoi(str)
}

// StrToUint32 string to uint32
func StrToUint32(str string) uint32 {
	v, _ := strconv.ParseUint(str, 10, 64)
	return uint32(v)
}

// StrToUint32E string to uint32 with error
func StrToUint32E(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint32(v), nil
}

// StrToUint64 string to uint64
func StrToUint64(str string) uint64 {
	v, _ := strconv.ParseUint(str, 10, 64)
	return v
}

// StrToUint64E string to uint64 with error
func StrToUint64E(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// StrToFloat32 string to float32
func StrToFloat32(str string) float32 {
	v, _ := strconv.ParseFloat(str, 32)
	return float32(v)
}

// StrToFloat32E string to float32 with error
func StrToFloat32E(str string) (float32, error) {
	v, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

// StrToFloat64 string to float64
func StrToFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

// StrToFloat64E string to float64 with error
func StrToFloat64E(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// IntToStr int to string
func IntToStr(v int) string {
	return strconv.Itoa(v)
}

// Uint64ToStr uint64 to string
func Uint64ToStr(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// Int64ToStr int64 to string
func Int64ToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}
