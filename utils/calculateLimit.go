package utils

import (
	"math"
	"time"
)

func CalculateLimit(startTime, endTime time.Time, interval int, limit float64) int64 {
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	var numerator float64

	if startUnix > endUnix {
		numerator = float64(startUnix - endUnix)
	} else {
		numerator = float64((endUnix - startUnix))
	}

	denominator := int64(interval) * 60

	return int64(math.Min(math.Ceil(numerator/float64(denominator)), limit))
}
