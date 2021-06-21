package models

import "time"

type Kline struct {
  OpenTime                 time.Time `json:"open_time"`
  Open                     float64   `json:"open"`
  High                     float64   `json:"high"`
  Low                      float64   `json:"low"`
  Close                    float64   `json:"close"`
  Volume                   float64   `json:"volume"`
  CloseTime                time.Time `json:"close_time"`
}