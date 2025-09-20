package tieba

import (
	"strconv"
	"time"
)

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Itoa
//
// 泛型的 strconv.Itoa，整数转字符串.
func Itoa[I Int](i I) string {
	return strconv.FormatInt(int64(i), 10)
}

// Timestamp
// 秒级时间戳的字符串.
func Timestamp() string {
	return Itoa(time.Now().Unix())
}
