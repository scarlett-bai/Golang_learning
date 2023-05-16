package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

const exchange = "go_ex"

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel() // 虚拟的connection
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		exchange, // 名字
		"fanout",
		true,  // durable
		false, // outoDelete
		false, // internla exchagne
		false, //  noWait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}

	go subscibe(conn, exchange)
	go subscibe(conn, exchange)

	i := 0
	for {
		i++
		err := ch.Publish(
			exchange, //exchange
			"",       //key   orange
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				Body: []byte(fmt.Sprintf("message %d", i)),
			},
		)
		if err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(200 * time.Millisecond)
	}

}

func subscibe(conn *amqp.Connection, ex string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // 名字 ,不设置 系统会自动分配
		false, // durable
		true,  // outoDelete
		false, // exclude
		false, //  noWait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}
	defer ch.QueueDelete(
		q.Name,
		false, // isUnused
		false, // isEmpty
		false, // noWait
	)
	err = ch.QueueBind(
		q.Name,
		"", //key
		ex,
		false, // noWait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}
	consume("c", ch, q.Name)
}

func consume(consumer string, ch *amqp.Channel, q string) {
	msgs, err := ch.Consume(
		q,        // name
		consumer, // 消费者的名字
		true,     // autoAck
		false,    //exclusive 独占这个queue
		false,    // noLocal
		false,    //noWait
		nil,      // args
	)
	if err != nil {
		panic(err)
	}
	for msg := range msgs {
		fmt.Printf("consumer:%s rec :%s\n", consumer, msg.Body)
	}

}
