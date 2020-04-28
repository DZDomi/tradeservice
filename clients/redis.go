package clients

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SetObject(prefix string, key string, object interface{}, duration time.Duration) error {
	keyToSet := prefix + ":" + key
	result, err := json.Marshal(object)
	if err != nil {
		log.Println("Error while marshaling: " + keyToSet + ", " + err.Error())
		return err
	}
	client.Set(keyToSet, result, duration)
	return nil
}

func GetObject(prefix string, key string, object interface{}) error {
	keyToGet := prefix + ":" + key
	val, err := client.Get(keyToGet).Result()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		log.Println("Error while trying to get redis key: " + keyToGet + ", " + err.Error())
		return err
	}
	if err = json.Unmarshal([]byte(val), object); err != nil {
		log.Println("Error while unmarshalling: " + err.Error())
		return err
	}
	return nil
}
