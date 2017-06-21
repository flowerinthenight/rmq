[![Build Status](https://travis-ci.org/flowerinthenight/rmq.svg?branch=master)](https://travis-ci.org/flowerinthenight/rmq)

# Overview

A simple wrapper to [streadway/amqp](https://github.com/streadway/amqp) for RabbitMQ with support for auto reconnections.

### Usage

The library maintains a single connection and channel. It also maintains a map of bindings of exchanges and queues. Each binding can be configured to be a producer, a consumer, or both.

See the 'examples' directory for a simple implementation.

### License

[The MIT License](./LICENSE.md)
