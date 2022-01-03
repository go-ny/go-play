package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-play/redisOp"
	"go-play/routers"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

	router := routers.Routers()

	rdb := redisOp.RedisNewClient()
	cn := rdb.Conn(context.Background())

	redis := router.Group("/redis")
	{
		redis.GET("/g", GetTest)
		redis.GET("/s", SetTest)
	}

	defer cn.Close()

	_ = router.Run(":4000")

}

func SetTest(c *gin.Context) {
	boolCmd := redisOp.RedisNewClient().Set(c, "123", "123", 20*time.Second)
	fmt.Println("boolCmd: ", boolCmd)
	if err := boolCmd.Err(); err != nil {
		fmt.Println("redis HSet", err)
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"msg": "success",
	})
}

func GetTest(c *gin.Context) {
	strCmd := redisOp.RedisNewClient().Get(c, "123")
	if err := strCmd.Err(); err != nil {
		errBody := fmt.Sprintf("%v", err)
		if info := strings.Split(errBody, ":")[1]; info == " nil" {
			c.JSON(http.StatusOK, gin.H {
				"msg": "expired",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H {
				"error": err,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"msg": strCmd.Val(),
	})
}


