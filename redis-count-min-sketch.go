package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	fmt.Printf("\n# COUNT-MIN\n")
	countMinSketch(ctx, rdb)

}
func countMinSketch(ctx context.Context, rdb *redis.Client) {
	if err := rdb.Do(ctx, "CMS.INITBYPROB", "count_min", 0.001, 0.01).Err(); err != nil {
		panic(err)
	}

	items := []string{"item1", "item2", "item3", "item4", "item5"}
	counts := make(map[string]int, len(items))

	for i := 0; i < 10000; i++ {
		n := rand.Intn(len(items))
		item := items[n]

		if err := rdb.Do(ctx, "CMS.INCRBY", "count_min", item, 1).Err(); err != nil {
			panic(err)
		}
		counts[item]++
	}

	for item, count := range counts {
		ns, err := rdb.Do(ctx, "CMS.QUERY", "count_min", item).Int64Slice()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: count-min=%d actual=%d\n", item, ns[0], count)
	}
}
