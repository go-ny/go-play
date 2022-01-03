package blogRequests

import (
	"context"
	"fmt"
	snowFlake "github.com/beinan/fastid"
	"github.com/gin-gonic/gin"
	"go-play/blog"
	"go-play/common/mongoHelper"
	"go-play/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func TestMongo(c *gin.Context) {
	log.Println("GET.testMongo")

	client, ctx := mongoHelper.MongoConnection(c)
	defer client.Disconnect(ctx)
	collection := client.Database("my_blogs").Collection("blogs")

	var blogs blog.Blog
	err := collection.FindOne(ctx, bson.D{{"views", 0}}).Decode(&blogs)
	if err != nil {
		//log.Fatal("testMongo.collection Find: \n", err)
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": "Mongo FindOne error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"blog": blogs,
	})
}

/* ######################## real functions ######################## */

// PostBlog after uploaded the cover image,
// should get the cover image URL from web page /*
func PostBlog(c *gin.Context) {
	log.Println("handlers.postBlog")

	// decode post json from r *http.Request
	var post blog.Blog

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "PostBlog",
			"error": "bindJson error",
		})
	}

	post.BlogId = snowFlake.CommonConfig.GenInt64ID()
	post.CreateTime = primitive.NewDateTimeFromTime(time.Now())
	post.UpdateTime = primitive.NewDateTimeFromTime(time.Now())

	client, collection, ctx, err := getDBCollectionBlogs(c)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)
	//
	doc, err := mongoHelper.ToDoc(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("bson.Unmarshal error: %v", err),
		})
	}

	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("postBlog InsertOne error: %v", err),
		})
	}

	c.JSON(http.StatusOK, gin.H {
		"result": result,
	})
}

// GetBlog if no document found in MongoDB, should return Status 500 to web page /*
func GetBlog(c *gin.Context) {
	log.Println("handlers.postBlog")

	rawBlogId := c.Query("blog-id")

	blogId, err := strconv.Atoi(rawBlogId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("Atoi error: %v", err),
		})
		return
	}

	client, ctx := mongoHelper.MongoConnection(c)
	defer client.Disconnect(ctx)
	collection := client.Database("my_blogs").Collection("post")

	ipAddr := c.GetHeader("X-Forwarded-For")

	fmt.Println("idAddr: \n", ipAddr)

	var blogs blog.Blog

	filter := bson.D{{ "blog_id", bsonx.Int64(int64(blogId)) }}

	err = collection.FindOne(ctx, filter).Decode(&blogs)
	if err != nil {
		fmt.Println("err.Error(): ", err.Error())
		if strings.Contains(err.Error(), "no document") {
			blogs.BlogId = -1
			c.JSON(http.StatusBadRequest, gin.H {
				"error": "no document",
			})
			return
		} else {
			//log.Fatal("testMongo.collection Find: \n", err)
			c.JSON(http.StatusInternalServerError, gin.H {
				"error": fmt.Sprintf("Mongo FindOne error: %v\n", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H {
		"blog": blogs,
	})
}

// GetBlogList get pagination blogs /*
func GetBlogList(c *gin.Context) {
	log.Println("handlers.getBlogs")

	rawStart, _ := c.GetQuery("start")
	rawSize, _ := c.GetQuery("size")

	start, err1 := strconv.Atoi(rawStart)
	size, err2 := strconv.Atoi(rawSize)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": "Atoi error",
		})
	}

	client, collection, ctx, err := getDBCollectionBlogs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("MongoDB error: %v", err),
		})
		return
	}

	defer client.Disconnect(ctx)

	var startPage int64
	if start > 0 {
		startPage = int64((start - 1) * size)
	} else {
		start = 0
	}

	cur, err := collection.Find(context.TODO(), bson.D{}, options.Find().SetSkip(startPage).SetLimit(int64(size)))

	var result []blog.Blog

	if err = cur.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("mongo All iterator error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"blogs": result,
	})
}

func UpdateBlog(c *gin.Context) {
	log.Println("handlers.updateBlog")

	var post blog.Blog

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "PostBlog",
			"error": "bindJson error",
		})
	}

	filter := bson.D{{"blog_id", post.BlogId}}
	update := bson.D{{"$set",
		bson.D{
		{ "update_time", primitive.NewDateTimeFromTime(time.Now()) },
		{ "cover", post.Cover },
		{ "title", post.Title },
		{ "content", post.Content },
		},
	}}


	client, collection, ctx, err := getDBCollectionBlogs(c)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("Mongo updateOne error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"msg": result,
	})
}

func DeleteBlog(c *gin.Context) {
	log.Println("handlers.deleteBlog")

	rawBlogId := c.Query("blog-id")

	blogId, err := strconv.Atoi(rawBlogId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("Atoi error: %v", err),
		})
		return
	}

	client, collection, ctx, err := getDBCollectionBlogs(c)
	if err != nil {
		log.Fatal("handlers.testDeleteBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	filter := bson.D{{"blog_id", bsonx.Int64((int64)(blogId))}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H {
		"msg": result,
	})
}

func getDBCollectionBlogs(c *gin.Context) (*mongo.Client, *mongo.Collection, context.Context, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.GetMongoAPI()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "getDBCollectionBlogs",
			"error": "Mongo NewClient error",
			"msg": "please check in mongodb console",
		})

		return nil, nil, nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "getDBCollectionBlogs",
			"error": "Mongo client.Connect error",
			"msg": "please check client connection in mongo console",
		})

		return client, nil, ctx, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"position": "getDBCollectionBlogs",
			"msg": "client ping failed",
		})

		return client, nil, ctx, err
	}

	collection := client.Database("my_blogs").Collection("blogs")

	return client, collection, ctx, nil
}