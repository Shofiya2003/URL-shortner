package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func StoreUrl(newUrl Url) error {
	redisClient, err := GetClient()

	if err != nil {
		fmt.Println("error in storing Url")
		return err
	}

	ctx := context.Background()
	key := newUrl.Key

	marshalledCurrency, err := json.Marshal(newUrl)

	if err != nil {
		fmt.Println("marshall error: ", err)
		return err
	}

	return redisClient.MSet(ctx, key, marshalledCurrency).Err()

}

func FetchUrl(key string) (Url, error) {
	redisClient, err := GetClient()

	url := Url{}
	if err != nil {
		fmt.Println("error in fetching url")
		return url, err
	}

	ctx := context.Background()

	res, err := redisClient.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return url, nil
		}
		fmt.Println("error in retrieving url from redis")
		return url, err
	}

	err = json.Unmarshal([]byte(res), &url)

	if err != nil {
		fmt.Println("error in unmarshalling url")
		return url, err
	}

	return url, nil
}
