package utils

import "math"

func RoundUpInt(x float64) int64 {
	return int64(math.Ceil(x))
}
