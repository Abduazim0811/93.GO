package leaderboard

import (
	"fmt"
	"log"
	"math/rand"

	repository "93.GO/internal/repository/redis"
	"github.com/go-redis/redis/v8"
)

func AddRandomUsers(rdb *redis.Client, count int) {
	for i := 1; i <= count; i++ {
		userID := fmt.Sprintf("user:%d", i)
		score := rand.Float64() * 100
		err := rdb.ZAdd(repository.Ctx, "leaderboard", &redis.Z{
			Score:  score,
			Member: userID,
		}).Err()
		if err != nil {
			log.Fatalf("did not add user: %v", err)
		}
	}
}

func DisplayLeaderboard(rdb *redis.Client, topN int) {
	users, err := rdb.ZRevRangeWithScores(repository.Ctx, "leaderboard", 0, int64(topN-1)).Result()
	if err != nil {
		log.Fatalf("did not get leaderboard: %v", err)
	}

	fmt.Println("Leaderboard:")
	for i, user := range users {
		fmt.Printf("%d. %s: %f\n", i+1, user.Member, user.Score)
	}
}

func UpdateScoreAtomic(rdb *redis.Client, userID string, increment float64) {
	err := rdb.Watch(repository.Ctx, func(tx *redis.Tx) error {
		currentScore, err := tx.ZScore(repository.Ctx, "leaderboard", userID).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		newScore := currentScore + increment

		_, err = tx.TxPipelined(repository.Ctx, func(pipe redis.Pipeliner) error {
			pipe.ZAdd(repository.Ctx, "leaderboard", &redis.Z{
				Score:  newScore,
				Member: userID,
			})
			return nil
		})

		return err
	}, "leaderboard")

	if err != nil {
		log.Fatalf("did not update user score: %v", err)
	}
}
