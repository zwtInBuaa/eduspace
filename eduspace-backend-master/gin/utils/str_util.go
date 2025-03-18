package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ParseStringToUint(str string) (uint, error) {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %v to uint: %v", str, err)
	}
	return uint(i), nil
}

func ParseStringToInt64(str string) (int64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %v to int64: %v", str, err)
	}
	return int64(i), nil
}

func ParseFloat64ToInt64(f float64) (int64, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("cannot convert %v to int64", f)
	}
	if f > float64(math.MaxInt64) || f < float64(math.MinInt64) {
		return 0, fmt.Errorf("%v out of int64 range", f)
	}
	return int64(f), nil
}

func ParseFloat64ToUint(f float64) (uint, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("cannot convert %v to int64", f)
	}
	if f > float64(math.MaxInt64) || f < float64(math.MinInt64) {
		return 0, fmt.Errorf("%v out of int64 range", f)
	}
	return uint(f), nil
}
