package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flowerinthenight/rmq"
)

func main() {
	term := make(chan int)
	go handleSignal(false, func(s os.Signal) {
		term <- 1
	})

	b := rmq.New(&rmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "jennah",
		Password: "jennah",
		Vhost:    "/",
	}, log.New(os.Stderr, "RMQ_[LOG] ", log.Lmicroseconds))

	err := b.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer b.Close()
	b.AddBinding(&rmq.ExchangeOptions{
		Name:    "test",
		Type:    "direct",
		Durable: true,
	}, &rmq.QueueOptions{
		QueueName: "qtest",
		Durable:   true,
	}, &rmq.QueueBindOptions{
		RoutingKey: "rk1",
	}, &rmq.ConsumeOptions{
		ClientTag: "consumer",
		FnCallback: func(b []byte) error {
			log.Println(fmt.Sprintf("[qtest] payload: %s", b))
			return nil
		},
	})

	<-term
}

func handleSignal(exit bool, callback func(s os.Signal)) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	exitChan := make(chan int)
	go func() {
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGHUP:
				log.Println("SIGHUP", s)
				callback(s)
			case syscall.SIGINT:
				log.Println("SIGINT", s)
				callback(s)
				exitChan <- 0
			case syscall.SIGTERM:
				log.Println("SIGTERM", s)
				callback(s)
				exitChan <- 0
			case syscall.SIGQUIT:
				log.Println("SIGQUIT", s)
				callback(s)
				exitChan <- 0
			default:
				log.Println("UNKNOWN", s)
				callback(s)
				exitChan <- 1
			}
		}
	}()

	code := <-exitChan
	if exit {
		os.Exit(code)
	}
}
