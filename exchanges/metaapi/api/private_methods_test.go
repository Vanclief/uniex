package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
)

func TestGetBalance(t *testing.T) {
	response, err := metaAPI.GetBalance()
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestGetOrders(t *testing.T) {
	// Open
	request := &exchanges.GetOrdersRequest{}
	response, err := metaAPI.GetOrders(request)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	// By IDs
	request = &exchanges.GetOrdersRequest{IDs: []string{"45105566"}}
	response, err = metaAPI.GetOrders(request)
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestCreateOrder(t *testing.T) {

	base := &market.Asset{Symbol: "BTC"}
	quote := &market.Asset{Symbol: "USD"}

	request := &market.OrderRequest{
		Action:   market.BuyAction,
		Type:     market.LimitOrder,
		Pair:     &market.Pair{Base: base, Quote: quote},
		Price:    50000,
		Quantity: 0.111,
	}
	response, err := metaAPI.CreateOrder(request)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.ID)
	assert.Equal(t, request.Price, response.Price)
}

func TestUpdateOrder(t *testing.T) {

	base := &market.Asset{Symbol: "BTC"}
	quote := &market.Asset{Symbol: "USD"}

	request := &market.OrderRequest{
		Action:   market.SellAction,
		Type:     market.LimitOrder,
		Pair:     &market.Pair{Base: base, Quote: quote},
		Price:    50000,
		Quantity: 0.1,
	}
	order, err := metaAPI.CreateOrder(request)
	assert.Nil(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, request.Price, order.Price)

	updateRequest := &exchanges.UpdateOrderRequest{
		Price: 51000,
	}
	order, err = metaAPI.UpdateOrder(order, updateRequest)
	assert.Nil(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, updateRequest.Price, order.Price)
}

func TestCancelOrder(t *testing.T) {

	base := &market.Asset{Symbol: "BTC"}
	quote := &market.Asset{Symbol: "USD"}

	// Create the order
	request := &market.OrderRequest{
		Action:   market.SellAction,
		Type:     market.LimitOrder,
		Pair:     &market.Pair{Base: base, Quote: quote},
		Price:    50000,
		Quantity: 0.1,
	}

	order, err := metaAPI.CreateOrder(request)
	assert.Nil(t, err)

	// Cancel the order
	id, err := metaAPI.CancelOrder(order)
	assert.Nil(t, err)
	assert.Equal(t, order.ID, id)
}

func TestGetPositions(t *testing.T) {
	// Open
	request := &exchanges.GetPositionsRequest{}
	response, err := metaAPI.GetPositions(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)

	// By IDs
	request = &exchanges.GetPositionsRequest{IDs: []string{"45424985"}}
	response, err = metaAPI.GetPositions(request)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.Len(t, response, 1)
}

func TestUpdatePosition(t *testing.T) {
	position := &market.Position{ID: "45099490"}
	request := &exchanges.UpdatePositionRequest{TakeProfit: 50000, StopLoss: 40000}

	response, err := metaAPI.UpdatePosition(position, request)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.TakeProfit, position.TakeProfit.Price)
	assert.Equal(t, request.StopLoss, position.StopLoss.Price)
}

func TestClosePosition(t *testing.T) {

	position := &market.Position{ID: "45075431"}

	// Cancel the order
	id, err := metaAPI.ClosePosition(position)
	assert.Nil(t, err)
	assert.Equal(t, position.ID, id)
}
