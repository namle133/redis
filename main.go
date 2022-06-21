package main

import (
	"fmt"
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v5"
)

type Object struct {
	Str string
	Num int
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)
	setCache()

	if er := rdb.Set("name", "Nam", 0).Err(); er != nil {
		return
	}

	if e := rdb.Set("fullname", "LeThanhNam", 0).Err(); e != nil {
		return
	}

	value, err := rdb.Get("name").Result()

	if err != nil {
		return
	}

	fmt.Println("name is value: ", value)

}

func setCache() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})

	codec := &cache.Codec{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	key := "mykey"

	obj := &Object{
		Str: "mystring",
		Num: 42,
	}

	codec.Set(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: time.Minute,
	})

	var wanted Object
	if err := codec.Get("mykey", &wanted); err == nil {
		fmt.Println(wanted)
	} else {
		fmt.Println(err)
	}
}
