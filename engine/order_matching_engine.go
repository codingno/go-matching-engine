package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order Order) ([]Trade, Order) {
	trades := make([]Trade, 0, 1)
	// var orderBookSide []Order
	// if order.Side == 1 {
	// 	orderBookSide = book.SellOrders
	// } else {
	// 	orderBookSide = book.BuyOrders
	// }
	orderBookSide := book.orderBookTemp(order.Side)
	n := len(orderBookSide)
	orderTemp := order
	if orderTemp.AmountTemp != 0 {
		orderTemp.Amount = orderTemp.AmountTemp
	}
	// check if we have at least one matching order
	if n != 0 {
		var bestPrice bool
		if order.Side == 1 {
			bestPrice = orderBookSide[n-1].Price <= orderTemp.Price
		} else {
			bestPrice = orderBookSide[n-1].Price >= orderTemp.Price
		}
		// if orderBookSide[n-1].Price <= orderTemp.Price {
		if bestPrice {
			// traverse all orders that match
			for i := 0; i < n; i++ {
				orderSide := orderBookSide[i]
				var skipPrice bool

				if order.Side == 1 {
					skipPrice = orderSide.Price > orderTemp.Price
				} else {
					skipPrice = orderSide.Price < orderTemp.Price
				}

				if skipPrice {
					break
				}

				if orderTemp.FillOrKill {
					if len(orderTemp.FillIndex) > 0 {
						if orderBookSide[orderTemp.FillIndex[0]].ID == orderSide.ID {
							continue
						}
						if orderTemp.FillIndex[0] > i {
							continue
						}
					}
					fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 97 ~ func ~ orderSide.Amount -= orderTemp.Amount`, orderTemp.Amount, orderSide.Amount)
					if orderSide.Amount >= orderTemp.Amount {
						orderSide.Amount -= orderTemp.Amount
						if orderSide.FillOrKill && orderSide.Amount != 0 {
							// calculate origin value

							orderSide = orderBookSide[i]
							if orderSide.Amount >= order.Amount {
								fmt.Println("asu", orderSide.Amount, order.Amount)
								orderSide.Amount -= order.Amount
							} else {
								order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
								order.IDCalculate = orderSide.ID
								orderSide.Amount = 0
								order.AmountTemp = 0
							}
							// else {
							// 	orderTemp = order
							// 	orderTemp.Amount -= orderSide.Amount
							// 	orderSide.Amount = 0
							// }

							orderSideAmount := orderSide.Amount

							if orderSide.Amount != 0 {
								orderSide = orderBookSide[i]
								orderSide.AmountTemp = orderSideAmount
								if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
									orderSide.FillIndex = append([]int{index}, order.FillIndex...)
								}
								fmt.Println("######################### orderside amount != 0 ######################### ")
								printJSON(orderSide)
								printJSON(order)
								fmt.Println("#########################                        ######################### ")
								moreTrades, moreOrder := book.Process(orderSide)
								if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
									orderSide = moreOrder
									orderSide.Amount = orderSide.AmountTemp
									trades = append(trades, moreTrades...)
								} else {
									continue
								}
								// continue
							}

							// orderTemp.Amount = order.Amount
							order.AmountTemp = orderTemp.Amount
							if order.AmountTemp == 0 {
								order.FillIndex = nil
							} else {
								order.AmountTemp = 0
							}
						} else {
							order.AmountTemp -= orderTemp.Amount
						}

						fmt.Println("========================================================")
						printJSON(order)
						printJSON(orderSide)
						fmt.Println("========================================================")
						var isReverseCalculate bool = false
						if order.ReverseCalculate != 0 {
							isReverseCalculate = true
							order.ReverseCalculate += int64(orderSide.Amount)
							if order.ReverseCalculate > 0 {

								orderSide.Amount -= uint64(order.ReverseCalculate)

								// if orderSide.FillOrKill && orderSide.Amount != 0 {
								// 	// calculate origin value

								// 	orderSide = orderBookSide[i]
								// 	if orderSide.Amount >= order.Amount {
								// 		orderSide.Amount -= order.Amount
								// 	} else {
								// 		order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
								// 		order.IDCalculate = orderSide.ID
								// 		orderSide.Amount = 0
								// 		order.AmountTemp = 0
								// 	}
								// 	// else {
								// 	// 	orderTemp = order
								// 	// 	orderTemp.Amount -= orderSide.Amount
								// 	// 	orderSide.Amount = 0
								// 	// }

								// 	orderSideAmount := orderSide.Amount
								// 	fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 175 ~ func ~ orderSideAmount`, orderSideAmount)

								// 	if orderSide.Amount != 0 {
								// 		orderSide = orderBookSide[i]
								// 		orderSide.AmountTemp = orderSideAmount
								// 		// orderSide.FillIndex = append([]int{i}, order.FillIndex...)
								// 		moreTrades, moreOrder := book.Process(orderSide)
								// 		if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
								// 			orderSide = moreOrder
								// 			orderSide.Amount = orderSide.AmountTemp
								// 			trades = append(trades, moreTrades...)
								// 		} else {
								// 			continue
								// 		}
								// 		// continue
								// 	}

								// 	// orderTemp.Amount = order.Amount
								// 	order.AmountTemp = orderTemp.Amount
								// 	if order.AmountTemp == 0 {
								// 		order.FillIndex = nil
								// 	} else {
								// 		order.AmountTemp = 0
								// 	}
								// }
								if orderSide.FillOrKill && orderSide.Amount != 0 {
									// calculate origin value

									orderSide = orderBookSide[i]
									if orderSide.Amount >= order.Amount {
										fmt.Println("asu", orderSide.Amount, order.Amount)
										orderSide.Amount -= order.Amount
									} else {
										order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
										order.IDCalculate = orderSide.ID
										orderSide.Amount = 0
										order.AmountTemp = 0
									}
									// else {
									// 	orderTemp = order
									// 	orderTemp.Amount -= orderSide.Amount
									// 	orderSide.Amount = 0
									// }

									orderSideAmount := orderSide.Amount

									if orderSide.Amount != 0 {
										orderSide = orderBookSide[i]
										orderSide.AmountTemp = orderSideAmount
										if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
											orderSide.FillIndex = append([]int{index}, order.FillIndex...)
										}
										fmt.Println("######################### orderside amount != 0 ######################### ")
										printJSON(orderSide)
										printJSON(order)
										fmt.Println("#########################                        ######################### ")
										moreTrades, moreOrder := book.Process(orderSide)
										if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
											orderSide = moreOrder
											orderSide.Amount = orderSide.AmountTemp
											trades = append(trades, moreTrades...)
										} else {
											continue
										}
										// continue
									}

									// orderTemp.Amount = order.Amount
									order.AmountTemp = orderTemp.Amount
									if order.AmountTemp == 0 {
										order.FillIndex = nil
									} else {
										order.AmountTemp = 0
									}
								}
								fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 153 ~ func ~ orderSide.Amount`, orderSide.Amount, order.ReverseCalculate)
								order.ReverseCalculate = 0
							}
						}

						if order.ReverseCalculate == 0 {
							trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderTemp.Amount, orderSide.Price, time.Now().String()})

							orderBookSide[i] = orderSide

							book.updateOrderBook(order.Side, orderBookSide)
						}

						// if order.Side == 1 {
						// 	book.SellOrders = orderBookSide
						// } else {
						// 	book.BuyOrders = orderBookSide
						// }

						if order.ReverseCalculate != 0 && isReverseCalculate {
							// order.ReverseCalculate = 0
							// order.IDCalculate = ""
							return trades, order
						}

						if len(order.FillIndex) > 0 {
							if order.FillIndex[0] == i {
								book.removeOrder(i, order.Side)
								i--
								n--
							}
						}

						if orderSide.Amount == 0 { // full match
							book.removeOrder(i, order.Side)
							i--
							n--
						}

						return trades, order
					}

					if orderSide.Amount < orderTemp.Amount {
						orderTemp.Amount -= orderSide.Amount
						trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSide.Amount, orderSide.Price, time.Now().String()})
						order.AmountTemp = orderTemp.Amount
						order.FillIndex = append([]int{i}, order.FillIndex...)
						var moreTrades []Trade
						var moreOrder Order
						if order.AmountTemp > 0 {
							moreTrades, moreOrder = book.Process(order)
							order = moreOrder
							orderBookSide = book.orderBookTemp(order.Side)
						}
						if order.AmountTemp != 0 {
							order.FillIndex = order.FillIndex[1:]
							trades = nil
							// if order.AmountTemp >= orderSide.Amount {
							// 	break
							// } else {
							// 	orderSide = orderBookSide[i]
							// 	orderSide.Amount -= order.AmountTemp
							// }
						}

						var isReverseCalculate bool = false
						if order.ReverseCalculate != 0 {
							isReverseCalculate = true
							// if order.ReverseCalculate >= orderSide.Amount {
							// 	order.ReverseCalculate -= orderSide.Amount
							// } else {
							// 	orderSide.Amount -= order.ReverseCalculate
							// 	order.ReverseCalculate = 0
							// }
							order.ReverseCalculate += int64(orderSide.Amount)
							// if order.ReverseCalculate > 0 {
							// 	orderSide.Amount -= uint64(order.ReverseCalculate)
							// 	order.ReverseCalculate = 0
							// }
							if order.ReverseCalculate > 0 {

								orderSide.Amount -= uint64(order.ReverseCalculate)

								if orderSide.FillOrKill && orderSide.Amount != 0 {
									// calculate origin value

									orderSide = orderBookSide[i]
									if orderSide.Amount >= order.Amount {
										orderSide.Amount -= order.Amount
									} else {
										order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
										order.IDCalculate = orderSide.ID
										orderSide.Amount = 0
										order.AmountTemp = 0
									}
									// else {
									// 	orderTemp = order
									// 	orderTemp.Amount -= orderSide.Amount
									// 	orderSide.Amount = 0
									// }

									orderSideAmount := orderSide.Amount
									fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 175 ~ func ~ orderSideAmount`, orderSideAmount)

									if orderSide.Amount != 0 {
										orderSide = orderBookSide[i]
										orderSide.AmountTemp = orderSideAmount
										// orderSide.FillIndex = append([]int{i}, order.FillIndex...)
										moreTrades, moreOrder := book.Process(orderSide)
										if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
											orderSide = moreOrder
											orderSide.Amount = orderSide.AmountTemp
											trades = append(trades, moreTrades...)
										} else {
											continue
										}
										// continue
									}

									// orderTemp.Amount = order.Amount
									order.AmountTemp = orderTemp.Amount
									if order.AmountTemp == 0 {
										order.FillIndex = nil
									} else {
										order.AmountTemp = 0
									}
								}
								fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 153 ~ func ~ orderSide.Amount`, orderSide.Amount, order.ReverseCalculate)
								order.ReverseCalculate = 0
							}
						}

						if order.ReverseCalculate == 0 {
							trades = append(trades, moreTrades...)

							orderBookSide[i] = orderSide

							book.updateOrderBook(order.Side, orderBookSide)
						}

						// if order.Side == 1 {
						// 	book.SellOrders = orderBookSide
						// } else {
						// 	book.BuyOrders = orderBookSide
						// }

						if len(order.FillIndex) > 0 && order.ReverseCalculate == 0 && !isReverseCalculate {
							if order.FillIndex[0] == i {
								book.removeOrder(i, order.Side)
								i--
								n--
							}
						}

						if len(trades) > 0 {
							order.FillIndex = order.FillIndex[1:]
						} else {
							break
						}

						return trades, order
					}
					break
				}

				// fill the entire order
				if orderSide.Amount >= orderTemp.Amount {
					orderSide.Amount -= orderTemp.Amount
					if orderSide.FillOrKill && orderSide.Amount != 0 {
						// calculate origin value

						orderSide = orderBookSide[i]
						if orderSide.Amount >= order.Amount {
							orderSide.Amount -= order.Amount
						} else {
							order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
							orderSide.Amount = 0
							order.AmountTemp = 0
						}
						// else {
						// 	orderTemp = order
						// 	orderTemp.Amount -= orderSide.Amount
						// 	orderSide.Amount = 0
						// }

						orderSideAmount := orderSide.Amount

						if orderSide.Amount != 0 {
							// fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 91 ~ func ~ orderSide`, orderSide)
							orderSide = orderBookSide[i]
							orderSide.AmountTemp = orderSideAmount
							// orderSide.FillIndex = append([]int{i}, order.FillIndex...)
							moreTrades, moreOrder := book.Process(orderSide)
							fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 123 ~ func ~ moreTrades`, moreTrades)
							if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
								orderSide = moreOrder
								orderSide.Amount = orderSide.AmountTemp
								trades = append(trades, moreTrades...)
							} else {
								continue
							}
							// continue
						}

						// orderTemp.Amount = order.Amount
						order.AmountTemp = orderTemp.Amount
						if order.AmountTemp == 0 {
							order.FillIndex = nil
						} else {
							order.AmountTemp = 0
						}
					} else {
						order.AmountTemp -= orderTemp.Amount
					}
					trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderTemp.Amount, orderSide.Price, time.Now().String()})
					// order = orderTemp
					// order.AmountTemp -= orderTemp.Amount

					orderBookSide[i] = orderSide

					book.updateOrderBook(order.Side, orderBookSide)

					if orderSide.Amount == 0 {
						book.removeOrder(i, order.Side)
						i--
						n--
					}
					return trades, order
				}
				// fill a partial order and continue
				log.Println(61, "\t", i, "\t", order.Amount, order.AmountTemp, orderSide.Amount)
				if orderSide.Amount < order.Amount {
					orderTemp.Amount -= orderSide.Amount
					trades = append(trades, Trade{order.ID, orderSide.ID, orderSide.Amount, orderSide.Price, time.Now().String()})
					order = orderTemp
					order.AmountTemp = orderTemp.Amount
					book.removeOrder(i, order.Side)
					// orderBookSide = append(orderBookSide[:i], orderBookSide[i+1:]...)
					orderBookSide = book.orderBookTemp(order.Side)
					i--
					n--
					continue
				}

			}
		}
	}
	// finally add the remaining order to the list
	if len(order.FillIndex) == 0 {
		var orderBookAdd []Order
		if order.Side == 1 {
			orderBookAdd = book.BuyOrders
		} else {
			orderBookAdd = book.SellOrders
		}

		if !book.contains(orderBookAdd, order.ID) {
			book.addOrder(order)
		}
	}
	return trades, order
}
