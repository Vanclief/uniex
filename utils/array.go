package utils

import (
	"math"
	"time"
)

// TimeInterval contains the UNIX timestamps to make a Binance API call with startTime and endTime, the maximum
// difference between endTime and startTime is 1000 minutes or 16 hours, 40 minutes
type TimeInterval struct {
	StartTime int64
	EndTime   int64
}

// CreateArrayOfTimestamps takes the start and end time as time.Time and calculates how many Binance API calls are
// necessary to cover the time period in minutes, every API call interval is stored in an Interval struct
func CreateArrayOfTimestamps(startTime, endTime time.Time) (timestamps []TimeInterval) {
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()
	delta := int64(60 * 1000)
	loops := math.Ceil(float64(endUnix-startUnix) / float64(delta))
	if loops == 0 {
		return append(timestamps, TimeInterval{StartTime: startUnix, EndTime: endUnix})
	}

	startIndex := startUnix
	endIndex := int64(math.Min(float64(startUnix+delta), float64(endUnix)))

	for i := 0; i < int(loops); i++ {
		timestamps = append(timestamps, TimeInterval{
			StartTime: startIndex,
			EndTime:   endIndex,
		})
		startIndex += delta
		endIndex = int64(math.Min(float64(endIndex+delta), float64(endUnix)))
		//startIndex += delta + 1
		//endIndex = int64(math.Min(float64(endIndex+delta+1), float64(endUnix)))
	}
	return timestamps
}
