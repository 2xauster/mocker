package data

import (
	"context"
	"errors"

	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/redis/go-redis/v9"
)

var (
	MissingEntityNameError = errs.NewError(errors.New("missing entity name"), errs.RedisErrorType, errs.ErrDataIllegal)
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(connString string) (*Redis, error) {
	if connString == "" {
		return nil, errors.New("connection string is empty")
	}

	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, err 
	}

	client := redis.NewClient(opt)
	return &Redis{
		Client: client,
	}, nil
}

func HSet[Entity any](ctx context.Context, client *redis.Client, id string, entityData Entity) (int64, error) {
	meta := ExtractMeta(entityData, true)

	identifier := meta.Name + ":" + id 
	hashFields := make(map[string]interface{}, len(meta.Fields))
	
	for _, v := range meta.Fields {
		hashFields[v.Name] = v.Value
	}

	res, err := client.HSet(ctx, identifier, hashFields).Result()
	err = RedisErrorComparator(err)
	if err != nil {
		return res, err
	}

	return res, nil
}

func HGet[Entity any](ctx context.Context, client *redis.Client, id string, entity Entity) (map[string]string, error) {
	meta := ExtractMeta(entity, false)
	if meta.Name == "" {
		return nil, MissingEntityNameError
	}

	key := meta.Name + ":" + id

	data, err := client.HGetAll(ctx, key).Result()
	err = RedisErrorComparator(err)

	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, redis.Nil
	}

	return data, nil
}

func HUpdate[Entity any](ctx context.Context, client *redis.Client, id string, updates Entity) (int64, error) {
	meta := ExtractMeta(updates, false)
	if meta.Name == "" {
		return 0, MissingEntityNameError
	}

	key := meta.Name + ":" + id

	fields := make(map[string]interface{}, len(meta.Fields))
	for _, f := range meta.Fields {
		if f.Value != nil {
			fields[f.Name] = f.Value
		}
	}

	if len(fields) == 0 {
		return 0, errs.NewError(errors.New("no fields to update"), errs.RedisErrorType, errs.ErrDataIllegal)
	}

	res, err := client.HSet(ctx, key, fields).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}
