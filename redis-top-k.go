package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	fmt.Printf("\n# TOP-K\n")
	topK(ctx, rdb)
}

func topK(ctx context.Context, rdb *redis.Client) {
	if err := rdb.Do(ctx, "TOPK.RESERVE", "top_items", 3).Err(); err != nil {
		panic(err)
	}

	counts := map[string]int{
		"item1": 1000,
		"item2": 2000,
		"item3": 3000,
		"item4": 4000,
		"item5": 5000,
		"item6": 6000,
	}

	for item, count := range counts {
		for i := 0; i < count; i++ {
			if err := rdb.Do(ctx, "TOPK.INCRBY", "top_items", item, 1).Err(); err != nil {
				panic(err)
			}
		}
	}

	items, err := rdb.Do(ctx, "TOPK.LIST", "top_items").StringSlice()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		ns, err := rdb.Do(ctx, "TOPK.COUNT", "top_items", item).Int64Slice()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: top-k=%d actual=%d\n", item, ns[0], counts[item])
	}
}
