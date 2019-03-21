package main

import (
	"flag"
	"log"

	bench "github.com/aleiphoenix/redis-list-bench/bench"
)

var (
	consumer = flag.Bool("consumer", false, "act as consumer")
	producer = flag.Bool("producer", false, "act as producer")
	addr     = flag.String("redis", "127.0.0.1:6379", "redis server address")
	password = flag.String("password", "", "redis auth password")
	pipe     = flag.Int("pipe", 1, "pipe size")
	db       = flag.Int("db", 0, "redis db")
	key      = flag.String("key", "mylist", "key name")
	worker   = flag.Int("worker", 1, "worker connection & threads")
	msgsize  = flag.Int("msgsize", 3, "data size")
)

func main() {
	flag.Parse()

	option := new(bench.Option)
	option.Addr = *addr
	option.Pass = *password
	option.Worker = *worker
	option.Key = *key
	option.DB = *db
	option.Pipe = *pipe
	option.Msgsize = *msgsize

	if *consumer {
		bench.Consumer(option)
	}

	if *producer {
		bench.Producer(option)
	}

	log.Fatal("please select consumer or producer mode")
}
