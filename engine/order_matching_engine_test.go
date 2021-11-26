package engine

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type RequestOrder struct {
	side       int8
	amount     uint64
	price      uint64
	fillOrKill bool
}

type TestMachine struct {
	name    string
	request []RequestOrder
}

func TestProcess(t *testing.T) {

	createdAt := time.Now().UTC().String()
	var side int8 = 1

	uniqueID := strings.Replace(uuid.New().String(), "-", "", -1)
	order := Order{
		Amount:           1,
		Price:            7400,
		ID:               uniqueID,
		Side:             side,
		CreatedAt:        createdAt,
		FillOrKill:       false,
		AmountTemp:       0,
		FillIndex:        []int{},
		ReverseCalculate: 0,
		IDCalculate:      "",
	}

	tests := []TestMachine{
		{
			name: "ZeroBookAddOneBuyOrder",
			request: []RequestOrder{
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "ZeroBookAddOneSellOrder",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "SellBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "BuySell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "SellHighBuyDiffPrice",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
			},
		},
		{
			name: "SellBuyHighDiffPrice",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "BuyHighSellDiffPrice",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
			},
		},
		{
			name: "BuySellHighDiffPrice",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var orderQueue []Order

			initTrade := []Trade{}

			initBook := OrderBook{
				BuyOrders:  []Order{},
				SellOrders: []Order{},
			}

			result, book, resultTrades, trades := initBook, initBook, initTrade, initTrade

			for _, request := range test.request {

				uniqueID = strings.Replace(uuid.New().String(), "-", "", -1)
				createdAt = time.Now().UTC().String()
				order.ID = uniqueID
				order.CreatedAt = createdAt
				order.Side = request.side
				order.Amount = request.amount
				order.Price = request.price
				order.FillOrKill = request.fillOrKill

				orderQueue = append(orderQueue, order)

			}

			for _, order := range orderQueue {
				trades, order = book.Process(order)
			}

			for i, trade := range trades {
				if trade.CreatedAt != "" {
					trades[i].CreatedAt = createdAt
				}
			}

			// ZeroBookAddOneBuyOrder
			if i == 0 {
				result.BuyOrders = append(result.BuyOrders, order)
			}

			// ZeroBookAddOneSellOrder
			if i == 1 {
				result.SellOrders = append(result.SellOrders, order)
			}

			// SellBuy
			if i == 2 {
				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        order.Price,
						CreatedAt:    createdAt,
					},
				}
				order = orderQueue[0]
				order.Amount = 1
				result.SellOrders = append(result.SellOrders, order)
			}

			// BuySell
			if i == 3 {
				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        order.Price,
						CreatedAt:    createdAt,
					},
				}
				order = orderQueue[0]
				order.Amount = 1
				result.BuyOrders = append(result.BuyOrders, order)
			}

			// SellHighBuyDiffPrice
			if i == 4 {
				result.SellOrders = append(result.SellOrders, orderQueue[0])
				result.BuyOrders = append(result.BuyOrders, orderQueue[1])
			}

			// SellBuyHighDiffPrice
			if i == 5 {
				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       3,
						Price:        orderQueue[0].Price,
						CreatedAt:    createdAt,
					},
				}
				// Orderbook Zero
			}

			// BuyHighSellDiffPrice
			if i == 6 {
				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       3,
						Price:        orderQueue[0].Price,
						CreatedAt:    createdAt,
					},
				}
				// Orderbook Zero
			}

			// BuySellHighDiffPrice
			if i == 7 {
				result.BuyOrders = append(result.BuyOrders, orderQueue[0])
				result.SellOrders = append(result.SellOrders, orderQueue[1])
			}

			assert.Equal(t, resultTrades, trades)
			assert.Equal(t, result, book)

		})
	}
}
