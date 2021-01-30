package gobatch_redis

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MasterOfBinary/gobatch-redis/batchhll"
	"github.com/MasterOfBinary/gobatch/batch"
	"github.com/MasterOfBinary/gobatch/source"
	"github.com/go-redis/redis/v8"
)

func ExampleHLL() {
	ctx := context.Background()

	key := "somekey"

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rand.Seed(101)

	hll := batchhll.New(rdb, key)
	ch := make(chan interface{})
	s := &source.Channel{
		Input: ch,
	}

	config := batch.NewConstantConfig(&batch.ConfigValues{
		MinItems: 10,
		MaxItems: 20,
		MaxTime:  5 * time.Millisecond,
	})
	b := batch.New(config)

	errs := b.Go(ctx, s, hll)

	go genData(ch)

	for err := range errs {
		panic(err)
	}

	res, err := rdb.PFCount(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Count of %v is %v", key, res)

	err = rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}

func genData(ch chan<- interface{}) {
	timer := time.NewTimer(1 * time.Second)

	var done bool
	for {
		select {
		case <-timer.C:
			done = true
			break
		default:
			ch <- rand.Int63n(10000000)
		}

		if done {
			break
		}
	}

	close(ch)
}
