package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKraken(t *testing.T) {

	kraken, err := NewKraken()
	assert.Nil(t, err)
	kraken.API.GetPositions()

}
