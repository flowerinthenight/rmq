package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	// Create a binding for exchange 'test' and queue 'qtest1'.
	// We don't define ConsumeOpt since we only send messages.
	// The return string is the binding id.
	bindId, err := b.AddBinding(&rmq.BindConfig{
		ExchangeOpt: &rmq.ExchangeOptions{
			Name:       "test-exchange",
			Type:       "direct",
			Durable:    false,
			AutoDelete: true,
		},
		QueueOpt: &rmq.QueueOptions{
			QueueName:  "queue1",
			Durable:    false,
			AutoDelete: true,
		},
		QueueBindOpt: &rmq.QueueBindOptions{
			RoutingKey: "rk1",
		},
	})

	log.Printf("binding added (id = %v)", bindId)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			b.Send(
				bindId, // binding id to send to
				"rk1",  // route key
				[]byte(fmt.Sprintf("for qtest1: %s", time.Now().String())), // message
			)

			time.Sleep(1 * time.Second)
		}
	}()

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
