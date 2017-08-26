package utilities

import "math"

func Round(f float64, prec int) float64 {
	pow10_n := math.Pow10(prec)
	return math.Trunc((f + 0.5/pow10_n) * pow10_n ) / pow10_n
}
