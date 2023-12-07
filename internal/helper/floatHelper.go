package helper

import (
	"math/big"
	"strconv"
	"strings"
)

func StringToFloat(val string) float64 {
	result, _ := strconv.ParseFloat(val, 64)

	return result
}

func FloatToString(val float64) string {
	return strconv.FormatFloat(val, 'f', 6, 64)
}

func StringToBigFloat(val string) *big.Float {
	result, _ := new(big.Float).SetString(val)

	return result
}

func BigFloatToString(val *big.Float) string {
	return val.String()
}

func FloatToBigFloat(val float64) *big.Float {
	return new(big.Float).SetFloat64(val)
}

func GetFloatPerc(val float64) int {
	s := strconv.FormatFloat(val, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}
