package data_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const TIMEOUT = 45 * time.Second

type TestRedisEntity struct {
	Datatype string 
	Value string 
}

func createTestRedisClient() (*redis.Client, error) {
	// Please for god's sake, do not use a production instance.
	env := filepath.Join("..", "..", ".env.test")

	if err := godotenv.Load(env); err != nil {
		return nil, err 
	} 

	connString := os.Getenv("REDIS_CONN_STRING")
	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, err 
	}

	client := redis.NewClient(opt)
	return client, nil
}

func TestHSet(t *testing.T) {
	entityD := TestRedisEntity{
		Datatype: "[]string",
		Value: "['h', 'i']",
	}
	client, err := createTestRedisClient()
	if err != nil {
		t.Errorf("creating test redis client failed :: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	id := uuid.NewString()
	identifier := "testredisentity:" + id 

	_, err = data.HSet(ctx, client, id, entityD)
	if err != nil {
		t.Errorf("Failed to set hash :: %v", err)
	}

	data, err := client.HGetAll(ctx, identifier).Result()
	if err != nil {
		t.Errorf("failed to fetch data with [ID : %s] :: %v", id, err)
	}

	t.Logf("%v", data)
	if _, ok := data["datatype"]; !ok {
		t.Error("Data not found :: expected [Field : datatype]")
	}

	if _, ok := data["value"]; !ok {
		t.Error("Data not found :: expected [Field : value]")
	}
}