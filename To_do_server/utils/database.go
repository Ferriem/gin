package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type key_value_pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Add(id string, title string, description string) error {
	ctx := context.Background()
	key_value := key_value_pair{
		Key:   id,
		Value: title,
	}
	err := rdb.RPush(ctx, id, key_value).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFirst(client *redis.Client, id string) (string, error) {
	ctx := context.Background()
	value, err := client.LRange(ctx, id, 0, 0).Result()
	if err != nil {
		return "", err
	}
	return value[0], nil
}

func GetInfo(client *redis.Client, id string) ([]string, error) {
	ctx := context.Background()
	value, err := client.LRange(ctx, id, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return value, nil
}

func Delete(client *redis.Client, id string) error {
	ctx := context.Background()
	err := client.Del(ctx, id).Err()
	if err != nil {
		return err
	}
	return nil
}

func Update(client *redis.Client, id string, title string, description string) error {
	ctx := context.Background()
	key_value := key_value_pair{
		Key:   id,
		Value: title,
	}
	err := client.LSet(ctx, id, 0, key_value).Err()
	if err != nil {
		return err
	}
	return nil
}
