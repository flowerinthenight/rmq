[![Build Status](https://travis-ci.org/flowerinthenight/rmq.svg?branch=master)](https://travis-ci.org/flowerinthenight/rmq)

# Overview

A simple wrapper to [streadway/amqp](https://github.com/streadway/amqp) for RabbitMQ with support for auto reconnections.

### Usage

The library maintains a single connection and channel. It also maintains a map of bindings of exchanges and queues added by the user. Each binding can be configured to be a producer, a consumer, or both.

The connection object can be initialized using the following code snippet:


```go
c := rmq.New(&rmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Vhost:    "/",
})
```

See the [examples](./examples) directory for a simple implementation.

### License

[The MIT License](./LICENSE.md)
