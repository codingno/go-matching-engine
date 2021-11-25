package engine

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderBookFromJSON(t *testing.T) {
	var book OrderBook = OrderBook{}
	var orderBookJSON string
	var orderBookBuyJSON, orderBookSellJSON string = "[", "["
	var orderBook OrderBook = OrderBook{
		BuyOrders:  []Order{},
		SellOrders: []Order{},
	}
	var orderBookBuy, orderBookSell []Order
	var orderTakerJSON, orderMakerJSON string
	var orderTaker, orderMaker Order
	n := 10
	for i := 0; i < n; i++ {
		uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
		uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
		createdAt := time.Now().UTC().String()
		price := uint64(7400 + (i * 10))
		coma := ","
		if i == n-1 {
			coma = ""
		}
		orderTakerJSON = fmt.Sprintf(`
		{
			"id"								:	"%v",
			"price"							:	%v,
			"amountTemp"				:	0,
			"amount"						:	1,
			"side" 							:	0,
			"createdAt" 				:	"%v",
			"fillOrKill"				:	false,
			"fillIndex" 				:	[],
			"reverseCalculate"	: 0,
			"idCalculate"				: ""
		}
		`, uniqueTaker, price, createdAt)

		orderBookSellJSON = fmt.Sprintf(`%v%v%v`, orderBookSellJSON, orderTakerJSON, coma)

		orderTaker = Order{1, price, uniqueTaker, 0, createdAt, false, 0, []int{}, 0, ""}

		orderBookSell = append(orderBookSell, orderTaker)

		orderMakerJSON = fmt.Sprintf(`
		{
			"id"								:	"%v",
			"price"							:	%v,
			"amountTemp"				:	0,
			"amount"						:	1,
			"side" 							:	1,
			"createdAt" 				:	"%v",
			"fillOrKill"				:	false,
			"fillIndex" 				:	[],
			"reverseCalculate"	: 0,
			"idCalculate"				: ""
		}
		`, uniqueMaker, price, createdAt)

		orderBookBuyJSON = fmt.Sprintf(`%v%v%v`, orderBookBuyJSON, orderMakerJSON, coma)

		orderMaker = Order{1, price, uniqueMaker, 1, createdAt, false, 0, []int{}, 0, ""}

		orderBookBuy = append(orderBookBuy, orderMaker)

	}
	orderBookSellJSON = fmt.Sprintf(`%v]`, orderBookSellJSON)
	orderBookBuyJSON = fmt.Sprintf(`%v]`, orderBookBuyJSON)
	orderBookJSON = fmt.Sprintf(`{"BuyOrders":%v,"SellOrders":%v}`, orderBookBuyJSON, orderBookSellJSON)

	orderBook.SellOrders = orderBookSell
	orderBook.BuyOrders = orderBookBuy

	if err := book.FromJSON([]byte(orderBookJSON)); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, orderBook, book, "this orderbook should be equal")
}

func TestOrderBookToJSON(t *testing.T) {
	var orderBookJSON string
	var orderBookBuyJSON, orderBookSellJSON string = "[", "["
	var orderBook OrderBook = OrderBook{
		BuyOrders:  []Order{},
		SellOrders: []Order{},
	}
	var orderBookBuy, orderBookSell []Order
	var orderTakerJSON, orderMakerJSON string
	var orderTaker, orderMaker Order
	n := 10
	for i := 0; i < n; i++ {
		uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
		uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
		createdAt := time.Now().UTC().String()
		price := uint64(7400 + (i * 10))
		coma := ","
		if i == n-1 {
			coma = ""
		}
		orderTakerJSON = fmt.Sprintf(`{"amount":1,"price":%v,"id":"%v","side":0,"createdAt":"%v","fillOrKill":false,"amountTemp":0,"fillIndex":[],"reverseCalculate":0,"idCalculate":""}`, price, uniqueTaker, createdAt)

		orderBookSellJSON = fmt.Sprintf(`%v%v%v`, orderBookSellJSON, orderTakerJSON, coma)

		orderTaker = Order{1, price, uniqueTaker, 0, createdAt, false, 0, []int{}, 0, ""}

		orderBookSell = append(orderBookSell, orderTaker)

		orderMakerJSON = fmt.Sprintf(`{"amount":1,"price":%v,"id":"%v","side":1,"createdAt":"%v","fillOrKill":false,"amountTemp":0,"fillIndex":[],"reverseCalculate":0,"idCalculate":""}`, price, uniqueMaker, createdAt)

		orderBookBuyJSON = fmt.Sprintf(`%v%v%v`, orderBookBuyJSON, orderMakerJSON, coma)

		orderMaker = Order{1, price, uniqueMaker, 1, createdAt, false, 0, []int{}, 0, ""}

		orderBookBuy = append(orderBookBuy, orderMaker)

	}
	orderBookSellJSON = fmt.Sprintf(`%v]`, orderBookSellJSON)
	orderBookBuyJSON = fmt.Sprintf(`%v]`, orderBookBuyJSON)
	orderBookJSON = fmt.Sprintf(`{"BuyOrders":%v,"SellOrders":%v}`, orderBookBuyJSON, orderBookSellJSON)

	orderBook.SellOrders = orderBookSell
	orderBook.BuyOrders = orderBookBuy

	assert.Equal(t, orderBookJSON, string(orderBook.ToJSON()), "this orderbook should be equal")
}

func TestOrderBookTemp(t *testing.T) {
	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false
	book := OrderBook{
		BuyOrders: []Order{
			Order{
				3, price, uniqueMaker, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
			},
		},
		SellOrders: []Order{
			Order{
				3, price, uniqueTaker, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
			},
		},
	}

	orderBookSide := book.orderBookTemp(1)

	assert.Equal(t, orderBookSide, book.SellOrders, "This should be equal")

	orderBookSide = book.orderBookTemp(0)

	assert.Equal(t, orderBookSide, book.BuyOrders, "This should be equal")
}

func TestUpdateOrderBook(t *testing.T) {
	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false

	for i := 0; i < 2; i++ {
		testName := "Sell"

		if i == 1 {
			testName = "Buy"
		}

		t.Run(testName, func(t *testing.T) {

			book := OrderBook{
				BuyOrders: []Order{
					Order{
						3, price, uniqueMaker, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
					},
				},
				SellOrders: []Order{
					Order{
						3, price, uniqueTaker, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
					},
				},
			}

			uniqueID := strings.Replace(uuid.New().String(), "-", "", -1)
			order := Order{
				3, price, uniqueID, int8(i), createdAt, fillOrKill, 0, []int{}, 0, "",
			}

			result := append(book.SellOrders, order)

			book.updateOrderBook(int8(i), result)

			var checkResult []Order

			if i == 0 {
				checkResult = book.BuyOrders
			} else {
				checkResult = book.SellOrders
			}

			assert.Equal(t, result, checkResult, "This should be equal")

		})
	}
}

func TestOrderBookSideContainsByID(t *testing.T) {

	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false

	orders := []Order{
		Order{
			3, price, uniqueMaker, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
		Order{
			3, price, uniqueTaker, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
	}

	var book OrderBook

	assert.True(t, book.contains(orders, uniqueMaker))

	assert.True(t, book.contains(orders, uniqueTaker))

	assert.False(t, book.contains(orders, "testID"))
}

func TestOrderBookGetIndexByID(t *testing.T) {

	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueMaker2 := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false

	book := OrderBook{
		BuyOrders: []Order{
			Order{
				3, price, uniqueMaker, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
			},
			Order{
				3, price, uniqueMaker2, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
			},
		},
		SellOrders: []Order{
			Order{
				3, price, uniqueTaker, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
			},
		},
	}

	index, err := book.getIndexByID(uniqueMaker, 1)
	assert.False(t, err)
	assert.Equal(t, uniqueMaker, book.BuyOrders[index].ID)

	_, err = book.getIndexByID(uniqueMaker, 0)
	assert.True(t, err)

	index, err = book.getIndexByID(uniqueTaker, 0)
	assert.False(t, err)
	assert.Equal(t, uniqueTaker, book.SellOrders[index].ID)
}

func TestOrderBookAddOrder(t *testing.T) {
	uniqueID := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false

	orders := []Order{
		Order{
			3, price, uniqueID, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
		Order{
			3, price, uniqueID, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
	}

	for _, order := range orders {
		testName := "sell"

		if order.Side == 1 {
			testName = "buyer"
		}

		t.Run(testName, func(t *testing.T) {

			var book = OrderBook{
				BuyOrders:  []Order{},
				SellOrders: []Order{},
			}

			orderArray := []Order{
				order,
			}

			var result OrderBook
			if order.Side == 1 {
				result = OrderBook{
					BuyOrders:  orderArray,
					SellOrders: []Order{},
				}
			} else {
				result = OrderBook{
					BuyOrders:  []Order{},
					SellOrders: orderArray,
				}
			}

			book.addOrder(order)

			assert.Equal(t, result, book, "this should be equal")

		})
	}
}

func TestOrderBookRemoveOrder(t *testing.T) {
	uniqueID := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()
	price := uint64(7400)
	fillOrKill := false

	orders := []Order{
		Order{
			3, price, uniqueID, 1, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
		Order{
			3, price, uniqueID, 0, createdAt, fillOrKill, 0, []int{}, 0, "",
		},
	}

	for _, order := range orders {
		testName := "sell"
		removeSide := 1

		if order.Side == 1 {
			testName = "buyer"
			removeSide = 0
		}

		t.Run(testName, func(t *testing.T) {

			subTestNames := []string{
				"byID",
				"byIndex",
			}

			for _, subTestName := range subTestNames {
				t.Run(subTestName, func(t *testing.T) {

					var result = OrderBook{
						BuyOrders:  []Order{},
						SellOrders: []Order{},
					}

					orderArray := []Order{
						order,
					}

					var book OrderBook
					if order.Side == 1 {
						book = OrderBook{
							BuyOrders:  orderArray,
							SellOrders: []Order{},
						}
					} else {
						book = OrderBook{
							BuyOrders:  []Order{},
							SellOrders: orderArray,
						}
					}

					if subTestName == "byID" {
						book.removeByID(order.ID, int8(removeSide))
					} else {
						book.removeOrder(0, int8(removeSide))
					}

					assert.Equal(t, result, book, "this should be equal")
				})
			}

		})
	}
}
