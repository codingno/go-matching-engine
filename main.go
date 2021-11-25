package main

import (
	"encoding/json"
	"fmt"
	"log"
	"matching-engine/cache"
	"matching-engine/engine"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

func main() {

	// create the consumer and listen for new order messages
	// consumer := createConsumer()

	// workerConn := createWorker()
	workerConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer workerConn.Close()
	workerChannel, err := workerConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer workerChannel.Close()
	q, err := workerChannel.QueueDeclare(
		"order", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = workerChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := workerChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// create the producer of trade messages
	// producer := createProducer()

	// queue, , _ := createPublisher()
	publisherConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer publisherConn.Close()

	publisherChannel, err := publisherConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer publisherChannel.Close()

	queue, err := publisherChannel.QueueDeclare(
		"trade", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = publisherChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}

	getOrderBookCache, err := cache.Get("orderbook")
	if err == redis.Nil {
		cache.Set("orderbook", string(book.ToJSON()))
		getOrderBookCache, _ = cache.Get("orderbook")
	}
	book.FromJSON([]byte(getOrderBookCache))
	// printOrderBook(book)
	fmt.Println("########################################################################################")
	fmt.Println("################################ After Trade ###########################################")
	fmt.Println("########################################################################################")
	fmt.Println("# ID\tQueue\tF/K\tAmount\tBid\t||   Ask\tAmount\tF/K\tQueue\tID     #")
	fmt.Println("########################################################################################")
	for i := 0; i < len(book.SellOrders); i++ {
		n := len(book.SellOrders) - 1 - i
		sellOrder := book.SellOrders[n]
		fmt.Println("# .....\t......\t......\t......\t......\t||  ", sellOrder.Price, "\t", sellOrder.Amount, "\t", sellOrder.FillOrKill, "\t", n, "\t", sellOrder.ID[len(sellOrder.ID)-4:], " #")
	}
	for i, buyerOrder := range book.BuyOrders {
		fmt.Println("#", buyerOrder.ID[len(buyerOrder.ID)-4:], "\t", i, "\t", buyerOrder.FillOrKill, "\t", buyerOrder.Amount, "\t", buyerOrder.Price, "\t||  ......\t......\t......\t......\t.....  #")
	}
	fmt.Println("########################################################################################")
	fmt.Println("buyer order len : ", len(book.BuyOrders), "sell order len : ", len(book.SellOrders))

	// create a signal channel to know when we are done
	done := make(chan bool)

	go func() {
		for d := range msgs {
			var order engine.Order
			// decode the message
			order.FromJSON(d.Body)
			printJSON(order)
			// process the order
			// log.Println("##########################################")
			// log.Println("############## Before Trade ##############")
			// log.Println("##########################################")
			// log.Println("  Amount\tBid\t||\tAsk\tAmount")
			// log.Println("##########################################")
			// for i := 0; i < len(book.SellOrders); i++ {
			// 	sellOrder := book.SellOrders[len(book.SellOrders)-1-i]
			// 	log.Println("  ......\t\t||\t", sellOrder.Price, "\t", sellOrder.Amount)
			// }
			// for _, buyerOrder := range book.BuyOrders {
			// 	log.Println("  ", buyerOrder.Amount, "\t", buyerOrder.Price, "\t||\t\t......")
			// }
			// log.Println("##########################################")
			trades, _ := book.Process(order)
			cache.Set("orderbook", string(book.ToJSON()))
			// fmt.Println(`🚀 ~ file: main.go ~ line 97 ~ gofunc ~ trades`, trades)
			// log.Println("orderbook : ", book)
			// send trades to message queue
			for _, trade := range trades {
				rawTrade := trade.ToJSON()
				err := publisherChannel.Publish(
					"",         // exchange
					queue.Name, // routing key
					false,      // mandatory
					false,      // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(rawTrade),
					})
				failOnError(err, "Failed to publish a message")
			}
			// printOrderBook(book)
			fmt.Println("########################################################################################")
			fmt.Println("################################ After Trade ###########################################")
			fmt.Println("########################################################################################")
			fmt.Println("# ID\tQueue\tF/K\tAmount\tBid\t||\tAsk\tAmount\tF/K\tQueue\tID     #")
			fmt.Println("########################################################################################")
			for i := 0; i < len(book.SellOrders); i++ {
				n := len(book.SellOrders) - 1 - i
				sellOrder := book.SellOrders[n]
				fmt.Println("# .....\t......\t......\t......\t......\t||\t", sellOrder.Price, "\t", sellOrder.Amount, "\t", sellOrder.FillOrKill, "\t", n, "\t", sellOrder.ID[len(sellOrder.ID)-4:], " #")
			}
			for i, buyerOrder := range book.BuyOrders {
				fmt.Println("#", buyerOrder.ID[len(buyerOrder.ID)-4:], "\t", i, "\t", buyerOrder.FillOrKill, "\t", buyerOrder.Amount, "\t", buyerOrder.Price, "\t||\t......\t......\t......\t......\t.....  #")
			}
			fmt.Println("########################################################################################")
			fmt.Println("buyer order len : ", len(book.BuyOrders), "sell order len : ", len(book.SellOrders))
		}
	}()

	<-done
}
