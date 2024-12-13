package utils

type INumber interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 |
		~int64 | ~uint64 | ~float32 | ~float64 | uintptr
}

// Abs 计算绝对值
func Abs[T INumber](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
