package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/uniex/exchanges"
)

func TestImplementsInterface(t *testing.T) {
	var exchange exchanges.ExchangeAPI
	var err error

	exchange, err = New("", "")
	assert.Nil(t, err)
	assert.NotNil(t, exchange)
}
