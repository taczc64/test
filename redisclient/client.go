package main

import (
	"fmt"
	"gopkg.in/redis.v3"
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
	pubsub, err := client.Subscribe("eventmessage")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("msg channel :", msg.Channel, " message data :", msg.Payload)
	}
}
