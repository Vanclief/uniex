package models

import (
	"errors"
	"strconv"

	"github.com/vanclief/ez"
)

func getFloat64FromStr(value interface{}) (float64, error) {
	const op = "models.getFloat64FromStr"

	str, ok := value.(string)
	if !ok {
		return .0, errors.New("Field must be a string")
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return .0, err
	}
	return f, nil
}

func getFloat64(value interface{}) (float64, error) {
	const op = "models.getFloat64"

	f, ok := value.(float64)
	if !ok {
		return .0, errors.New("Field must be a float64")
	}
	return f, nil
}

func getTimestamp(value interface{}) (int64, error) {
	const op = "models.getTimestamp"

	f, ok := value.(float64)
	if !ok {
		return 0, ez.New(op, ez.EINVALID, "Field must be a float64", nil)
	}

	return int64(f), nil
}
