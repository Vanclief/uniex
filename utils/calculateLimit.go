package utils

import (
  "math"
  "time"
)

func CalculateLimit(startTime, endTime time.Time, interval int) int64 {
  startUnix := startTime.Unix()
  endUnix := endTime.Unix()
  if startUnix > endUnix {
    return int64(math.Min(float64((startUnix-endUnix)/(int64(interval)*60)), 1000))
  }
  return int64(math.Min(float64((endUnix-startUnix)/(int64(interval)*60)), 1000))
}
