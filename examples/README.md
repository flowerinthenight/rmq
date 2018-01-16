```bash
# you have to set the following environment variables with your own values
$ export RABBITMQ_HOST=localhost
$ export RABBITMQ_HOST=5672
$ export RABBITMQ_HOST=guest
$ export RABBITMQ_HOST=guest

# build and run the receiver
$ cd recvr
$ go build -v
$ ./recvr

# then build the sender
$ cd sender
$ go build -v
$ ./sender
```
