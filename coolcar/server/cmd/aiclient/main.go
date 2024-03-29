package main

import (
	"context"
	"coolcar/car/mq/amqpclt"
	coolenvpb "coolcar/shared/coolenv"
	"coolcar/shared/server"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:18001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	ac := coolenvpb.NewAIServiceClient(conn)
	c := context.Background()
	// 计算距离
	res, err := ac.MeasureDistance(c, &coolenvpb.MeasureDistanceRequest{
		From: &coolenvpb.Location{
			Latitude:  30,
			Longitude: 120,
		},
		To: &coolenvpb.Location{
			Latitude:  31,
			Longitude: 121,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)

	// Licsense recognition
	idRes, err := ac.LicIdentity(c, &coolenvpb.IdentityRequest{
		Photo: []byte{1, 2, 3, 4, 5},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", idRes)

	// Car position simulation
	_, err = ac.SimulateCarPos(c, &coolenvpb.SimulateCarPosRequest{
		CarId: "car123",
		InitialPos: &coolenvpb.Location{
			Latitude:  30,
			Longitude: 120,
		},
		Type: coolenvpb.PosType_RANDOM,
	})
	if err != nil {
		fmt.Println("SimulateCarPos error")
		panic(err)
	}

	logger, err := server.NewZapLogger()
	if err != nil {
		panic(err)
	}

	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	sub, err := amqpclt.NewSubscriber(amqpConn, "pos_sim", logger)
	if err != nil {
		panic(err)
	}

	ch, cleanUp, err := sub.SubscribeRaw(c)
	defer cleanUp()

	if err != nil {
		panic(err)
	}

	tm := time.After(10 * time.Second)
	for {
		shouldStop := false
		select {
		case msg := <-ch:
			fmt.Printf("%s\n", msg.Body)
		case <-tm:
			shouldStop = true
		}
		if shouldStop {
			break
		}
	}

	_, err = ac.EndSimulateCarPos(c, &coolenvpb.EndSimulateCarPosRequest{
		CarId: "car123",
	})
	if err != nil {
		panic(err)
	}
}
