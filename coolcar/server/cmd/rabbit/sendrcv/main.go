package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel() // 虚拟的connection
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"go_q1", // 名字
		true,    // durable
		false,   // outoDelete
		false,   // exclude
		false,   //  noWait
		nil,     // args
	)
	if err != nil {
		panic(err)
	}

	go consume("c1", conn, q.Name)
	go consume("c2", conn, q.Name)

	i := 0
	for {
		i++
		err := ch.Publish(
			"",     //exchange
			q.Name, //key
			false,  // mandatory
			false,  // immediate
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

func consume(consumer string, conn *amqp.Connection, q string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
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
