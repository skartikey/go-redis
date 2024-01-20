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

	fmt.Printf("# BLOOM\n")
	bloomFilter(ctx, rdb)

}

func bloomFilter(ctx context.Context, rdb *redis.Client) {
	inserted, err := rdb.Do(ctx, "BF.ADD", "bf_key", "item0").Bool()
	if err != nil {
		panic(err)
	}
	if inserted {
		fmt.Println("item0 was inserted")
	} else {
		fmt.Println("item0 already exists")
	}

	for _, item := range []string{"item0", "item1"} {
		exists, err := rdb.Do(ctx, "BF.EXISTS", "bf_key", item).Bool()
		if err != nil {
			panic(err)
		}
		if exists {
			fmt.Printf("%s does exist\n", item)
		} else {
			fmt.Printf("%s does not exist\n", item)
		}
	}

	bools, err := rdb.Do(ctx, "BF.MADD", "bf_key", "item1", "item2", "item3").BoolSlice()
	if err != nil {
		panic(err)
	}
	fmt.Println("adding multiple items:", bools)
}
