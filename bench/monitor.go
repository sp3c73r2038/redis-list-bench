package bench

import (
	"log"
	"time"
)

func Monitor(q chan int) {
	var cnt int64 = 0

	go func(c *int64) {
		var current int64
		var prev int64
		var delta int64
		for {
			current = *c
			delta = current - prev
			log.Printf("cnt: %d\n", delta)
			prev = current
			time.Sleep(1 * time.Second)
		}
	}(&cnt)

	for _ = range q {
		cnt++
	}
}
