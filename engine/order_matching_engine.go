package engine

import (
	"encoding/json"
	"fmt"
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
			bestPrice = orderBookSide[0].Price <= orderTemp.Price
		} else {
			bestPrice = orderBookSide[0].Price >= orderTemp.Price
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
					if orderSide.Amount >= orderTemp.Amount {
						orderSide.Amount -= orderTemp.Amount
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

							orderSideAmount := orderSide.Amount

							if orderSide.Amount != 0 {
								orderSide = orderBookSide[i]
								orderSide.AmountTemp = orderSideAmount
								if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
									orderSide.FillIndex = append([]int{index}, order.FillIndex...)
								}
								moreTrades, moreOrder := book.Process(orderSide)
								if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 && moreOrder.ReverseCalculate == 0 {
									orderSide = moreOrder
									orderSide.Amount = orderSide.AmountTemp
									trades = append(trades, moreTrades...)
								} else {
									// fmt.Println("111 ######################################################################### ")
									if moreOrder.ReverseCalculate < 0 && len(moreOrder.FillIndex) > 0 {
										// order.ReverseCalculate = -1 * moreOrder.ReverseCalculate
										// order.IDCalculate = moreOrder.ID
										if orderMaker, ok := book.getOrderByID(moreOrder.IDCalculate, moreOrder.Side); ok {
											trades = append(trades, Trade{moreOrder.ID, orderMaker.ID, orderMaker.Amount, orderMaker.Price, time.Now().String()})

											book.removeByID(orderMaker.ID, moreOrder.Side)

											orderSide.Amount -= orderMaker.Amount

											orderSide.Amount -= orderTemp.Amount

											trades = append(trades, Trade{order.ID, orderSide.ID, orderTemp.Amount, order.Price, time.Now().String()})

											book.removeOrder(i, order.Side)

											orderTemp.Amount = 0

											order.AmountTemp = orderTemp.Amount

											// order.FillReverse = append(order.FillReverse, FillReverse{
											// 	ID: orderMaker.ID,
											// })

											return trades, order
										}
										// trades = append(trades)
										printJSON(order)
										printJSON(moreOrder)
										printJSON(orderSide)
									}
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

						var isReverseCalculate bool = false
						if order.ReverseCalculate != 0 && order.IDCalculate != orderSide.ID {
							// println("masuk pak")
							isReverseCalculate = true
							order.ReverseCalculate += int64(orderSide.Amount)
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

									if orderSide.Amount != 0 {
										orderSide = orderBookSide[i]
										orderSide.AmountTemp = orderSideAmount
										if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
											orderSide.FillIndex = append([]int{index}, order.FillIndex...)
										}
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
								order.ReverseCalculate = 0
							}
						}

						if order.ReverseCalculate != 0 && !isReverseCalculate {
							// order.ReverseCalculate = 0
							// order.IDCalculate = ""

							// fmt.Println("RENE O PAK")

							// printJSON(order)
							// orderSideTemp := orderBookSide[i]
							// printJSON(orderSide)
							// printJSON(orderSideTemp)

							// if orderSide.FillOrKill {

							// 	orderSideTemp := orderBookSide[i]
							// 	trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSideTemp.Amount, orderSide.Price, time.Now().String()})

							// 	orderBookSide[i] = orderSide

							// 	book.updateOrderBook(order.Side, orderBookSide)
							// }

							return trades, order
						}

						if order.ReverseCalculate == 0 {

							// orderSideTemp := orderBookSide[i]
							trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderTemp.Amount, orderSide.Price, time.Now().String()})

							orderBookSide = book.orderBookTemp(order.Side)
							orderBookSide[i] = orderSide

							book.updateOrderBook(order.Side, orderBookSide)
						} else if !isReverseCalculate && order.IDCalculate != orderSide.ID {

							orderSideTemp := orderBookSide[i]
							trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSideTemp.Amount, orderSide.Price, time.Now().String()})

							orderBookSide[i] = orderSide

							book.updateOrderBook(order.Side, orderBookSide)
						}

						if len(order.FillIndex) > 0 {
							if order.FillIndex[0] == i {
								book.removeOrder(i, order.Side)
								i--
								n--
							}
						}

						if orderSide.Amount == 0 && !isReverseCalculate && orderSide.ID != order.IDCalculate { // full match
							book.removeOrder(i, order.Side)
							i--
							n--
						}

						return trades, order
					}

					if orderSide.Amount < orderTemp.Amount {
						orderTemp.Amount -= orderSide.Amount
						order.AmountTemp = orderTemp.Amount
						order.FillIndex = append([]int{i}, order.FillIndex...)
						var moreTrades []Trade
						var moreOrder Order
						if order.AmountTemp > 0 {
							moreTrades, moreOrder = book.Process(order)
							// printJSON(moreTrades)
							// printJSON(moreOrder)
							// fmt.Println("######################################################################### ")
							order = moreOrder
							trades = append(trades, moreTrades...)
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
							orderReverseCalculateTemp := order.ReverseCalculate
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
									// if orderSide.Amount >= order.Amount {
									// 	orderSide.Amount -= order.Amount
									// } else {
									// 	order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
									// 	fmt.Println(`🚀 ~ file: order_matching_engine.go ~ line 278 ~ func ~ int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)`, int64(orderTemp.Amount), int64(orderBookSide[i].Amount), order.Amount)
									// 	fmt.Println(`🚀 ~ file: order_matching_engine.go ~ line 278 ~ func ~ order.ReverseCalculate`, order.ReverseCalculate)
									// 	order.IDCalculate = orderSide.ID
									// 	orderSide.Amount = 0
									// 	order.AmountTemp = 0
									// }
									// // else {
									// // 	orderTemp = order
									// // 	orderTemp.Amount -= orderSide.Amount
									// // 	orderSide.Amount = 0
									// // }

									// orderSideAmount := orderSide.Amount

									// if orderSide.Amount != 0 {
									// 	orderSide = orderBookSide[i]
									// 	orderSide.AmountTemp = orderSideAmount
									// 	printJSON(orderBookSide)
									// 	// orderSide.FillIndex = append([]int{i}, order.FillIndex...)
									// 	moreTrades, moreOrder := book.Process(orderSide)
									// 	if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 {
									// 		orderSide = moreOrder
									// 		orderSide.Amount = orderSide.AmountTemp
									// 		trades = append(trades, moreTrades...)
									// 	} else {
									// 		continue
									// 	}
									// 	// continue
									// }

									// // orderTemp.Amount = order.Amount
									// order.AmountTemp = orderTemp.Amount
									// if order.AmountTemp == 0 {
									// 	order.FillIndex = nil
									// } else {
									// 	order.AmountTemp = 0
									// }
									trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSide.Amount, orderSide.Price, time.Now().String()})
									book.removeByID(orderSide.ID, order.Side)
									order.ReverseCalculate = orderReverseCalculateTemp
								} else {
									trades = append(trades, Trade{orderTemp.ID, orderSide.ID, uint64(order.ReverseCalculate), orderSide.Price, time.Now().String()})
									order.ReverseCalculate = 0
								}
								index, ok := book.getIndexByID(order.IDCalculate, order.Side)
								if ok {
									reverseCalculateOrder := orderBookSide[index]
									trades = append(trades, Trade{orderTemp.ID, reverseCalculateOrder.ID, reverseCalculateOrder.Amount, reverseCalculateOrder.Price, time.Now().String()})

									book.removeOrder(index, order.Side)

									orderBookSide = book.orderBookTemp(order.Side)

								}
							} else {
								if len(order.FillIndex) > 0 {
									order.FillIndex = order.FillIndex[1:]
								}
							}
						} else {
							trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSide.Amount, orderSide.Price, time.Now().String()})
						}

						if order.ReverseCalculate == 0 {
							// trades = append(trades, moreTrades...)

							orderBookSide[i] = orderSide

							book.updateOrderBook(order.Side, orderBookSide)
						} else {
							orderSide.Amount = orderBookSide[i].Amount
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

						if len(trades) > 0 && len(order.FillIndex) > 0 {
							order.FillIndex = order.FillIndex[1:]
						} else {
							trades = nil
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
							order.IDCalculate = orderSide.ID
							orderSide.Amount = 0
							order.AmountTemp = 0
						}

						orderSideAmount := orderSide.Amount

						if orderSide.Amount != 0 {
							orderSide = orderBookSide[i]
							orderSide.AmountTemp = orderSideAmount
							if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
								orderSide.FillIndex = append([]int{index}, order.FillIndex...)
							}
							moreTrades, moreOrder := book.Process(orderSide)
							if len(moreTrades) > 0 && moreOrder.AmountTemp == 0 && moreOrder.ReverseCalculate == 0 {
								orderSide = moreOrder
								orderSide.Amount = orderSide.AmountTemp
								trades = append(trades, moreTrades...)
							} else {
								// fmt.Println("111 ######################################################################### ")
								if moreOrder.ReverseCalculate < 0 && len(moreOrder.FillIndex) > 0 {
									// order.ReverseCalculate = -1 * moreOrder.ReverseCalculate
									// order.IDCalculate = moreOrder.ID
									if orderMaker, ok := book.getOrderByID(moreOrder.IDCalculate, moreOrder.Side); ok {
										trades = append(trades, Trade{moreOrder.ID, orderMaker.ID, orderMaker.Amount, orderMaker.Price, time.Now().String()})

										book.removeByID(orderMaker.ID, moreOrder.Side)

										orderSide.Amount -= orderMaker.Amount

										orderSide.Amount -= orderTemp.Amount

										trades = append(trades, Trade{order.ID, orderSide.ID, orderTemp.Amount, order.Price, time.Now().String()})

										book.removeOrder(i, order.Side)

										orderTemp.Amount = 0

										order.AmountTemp = orderTemp.Amount

										// order.FillReverse = append(order.FillReverse, FillReverse{
										// 	ID: orderMaker.ID,
										// })

										return trades, order
									}
									// trades = append(trades)
									printJSON(order)
									printJSON(moreOrder)
									printJSON(orderSide)
								}
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

					var isReverseCalculate bool = false
					if order.ReverseCalculate != 0 && order.IDCalculate != orderSide.ID {
						// println("masuk pak")
						isReverseCalculate = true
						order.ReverseCalculate += int64(orderSide.Amount)
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

								if orderSide.Amount != 0 {
									orderSide = orderBookSide[i]
									orderSide.AmountTemp = orderSideAmount
									if index, ok := book.getIndexByID(order.ID, orderSide.Side); ok {
										orderSide.FillIndex = append([]int{index}, order.FillIndex...)
									}
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
							order.ReverseCalculate = 0
						}
					}

					if order.ReverseCalculate != 0 && !isReverseCalculate {
						// order.ReverseCalculate = 0
						// order.IDCalculate = ""

						// fmt.Println("RENE O PAK")

						// printJSON(order)
						// orderSideTemp := orderBookSide[i]
						// printJSON(orderSide)
						// printJSON(orderSideTemp)

						// if orderSide.FillOrKill {

						// 	orderSideTemp := orderBookSide[i]
						// 	trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSideTemp.Amount, orderSide.Price, time.Now().String()})

						// 	orderBookSide[i] = orderSide

						// 	book.updateOrderBook(order.Side, orderBookSide)
						// }

						return trades, order
					}

					if order.ReverseCalculate == 0 {

						// orderSideTemp := orderBookSide[i]
						trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderTemp.Amount, orderSide.Price, time.Now().String()})

						orderBookSide = book.orderBookTemp(order.Side)
						orderBookSide[i] = orderSide

						book.updateOrderBook(order.Side, orderBookSide)
					} else if !isReverseCalculate && order.IDCalculate != orderSide.ID {

						orderSideTemp := orderBookSide[i]
						trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderSideTemp.Amount, orderSide.Price, time.Now().String()})

						orderBookSide[i] = orderSide

						book.updateOrderBook(order.Side, orderBookSide)
					}

					if len(order.FillIndex) > 0 {
						if order.FillIndex[0] == i {
							book.removeOrder(i, order.Side)
							i--
							n--
						}
					}

					if orderSide.Amount == 0 && !isReverseCalculate && orderSide.ID != order.IDCalculate { // full match
						book.removeOrder(i, order.Side)
						i--
						n--
					}
					// if orderSide.FillOrKill && orderSide.Amount != 0 {
					// 	// calculate origin value

					// 	orderSide = orderBookSide[i]
					// 	if orderSide.Amount >= order.Amount {
					// 		orderSide.Amount -= order.Amount
					// 	} else {
					// 		order.ReverseCalculate = int64(orderTemp.Amount) - int64(orderBookSide[i].Amount)
					// 		orderSide.Amount = 0
					// 		order.AmountTemp = 0
					// 	}
					// 	// else {
					// 	// 	orderTemp = order
					// 	// 	orderTemp.Amount -= orderSide.Amount
					// 	// 	orderSide.Amount = 0
					// 	// }

					// 	orderSideAmount := orderSide.Amount

					// 	if orderSide.Amount != 0 {
					// 		// fmt.Println(`🚀 ~ file: order_book_limit_order.go ~ line 91 ~ func ~ orderSide`, orderSide)
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
					// } else {
					// 	order.AmountTemp -= orderTemp.Amount
					// }
					// trades = append(trades, Trade{orderTemp.ID, orderSide.ID, orderTemp.Amount, orderSide.Price, time.Now().String()})

					// orderBookSide[i] = orderSide

					// book.updateOrderBook(order.Side, orderBookSide)

					// if orderSide.Amount == 0 {
					// 	book.removeOrder(i, order.Side)
					// 	i--
					// 	n--
					// }
					return trades, order
				}
				// fill a partial order and continue
				if orderSide.Amount < orderTemp.Amount {
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
			if order.AmountTemp > 0 {
				order.AmountTemp = 0
			}
			book.addOrder(order)
		}
	}

	return trades, order
}
