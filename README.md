## kafka-toolkit
---

```
$ go get -v -u github.com/confluentinc/confluent-kafka-go/kafka
$ cd $GOPATH/src/kafka-toolkit
$ go build producer.go
$ go build consumer.go
```

### Help

```
➜  kafka-toolkit git:(master) ✗ ./consumer --help
Usage of ./consumer:
  -b string
    	Kafka Brokers (default "localhost")
  -g string
    	Group (default "test")
  -h	Help
  -s int
    	Session Timeout (default 60000)
  -t string
    	Topic (default "test")
➜  kafka-toolkit git:(master) ✗ ./consumer -
➜  kafka-toolkit git:(master) ✗ ./producer --help
Usage of ./producer:
  -b string
    	Kafka Brokers (default "localhost")
  -g string
    	Group (default "test")
  -h	Help
```
