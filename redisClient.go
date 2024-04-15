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

	redisClient.RPush(ctx, "urls", marshalledCurrency).Err()

	if err != nil {
		fmt.Println("error in adding to sorted set", err)
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

func FetchAllUrl(page int) ([]Url, error) {

	redisClient, err := GetClient()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	pageSize := 2
	start := (page - 1) * pageSize
	stop := start + pageSize - 1

	ctx := context.Background()

	result, err := redisClient.LRange(ctx, "urls", int64(start), int64(stop)).Result()

	if err != nil {
		fmt.Println("error in fetching values")
		fmt.Println(err)
		return nil, err
	}

	urls := make([]Url, len(result))

	for i, item := range result {

		var url Url

		json.Unmarshal([]byte(item), &url)

		urls[i] = url
	}

	return urls, nil
}

func DeleteUrl(key string) error {
	redisClient, err := GetClient()

	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := context.Background()

	_, err = redisClient.Del(ctx, key).Result()

	if err != nil {
		fmt.Println("error in deleting key from redis")
		return err
	}

	fmt.Printf("deleted key %s ", key)
	return nil
}
