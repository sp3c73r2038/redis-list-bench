package bench

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func Producer(option *Option) {

	log.Println(">>> redis addr:", option.Addr)
	log.Println(">>> pipe:", option.Pipe)
	log.Println(">>> worker:", option.Worker)

	q := make(chan int, 10000)
	go Monitor(q)

	iq := make(chan int, 100000)
	for i := 0; i < option.Worker; i++ {
		go func(q chan int) {
			for {
				iq <- 1
			}
		}(iq)
	}

	b := []byte
	for i := 0; i < option.Msgsize; i++ {
		b[i] = 255
	}

	lock := new(sync.WaitGroup)
	for i := 0; i < option.Worker; i++ {
		log.Println(">>>> start worker ", i)
		lock.Add(1)
		go producer(b, option, iq, q)
	}

	lock.Wait()
}

func producer(b []byte, option *Option, iq chan int, q chan int) {

	opt := &redis.Options{
		Addr:     option.Addr,
		Password: option.Pass,
		DB:       option.DB,
	}

	client := redis.NewClient(opt)
	if option.Pipe <= 1 {
		log.Println("!! pipe = 1, normal send")
		normalSend(client, option.Key, b, iq, q)
	} else {
		log.Println("!! pipe > 1, group send")
		groupSend(client, option.Key, b, option.Pipe, iq, q)
	}

}

func normalSend(
	client *redis.Client, key string, b []byte, iq chan int, q chan int) {
	for _ = range iq {
		client.RPush(key, b)
		q <- 1
	}
}

func send(
	client *redis.Client, key string, msgs []([]byte), rq chan int) {
	var err error
	// log.Printf("send: %d\n", len(msgs))
	pipe := client.Pipeline()
	for _, b := range msgs {
		// pipe.RPush("1", buf[i])
		pipe.RPush(key, b)
		rq <- 1
	}
	_, err = pipe.Exec()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("sent")
}

func groupSend(
	client *redis.Client, key string, b []byte,
	size int, q chan int, rq chan int) {

	buf := make([]([]byte), 0, size)
	var cnt int = 0
	var timestamp = time.Now().UnixNano()

	for {
		select {
		case _, ok := <-q:
			// log.Printf("got: %d\n", ii)
			if !ok {
				q = nil
				break
			}
			cnt++
			buf = append(buf, b)
			if cnt >= size {
				send(client, key, buf, rq)
				buf = make([]([]byte), 0, size)
				cnt = 0
			}
			timestamp = time.Now().UnixNano()
		default:
			now := time.Now().UnixNano()
			// log.Printf("delta: %d\n", now-timestamp)
			// log.Printf("len: %d\n", len(buf))
			if (now - timestamp) > 1e9 {
				send(client, key, buf, rq)
				buf = make([]([]byte), 0, size)
				cnt = 0
			}
			timestamp = time.Now().UnixNano()
			time.Sleep(1 * time.Second)
		}

		if q == nil && len(buf) <= 0 {
			break
		}
	}
}
