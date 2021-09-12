package api

import (
	"fmt"
	"time"

	"github.com/vanclief/ez"
	"github.com/vanclief/finmod/market"
	"github.com/vanclief/uniex/exchanges"
	"github.com/vanclief/uniex/exchanges/metaapi/client"
)

func (api *API) GetBalance() (*market.BalanceSnapshot, error) {
	const op = "MetaAPI.GetBalances"

	accountInfo, err := api.Client.GetAccountInformation()
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	fmt.Println("accountInfo", accountInfo)

	balance := &market.BalanceSnapshot{
		Balance:    accountInfo.Balance,
		Equity:     accountInfo.Equity,
		Margin:     accountInfo.Margin,
		FreeMargin: accountInfo.FreeMargin,
	}

	return balance, nil
}

func (api *API) GetAssets() (*market.AssetsSnashot, error) {
	const op = "MetaAPI.GetAssets"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

// GetOrders - Returns existing orders with their status
func (api *API) GetOrders(request *exchanges.GetOrdersRequest) ([]market.Order, error) {
	const op = "MetaAPI.Orders"

	orders := []market.Order{}

	if request.Status == exchanges.ClosedStatus {
		return nil, ez.New(op, ez.EINVALID, "THis API only returns open orders", nil)
	}

	metaOrders, err := api.Client.GetOrders()
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	for _, metaOrder := range metaOrders {

		openTime, err := time.Parse(time.RFC3339, metaOrder.Time)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		order := market.Order{
			ID: metaOrder.ID,
			// Pair: Unfilled for now as it is tricky having both symbols EURUSD & #US30
			Price:          metaOrder.OpenPrice,
			Volume:         metaOrder.Volume,
			ExecutedVolume: metaOrder.CurrentVolume,
			OpenTime:       openTime,
		}

		switch metaOrder.Type {
		case "ORDER_TYPE_BUY":
			order.Action = market.BuyAction
			order.Type = market.MarketOrder
		case "ORDER_TYPE_SELL":
			order.Action = market.SellAction
			order.Type = market.MarketOrder
		case "ORDER_TYPE_BUY_LIMIT":
			order.Action = market.BuyAction
			order.Type = market.LimitOrder
		case "ORDER_TYPE_SELL_LIMIT":
			order.Action = market.SellAction
			order.Type = market.LimitOrder
		}

		switch metaOrder.State {
		case "ORDER_STATE_PLACED":
			order.Status = market.UnfilledOrder
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// CreateOrder - Places a new order
func (api *API) CreateOrder(orderRequest *market.OrderRequest) (*market.Order, error) {
	const op = "MetaAPI.CreateOrder"

	// pair.Base.Symbol
	request := &client.MetatraderTrade{
		Symbol:    orderRequest.Pair.Symbol(),
		OpenPrice: orderRequest.Price,
		Volume:    orderRequest.Quantity,
	}

	switch orderRequest.Action {
	case market.BuyAction:
		if orderRequest.Type == market.MarketOrder {
			request.ActionType = "ORDER_TYPE_BUY"
		} else {
			request.ActionType = "ORDER_TYPE_BUY_LIMIT"
		}
	case market.SellAction:
		if orderRequest.Type == market.MarketOrder {
			request.ActionType = "ORDER_TYPE_SELL"
		} else {
			request.ActionType = "ORDER_TYPE_SELL_LIMIT"
		}
	}

	trade, err := api.Client.Trade(request)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	order := &market.Order{
		ID:       trade.OrderID,
		Action:   orderRequest.Action,
		Type:     orderRequest.Type,
		Pair:     orderRequest.Pair,
		Price:    orderRequest.Price,
		Volume:   orderRequest.Quantity,
		Status:   market.UnfilledOrder,
		OpenTime: time.Now(),
	}

	return order, nil
}

// CancelOrder - Cancels an existing order
func (api *API) CancelOrder(order *market.Order) (string, error) {
	const op = "MetaAPI.CancelOrder"

	request := &client.MetatraderTrade{
		OrderID:    order.ID,
		ActionType: "ORDER_CANCEL",
	}

	response, err := api.Client.Trade(request)
	if err != nil {
		return "", ez.Wrap(op, err)
	}

	return response.OrderID, nil
}

func (api *API) GetTrades(request *exchanges.GetTradesRequest) ([]market.Trade, error) {
	const op = "MetaAPI.GetPositions"
	return nil, ez.New(op, ez.ENOTIMPLEMENTED, "Not implemented", nil)
}

func (api *API) GetPositions(request *exchanges.GetPositionsRequest) ([]market.Position, error) {
	const op = "MetaAPI.GetPositions"

	positions := []market.Position{}

	if request.Status == exchanges.ClosedStatus {
		return nil, ez.New(op, ez.EINVALID, "THis API only returns open positions", nil)
	}

	metaPositions, err := api.Client.GetPositions()
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	for _, metaPosition := range metaPositions {

		openTime, err := time.Parse(time.RFC3339, metaPosition.Time)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}

		position := market.Position{
			ID:        metaPosition.ID,
			Open:      true,
			OpenPrice: metaPosition.OpenPrice,
			Quantity:  metaPosition.Volume,
			OpenTime:  openTime,
		}

		switch metaPosition.Type {
		case "POSITION_TYPE_BUY":
			position.Type = market.LongPosition
		case "POSITION_TYPE_SELL":
			position.Type = market.ShortPosition
		}

		positions = append(positions, position)
	}

	return positions, nil
}

func (api *API) UpdatePosition(position *market.Position, request *exchanges.UpdatePositionRequest) (*market.Position, error) {
	const op = "MetaAPI.UpdatePosition"

	metaRequest := &client.MetatraderTrade{
		PositionID: position.ID,
		ActionType: "POSITION_MODIFY",
	}

	if request.TakeProfit != 0 {
		metaRequest.TakeProfit = request.TakeProfit
		position.TakeProfit = market.PositionCloseOrder{Price: request.TakeProfit}
	}

	if request.StopLoss != 0 {
		metaRequest.StopLoss = request.StopLoss
		position.StopLoss = market.PositionCloseOrder{Price: request.StopLoss}
	}

	_, err := api.Client.Trade(metaRequest)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return position, nil
}

func (api *API) ClosePosition(position *market.Position) (string, error) {
	const op = "MetaAPI.ClosePosition"

	request := &client.MetatraderTrade{
		PositionID: position.ID,
		ActionType: "POSITION_CLOSE_ID",
	}

	response, err := api.Client.Trade(request)
	if err != nil {
		return "", ez.Wrap(op, err)
	}

	return response.PositionID, nil
}

// GetFundingAddress - Retrieves or generates a new deposit addresses for an asset

// WithdrawAsset - Places a withdrawal request

// CancelWithdraw - Cancels an asset withdrawal
