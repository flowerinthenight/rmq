package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/flowerinthenight/rmq"
)

func main() {
	// Setup CTRL+C handler for app termination.
	term := make(chan int)
	go handleSignal(true, func(s os.Signal) {
		term <- 0
	})

	port, err := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	if err != nil {
		log.Fatalln(err)
	}

	b := rmq.New(&rmq.Config{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     port,
		Username: os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASS"),
		Vhost:    "/",
	})

	err = b.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer b.Close()

	// Create a binding for exchange 'test' and queue 'qtest1'. We are providing all
	// the options here so it can send and consume messages at the same time.
	// The return string is the binding id.
	_, err = b.AddBinding(&rmq.BindConfig{
		&rmq.ExchangeOptions{
			Name:    "test",
			Type:    "direct",
			Durable: false,
		},
		&rmq.QueueOptions{
			QueueName: "qtest1",
			Durable:   false,
		},
		&rmq.QueueBindOptions{
			RoutingKey: "rk1",
		},
		&rmq.ConsumeOptions{
			ClientTag: "consumer1",
			FnCallback: func(b []byte) error {
				log.Printf("[qtest1] payload: %s", b)
				return nil
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

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
