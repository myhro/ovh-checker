package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/api/token"
	"github.com/myhro/ovh-checker/storage"
)

func atoi(s string, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func sleep() {
	time.Sleep(60 * time.Second)
}

func main() {
	cache, err := storage.NewCache()
	if err != nil {
		log.Fatal(err)
	}

	for {
		now := fmt.Sprintf("%v", storage.Now().Unix())
		zrange := redis.ZRangeBy{
			Min: "0",
			Max: now,
		}

		log.Print("Fetching list of expired tokens")
		list, err := cache.ZRangeByScore(token.SessionSetExpirationKey, zrange)
		if err != nil {
			log.Print(err)
			sleep()
			continue
		}
		log.Print("Done.")

		if len(list) == 0 {
			sleep()
			continue
		}

		log.Print("Removing expired tokens")
		for _, t := range list {
			log.Print("Removing: ", t)

			tokenKey := token.SessionTokenKey(t)
			userID, err := atoi(cache.HGet(tokenKey, "user_id"))
			if err != nil {
				log.Print(err)
				continue
			}
			setKey := token.SessionTokenSetKey(userID)

			tx := cache.TxPipeline()
			tx.Del(tokenKey)
			tx.SRem(setKey, t)
			tx.ZRem(token.SessionSetExpirationKey, t)
			_, err = tx.Exec()
			if err != nil {
				log.Print(err)
			}
		}

		sleep()
	}
}
