package main

import (
	"fmt"
	"gopkg.in/redis.v3"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	// pubsub, err := client.Subscribe("test-channel")
	// if err != nil {
	// fmt.Println(err)
	// return
	// }
	// defer pubsub.Close()

	// go func() {
	// var err error
	for {
		err = client.Publish("test-channel", "hello").Err()
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second * 1)
	}
	// }()

	// for {
	// msg, err := pubsub.ReceiveMessage()
	// if err != nil {
	// fmt.Println(err)
	// return
	// }
	// fmt.Println("msg channel :", msg.Channel, " message data :", msg.Payload)
	// }
}
