package helper

import (
	"math/big"
)

func SumWithPercent(val float64, percent float64) float64 {
	return val + getDiffValue(val, percent)
}

func DiffWithPercent(val float64, percent float64) float64 {
	return val - getDiffValue(val, percent)
}

func BigSumWithPercent(val *big.Float, percent *big.Float) *big.Float {
	diffval := new(big.Float).Copy(val)
	diffValue := getBigDiffValue(diffval, percent)

	return val.Add(val, diffValue)
}

func BigDiffWithPercent(val *big.Float, percent *big.Float) *big.Float {
	diffval := new(big.Float).Copy(val)
	diffValue := getBigDiffValue(diffval, percent)

	return val.Sub(val, diffValue)
}

func getBigDiffValue(val *big.Float, percent *big.Float) *big.Float {
	all := new(big.Float).SetFloat64(100)
	diffValue := val.Mul(val, percent)

	return diffValue.Quo(diffValue, all)
}

func getDiffValue(val float64, percent float64) float64 {
	return (val * percent) / 100
}
