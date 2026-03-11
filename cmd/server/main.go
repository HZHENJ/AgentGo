package main

import (
	"log"
	"agentgo/pkg/conf"
	db "agentgo/internal/common/mysql"
	rds "agentgo/internal/common/redis"
	"agentgo/internal/routes"
)

func main() {
	conf.Init()
	if err := db.InitDB(); err != nil {
		log.Fatalf("init db error: %v", err)
	}
	// 可选：初始化 Redis（如未使用可移除）
	if err := rds.InitRedis(); err != nil {
		log.Printf("init redis error: %v", err)
	}

	r := routes.NewRouter()

	r.Run(":8080")
}
