package main

import (
	"log"
	"os"

	"github.com/flowerinthenight/rmq"
)

func main() {
	log.Println("hello")
	b := rmq.New(&rmq.Config{
		Host:     rmqhost,
		Port:     port,
		Username: rmquser,
		Password: rmqpass,
		Vhost:    "/",
	}, log.New(os.Stderr, "RMQ_[TEMP] ", log.Lmicroseconds))
}
