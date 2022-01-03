package routers

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"go-play/awsOp"
	"go-play/blog/blogRequests"
	//"go-play/middleware"
)

func Routers() *gin.Engine {
	sess := awsOp.ConnectAws()

	// Logging to a file.
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		cors.Default()
		c.Set("sess", sess)
		c.Next()
		gin.Logger()
		gin.Recovery()
		//middleware.CORSMiddleware()
	})

	mongo := router.Group("/mongo")
	{
		mongo.GET("/testMongo", blogRequests.TestMongo)
		mongo.POST("/post-blog", blogRequests.PostBlog)
		mongo.GET("/blog", blogRequests.GetBlog)
		mongo.GET("/blogs", blogRequests.GetBlogList)
		mongo.POST("/edit", blogRequests.UpdateBlog)
		mongo.GET("/delete", blogRequests.DeleteBlog)
	}

	router.POST("/upload", awsOp.UploadImage)

	return router
}
