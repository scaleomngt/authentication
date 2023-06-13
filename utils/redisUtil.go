package utils

import (
	"github.com/go-redis/redis"
	"id-card-server/config"
	"log"
	"sync"
)

var redisClient *redis.Client
var mtx sync.Mutex

var RedisServer = config.Config.GetString("RedisServer")
var Address = config.Config.GetString("Address")

func init() {
	doInit()
}

func doInit() {
	if redisClient != nil {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if redisClient != nil {
		return
	}

	if RedisServer != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     RedisServer,
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		})
	} else {
		log.Panic("RedisServer is not null")
	}
}

func GetId() (string, error) {
	doInit()
	id, err := redisClient.Get(Address).Result()
	// 判断查询是否出错
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("id的值：", id)
	return id, nil
}

func SetId(val string) error {
	doInit()
	err := redisClient.Set(Address, val, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SetHash(key, key1, value string) error {
	doInit()
	err := redisClient.HSet(key, key1, value).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetOneHash(key, key1 string) (string, error) {
	doInit()
	result, err := redisClient.HGet(key, key1).Result()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return result, nil
}

func GetAllHash(key string) ([]string, error) {
	doInit()
	result := make([]string, 0)

	data, err := redisClient.HGetAll(key).Result()
	if err != nil {
		log.Println(err)
		return result, err
	}

	for field, val := range data {
		log.Println(field, val)
		result = append(result, val)
	}
	return result, nil
}
