package bench

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func Consumer(option *Option) {

	log.Println(">>> redis:", option.Addr)
	log.Println(">>> pipe:", option.Pipe)
	log.Println(">>> worker:", option.Worker)

	q := make(chan int, 10000)
	go Monitor(q)

	lock := new(sync.WaitGroup)
	for i := 0; i < option.Worker; i++ {
		lock.Add(1)
		log.Println(">>>> start worker", i)
		go consumer(option, q)
	}
	lock.Wait()
}

func consumer(option *Option, q chan int) {
	opt := &redis.Options{
		Addr:     option.Addr,
		Password: option.Pass,
		DB:       option.DB,
	}

	batch := 1
	if option.Pipe > 1 {
		batch = option.Pipe
	}

	client := redis.NewClient(opt)
	defer client.Close()

	for {
		pipe := client.Pipeline()

		for i := 0; i < batch-1; i++ {
			pipe.LPop(option.Key)
		}
		pipe.BLPop(time.Second*1, option.Key)

		r, err := pipe.Exec()
		pipe.Close()

		for _, re := range r {
			switch re.(type) {
			case *redis.StringCmd:
				// log.Println(re.Err() == redis.Nil)
				if re.Err() != redis.Nil {
					// log.Println(re.(*redis.StringCmd).Val())
					q <- 1
				}
			case *redis.StringSliceCmd:
				// log.Println(re.Err() == redis.Nil)
				if re.Err() != redis.Nil {
					// log.Println(re.(*redis.StringSliceCmd).Val()[1])
					q <- 1
				}
			default:
			}
		}

		if err != nil {
			if err != redis.Nil {
				log.Println("error not nil")
				log.Println(err)
			}
		}

	}
}
