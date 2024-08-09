package main

import (
	"93.GO/internal/leaderboard"
	repository "93.GO/internal/repository/redis"
)

func main() {
	rdb := repository.SetupRedis()

	leaderboard.AddRandomUsers(rdb, 10)
	leaderboard.DisplayLeaderboard(rdb, 10)
	leaderboard.UpdateScoreAtomic(rdb, "user:1", 20)
}
