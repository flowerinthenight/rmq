[![Build Status](https://travis-ci.org/flowerinthenight/rmq.svg?branch=master)](https://travis-ci.org/flowerinthenight/rmq)

# Overview

A simple wrapper to [streadway/amqp](https://github.com/streadway/amqp) for RabbitMQ with support for auto reconnections.

## Usage

The library maintains a single connection and channel. It also maintains a map of bindings of exchanges and queues added by the user. Each binding can be configured to be a producer, a consumer, or both.

First, create the connection object with:


```go
// replace the values below with your own
c := rmq.New(&rmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Vhost:    "/",
})

err := c.Connect()
...
```

Next, we will create a binding for a direct-type exchange `test-exchange` and queue `qtest` with:

```go
bindId, err := c.AddBinding(&rmq.BindConfig{
		ExchangeOpt: &rmq.ExchangeOptions{
			Name:       "test-exchange",
			Type:       "direct",
			Durable:    false,
			AutoDelete: true,
		},
		QueueOpt: &rmq.QueueOptions{
			QueueName:  "qtest",
			Durable:    false,
			AutoDelete: true,
		},
		QueueBindOpt: &rmq.QueueBindOptions{
			RoutingKey: "rk1",
		},
		// when `ConsumeOpt` is provided, this binding is able to receive
		// messages from the specified exchange/queue
		ConsumeOpt: &rmq.ConsumeOptions{
			ClientTag:  "consumer1",
			FnCallback: func (b []byte) error {
				log.Printf("payload: %s", b)
				return nil
			},
		},
})

// send a message using the binding above
c.Send(bindId, "rk1", []byte("hello world"))
```

You can also create a send-only binding with:

```go
bindId, err := c.AddBinding(&rmq.BindConfig{
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

c.Send(bindId, "rk1", []byte("hello world"))
```

See the [examples](./examples) directory for a simple receiver/sender implementation.

## License

[The MIT License](./LICENSE.md)
