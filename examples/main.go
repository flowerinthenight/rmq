package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowerinthenight/rmq"
)

func main() {
	// Setup CTRL+C handler for app termination.
	term := make(chan int)
	go handleSignal(false, func(s os.Signal) {
		term <- 1
	})

	// The usual RabbitMQ defaults.
	b := rmq.New(
		&rmq.Config{
			Host:     "localhost",
			Port:     5672,
			Username: "guest",
			Password: "guest",
			Vhost:    "/",
		},
		log.New(os.Stderr, "RMQ-[LOG] ", log.Lmicroseconds))

	err := b.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer b.Close()

	// Create a binding for exchange 'test' and queue 'qtest1'. We are providing all
	// the options here so it can send and consume messages at the same time.
	// The return string is the binding id.
	bind1, err := b.AddBinding(
		&rmq.ExchangeOptions{
			Name:    "test",
			Type:    "direct",
			Durable: true,
		},
		&rmq.QueueOptions{
			QueueName: "qtest1",
			Durable:   true,
		},
		&rmq.QueueBindOptions{
			RoutingKey: "rk1",
		},
		&rmq.ConsumeOptions{
			ClientTag: "consumer1",
			FnCallback: func(b []byte) error {
				log.Println(fmt.Sprintf("[qtest1] payload: %s", b))
				return nil
			},
		})

	if err != nil {
		log.Fatalln(err)
	}

	// Create a binding for exchange 'test' and a queue with auto-generated name. We are
	// also providing all the options here so it can send and consume messages at the same time.
	// The return string is the binding id.
	bind2, err := b.AddBinding(
		&rmq.ExchangeOptions{
			Name:    "test",
			Type:    "direct",
			Durable: true,
		},
		&rmq.QueueOptions{
			QueueName: "", // auto-generate queue name
		},
		&rmq.QueueBindOptions{
			RoutingKey: "rk2",
		},
		&rmq.ConsumeOptions{
			ClientTag: "consumer2",
			FnCallback: func(b []byte) error {
				log.Println(fmt.Sprintf("[qtest2] payload: %s", b))
				return nil
			},
		})

	if err != nil {
		log.Fatalln(err)
	}

	// Fire a goroutine that alternately sends a message to 'bind1' and 'bind2'.
	go func() {
		flip := true
		for {
			if flip {
				b.Send(bind1, "rk1", false, false, []byte(fmt.Sprintf("for qtest1: %s", time.Now().String())))
			} else {
				b.Send(bind2, "rk2", false, false, []byte(fmt.Sprintf("for autogen queue: %s", time.Now().String())))
			}

			time.Sleep(1 * time.Second)
			flip = !flip
		}
	}()

	<-term
	os.Exit(0)
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
