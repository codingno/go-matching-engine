package engine

import (
	"encoding/json"
)

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

func (orderBook *OrderBook) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, orderBook)
}

func (orderBook *OrderBook) ToJSON() []byte {
	str, _ := json.Marshal(orderBook)
	return str
}

// func printOrder(data []Order) {
// 	b, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	fmt.Print(string(b))
// }

// Add an order to the order book
func (book *OrderBook) addOrder(order Order) {
	var orderBookSide []Order
	if order.Side == 1 {
		orderBookSide = book.BuyOrders
	} else {
		orderBookSide = book.SellOrders
	}
	order.FillIndex = []int{}

	n := len(orderBookSide)
	var i int
	if n == 0 {
		orderBookSide = append(orderBookSide, order)
		if order.Side == 1 {
			book.BuyOrders = orderBookSide
		} else {
			book.SellOrders = orderBookSide
		}
		return
	}
	// for i = n - 1; i >= 0; i-- {
	for i = 0; i < n; i++ {
		sideOrder := orderBookSide[i]
		if order.Side == 1 {
			if sideOrder.Price < order.Price {
				break
			}
		} else {
			if sideOrder.Price > order.Price {
				break
			}
		}
	}
	orderBookSide = append(orderBookSide, order)
	// if i == n-1 {
	// 	return
	// } else {
	copy(orderBookSide[i+1:], orderBookSide[i:])
	orderBookSide[i] = order
	// }
	if order.Side == 1 {
		book.BuyOrders = orderBookSide
	} else {
		book.SellOrders = orderBookSide
	}
}

// Add a buy order to the order book
func (book *OrderBook) addBuyOrder(order Order) {
	n := len(book.BuyOrders)
	var i int
	if n == 0 {
		book.BuyOrders = append(book.BuyOrders, order)
		return
	}
	// for i = n - 1; i >= 0; i-- {
	for i = 0; i < n; i++ {
		buyOrder := book.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}
	}
	book.BuyOrders = append(book.BuyOrders, order)
	// if i == n-1 {
	// 	return
	// } else {
	copy(book.BuyOrders[i+1:], book.BuyOrders[i:])
	book.BuyOrders[i] = order
	// }
}

// Add a sell order to the order book
func (book *OrderBook) addSellOrder(order Order) {

	n := len(book.SellOrders)
	var i int
	if n == 0 {
		book.SellOrders = append(book.SellOrders, order)
		return
	}
	for i = 0; i < n; i++ {
		sellOrder := book.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
	}
	// if i == n-1 {
	// 	book.SellOrders = append(book.SellOrders, order)
	// 	return
	// } else {
	book.SellOrders = append(book.SellOrders, order)
	copy(book.SellOrders[i+1:], book.SellOrders[i:])
	book.SellOrders[i] = order
	// }
}

// Remove an order from the order book at a given index
func (book *OrderBook) removeOrder(index int, side int8) {
	var orderBookSide []Order
	if side == 0 {
		orderBookSide = book.BuyOrders
	} else {
		orderBookSide = book.SellOrders
	}
	orderBookSide = append(orderBookSide[:index], orderBookSide[index+1:]...)
	if side == 0 {
		book.BuyOrders = orderBookSide
	} else {
		book.SellOrders = orderBookSide
	}
}

func (book *OrderBook) removeByID(ID string, side int8) {
	var orderBookSide []Order
	if side == 0 {
		orderBookSide = book.BuyOrders
	} else {
		orderBookSide = book.SellOrders
	}
	for index, a := range orderBookSide {
		if a.ID == ID {
			orderBookSide = append(orderBookSide[:index], orderBookSide[index+1:]...)
		}
	}
	if side == 0 {
		book.BuyOrders = orderBookSide
	} else {
		book.SellOrders = orderBookSide
	}
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
