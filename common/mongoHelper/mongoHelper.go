package mongoHelper

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-play/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func MongoConnection(c *gin.Context) (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.GetMongoAPI()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "MongoConnection",
			"error": "Mongo NewClient error",
		})
		return nil, nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "MongoConnection",
			"error": "context.WithTimeout error",
		})
		return client, nil
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "MongoConnection",
			"error": "Mongo client.Ping error",
		})
		return client, ctx
	}

	return client, ctx
}