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
		FillReverse:      []FillReverse{},
	}

	tests := []TestMachine{
		{
			name: "0000-ZeroBookAddOneBuyOrder",
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
			name: "0001-ZeroBookAddOneSellOrder",
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
			name: "0002-SellBuy",
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
			name: "0003-BuySell",
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
			name: "0004-SellHighBuyDiffPrice",
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
			name: "0005-SellBuyHighDiffPrice",
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
			name: "0006-BuyHighSellDiffPrice",
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
			name: "0007-BuySellHighDiffPrice",
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
		{
			name: "0008-ManySellBuyValueAndBestPrice",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0009-ManyBuySellValueAndBestPrice",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     2,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7400,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0010-ManySellBuyValueAndPriceWorst",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7400,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0011-ManyBuySellValueAndPriceWorst",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     2,
					price:      7400,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0012-FillOrKillBuySell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: true,
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
			name: "0013-FillOrKillSellBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: true,
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
			name: "0014-FillOrKillBuyManySell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     4,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     1,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0015-FillOrKillSellManyBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     4,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     1,
					price:      7500,
					fillOrKill: false,
				},
			},
		},
		{
			name: "0016-ManySellMoreFillOrKillBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     4,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0017-ManyBuyMoreFillOrKillSell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     4,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0018-ManySellFillOrMoreFillOrKillBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     7,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0019-ManyBuyFillOrMoreFillOrKillSell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     7,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0020-ManySellFillOrMoreFillOrKillBuy2",
			request: []RequestOrder{
				{
					side:       0,
					amount:     1,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     7,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0021-ManyBuyFillOrMoreFillOrKillSell2",
			request: []RequestOrder{
				{
					side:       1,
					amount:     1,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     7,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0022-ManySellFillOrMoreFillOrKillBuy3",
			request: []RequestOrder{
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     4,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     15,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0023-ManyBuyFillOrMoreFillOrKillSell3",
			request: []RequestOrder{
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     4,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     15,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0024-ManySell2FillOrMoreFillOrKillBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     4,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       0,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     15,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0025-ManyBuy2FillOrMoreFillOrKillSell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     4,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     3,
					price:      7500,
					fillOrKill: false,
				},
				{
					side:       1,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     15,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0026-Buy1Fill2FillSell",
			request: []RequestOrder{
				{
					side:       0,
					amount:     1,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     1,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0027-Sell1Fill2FillBuy",
			request: []RequestOrder{
				{
					side:       1,
					amount:     1,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     1,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0028-Sell1Fill2FillBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0029-Buy1Fill2FillSell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0030-Sell2Fill2FillBuy",
			request: []RequestOrder{
				{
					side:       0,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     1,
					price:      7500,
					fillOrKill: true,
				},
			},
		},
		{
			name: "0031-Buy2Fill2FillSell",
			request: []RequestOrder{
				{
					side:       1,
					amount:     6,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     5,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       0,
					amount:     2,
					price:      7500,
					fillOrKill: true,
				},
				{
					side:       1,
					amount:     1,
					price:      7500,
					fillOrKill: true,
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
				moreTrades, _ := book.Process(order)
				trades = append(trades, moreTrades...)
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

			// ManySellBuyValue & ManyBuySellValue
			if i == 8 || i == 9 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       3,
						Price:        orderQueue[0].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       2,
						Price:        orderQueue[1].Price,
						CreatedAt:    createdAt,
					},
				}
			}

			// ManySellBuyValueAndPriceWorst & ManyBuySellValueAndPriceWorst
			if i == 10 || i == 11 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       3,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
				}

				order = orderQueue[2]
				order.Amount = 2

				if i == 10 {
					result.BuyOrders = append(result.BuyOrders, order)
					result.SellOrders = append(result.SellOrders, orderQueue[1])
				} else {
					result.SellOrders = append(result.SellOrders, order)
					result.BuyOrders = append(result.BuyOrders, orderQueue[1])
				}
			}

			// FillOrKillBuySell & FillOrKillBuyManySell
			if i == 12 || i == 14 {
				result.BuyOrders = append(result.BuyOrders, orderQueue[0])
				result.SellOrders = append(result.SellOrders, orderQueue[1])
				if i == 14 {
					result.SellOrders = append(result.SellOrders, orderQueue[2])
				}
			}

			// FillOrKillSellBuy & FillOrKillSellManyBuy
			if i == 13 || i == 15 {
				result.SellOrders = append(result.SellOrders, orderQueue[0])
				result.BuyOrders = append(result.BuyOrders, orderQueue[1])
				if i == 15 {
					result.BuyOrders = append(result.BuyOrders, orderQueue[2])
				}
			}

			// ManySellMoreFillOrKillBuy & ManyBuyMoreFillOrKillSell
			if i == 16 || i == 17 {
				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       2,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
				}

				if i == 17 {
					order = orderQueue[1]
					order.Amount = 1
					result.BuyOrders = append(result.BuyOrders, order)
				}

			}

			// ManySellFillOrMoreFillOrKillBuy & ManyBuyFillOrMoreFillOrKillSell
			if i == 18 || i == 19 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       5,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
				}

				order = orderQueue[0]
				order.Amount = 1

				if i == 18 {
					result.SellOrders = append(result.SellOrders, order)
				} else {
					result.BuyOrders = append(result.BuyOrders, order)
				}
			}

			// ManySellFillOrMoreFillOrKillBuy2 & ManyBuyFillOrMoreFillOrKillSell2
			if i == 20 || i == 21 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[3].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       1,
						Price:        orderQueue[3].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[3].ID,
						MakerOrderID: orderQueue[2].ID,
						Amount:       5,
						Price:        orderQueue[3].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[3].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       1,
						Price:        orderQueue[3].Price,
						CreatedAt:    createdAt,
					},
				}

				order = orderQueue[1]
				order.Amount = 2

				if i == 20 {
					result.SellOrders = append(result.SellOrders, order)
				} else {
					result.BuyOrders = append(result.BuyOrders, order)
				}
			}

			// ManySellFillOrMoreFillOrKillBuy3 & ManyBuyFillOrMoreFillOrKillSell3
			if i == 22 || i == 23 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[2].ID,
						Amount:       3,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[4].ID,
						Amount:       6,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       4,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
				}

				order = orderQueue[2]
				order.Amount = 2

				if i == 22 {
					result.SellOrders = append(result.SellOrders, order, orderQueue[3])
				} else {
					result.BuyOrders = append(result.BuyOrders, order, orderQueue[3])
				}
			}

			// ManySell2FillOrMoreFillOrKillBuy & ManyBuy2FillOrMoreFillOrKillSell
			if i == 24 || i == 25 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[2].ID,
						Amount:       5,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[4].ID,
						Amount:       6,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       2,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[5].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       2,
						Price:        orderQueue[5].Price,
						CreatedAt:    createdAt,
					},
				}

				order = orderQueue[1]
				order.Amount = 2

				if i == 24 {
					result.SellOrders = append(result.SellOrders, order, orderQueue[3])
				} else {
					result.BuyOrders = append(result.BuyOrders, order, orderQueue[3])
				}
			}

			// ManyBuy2FillOneFillSellOneFalseSell
			if i == 26 || i == 27 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       1,
						Price:        orderQueue[0].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[2].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       1,
						Price:        orderQueue[1].Price,
						CreatedAt:    createdAt,
					},
				}
			}

			// Sell1Fill2FillBuy
			if i == 28 {

				result.BuyOrders = append(result.BuyOrders, orderQueue[1], orderQueue[2])
				result.SellOrders = append(result.SellOrders, orderQueue[0])
			}

			// Buy1Fill2FillSell
			if i == 29 {

				result.BuyOrders = append(result.BuyOrders, orderQueue[0])
				result.SellOrders = append(result.SellOrders, orderQueue[1], orderQueue[2])
			}

			if i == 30 || i == 31 {

				resultTrades = []Trade{
					{
						TakerOrderID: orderQueue[0].ID,
						MakerOrderID: orderQueue[2].ID,
						Amount:       2,
						Price:        orderQueue[2].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[1].ID,
						MakerOrderID: orderQueue[0].ID,
						Amount:       4,
						Price:        orderQueue[0].Price,
						CreatedAt:    createdAt,
					},
					{
						TakerOrderID: orderQueue[3].ID,
						MakerOrderID: orderQueue[1].ID,
						Amount:       1,
						Price:        orderQueue[1].Price,
						CreatedAt:    createdAt,
					},
				}

			}

			// fmt.Println("================================================")
			// for i, v := range orderQueue {
			// 	fmt.Println(v.ID, v.Amount, v.Side, i)
			// }
			// fmt.Println("================== EXPECT ======================")
			// printJSON(resultTrades)
			// printJSON(result)
			// fmt.Println("================== ACTUAL ======================")
			// printJSON(trades)
			// printJSON(book)
			// fmt.Println("================================================")

			assert.Equal(t, resultTrades, trades)
			assert.Equal(t, result, book)

		})
	}
}
