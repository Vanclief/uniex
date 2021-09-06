package utils

import (
  "fmt"
  "testing"
  "time"
)

func TestCalculateLimit(t *testing.T) {
  endTime := time.Now()
  startTime := endTime.Add(time.Hour * -3)
  answer := CalculateLimit(startTime, endTime, 2)
  fmt.Println(answer)
}
