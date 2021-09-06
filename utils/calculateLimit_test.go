package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateLimit(t *testing.T) {
	endTime := time.Now()
	startTime := endTime.Add(time.Hour * -3)
	limit := CalculateLimit(startTime, endTime, 2, 1000)

	assert.Equal(t, int64(90), limit)

	endTime = time.Now()
	startTime = endTime.Add(time.Minute * -10)
	limit = CalculateLimit(startTime, endTime, 2, 1000)

	assert.Equal(t, int64(5), limit)

	endTime = time.Now()
	startTime = endTime.Add(time.Minute * -11)
	limit = CalculateLimit(startTime, endTime, 2, 1000)

	assert.Equal(t, int64(6), limit)

	endTime = time.Now()
	startTime = endTime.Add(time.Hour * -48)
	limit = CalculateLimit(startTime, endTime, 2, 1000)

	assert.Equal(t, int64(1000), limit)
}
