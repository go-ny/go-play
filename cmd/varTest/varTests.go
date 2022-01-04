package main

import (
	"context"
	"github.com/joho/godotenv"
	"go-play/redisOp"
	"go-play/routers"
	"log"
	"os"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	LoadEnv()

	//awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	//fmt.Println("My access key ID is ", awsAccessKeyID)

	router := routers.Routers()

	rdb := redisOp.RedisNewClient()
	cn := rdb.Conn(context.Background())

	defer cn.Close()

	_ = router.Run(":4000")
}
