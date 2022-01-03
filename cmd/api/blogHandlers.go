package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-play/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"strconv"
	"time"

	snowFlake "github.com/beinan/fastid"
	"go-play/common/httpResponse"
	"go-play/common/mongoHelper"
)

type Blog struct {
	BlogId int64 `bson:"blog_id"`
	Views int `bson:"views"`
	Cover string `bson:"cover"`
	Title string `bson:"title"`
	Content string `bson:"content"`
	CreateTime primitive.DateTime `bson:"create_time"`
	UpdateTime primitive.DateTime `bson:"update_time"`
}

const errHeader = "api/handlers"

/* (http.MethodGet, "/test", app.test)*/
func (app *application) test(w http.ResponseWriter, r *http.Request) {
	log.Println("GET.test")

	currentStatus := AppStatus{
		Status:      "Available",
		Environment: "app.config.env",
		Version:     "1.0.0",
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

/* (http.MethodGet, "/testMongo", app.testMongo) */
func (app *application) testMongo(w http.ResponseWriter, r *http.Request) {
	log.Println("GET.testMongo")

	client, ctx := MongoConnection()
	defer client.Disconnect(ctx)
	collection := client.Database("my_blogs").Collection("blogs")

	var blogs Blog
	err := collection.FindOne(ctx, bson.D{{"view", 0}}).Decode(&blogs)
	if err != nil { log.Fatal("testMongo.collection Find: \n", err) }

	fmt.Println(blogs)

	httpResponse.ReturnSuccessStatus(w, r, blogs)
}

/* (http.MethodGet, "/testUpdate/:id", app.testUpdateBlog) */
func (app *application) testUpdateBlog(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.testUpdateBlog")

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		log.Fatalln(errHeader + ".testUpdateBlog strconv.Atoi error: \n", err)
		return
	}

	client, collection, ctx, err := getDBCollectionBlogs(w, r)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	//
	filter := bson.D{{"blog_id", id}}
	update := bson.D{{"$set",
		bson.D{{"blog_id", snowFlake.CommonConfig.GenInt64ID()}}}}

	cur, err := collection.CountDocuments(ctx, filter)
	if cur < 1 {
		httpResponse.ReturnSuccessStatus(w, r, "No such document")
		return
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatalln("handlers.testUpdateBlog error in UpdateOne: \n", err)
	}
	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
	//

	httpResponse.ReturnSuccessStatus(w, r, result)
}

/* (http.MethodGet, "/testInsert", app.testInsertBlog) */
func (app *application) testInsertBlog(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.testInsertBlog")

	client, collection, ctx, err := getDBCollectionBlogs(w, r)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	blog := Blog {
		BlogId: snowFlake.CommonConfig.GenInt64ID(),
		Views: 0,
		Cover: "6343",
		Title: "tititi",
		Content: "content",
		CreateTime: primitive.NewDateTimeFromTime(time.Now()),
		UpdateTime: primitive.NewDateTimeFromTime(time.Now()),

	}
	//
	doc, err := mongoHelper.ToDoc(blog)
	if err != nil {
		log.Fatalln(errHeader + ".testInsertBlog mongoHelper.ToDoc error: \n", err)
	}

	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatalln(errHeader + ".testInsertBlog collection.InsertOne error: \n", err)
	}

	httpResponse.ReturnSuccessStatus(w, r, result)
	//
}

func (app *application) testDeleteBlog(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.testInsertBlog")

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		log.Fatalln(errHeader + ".testDeleteBlog strconv.Atoi error: \n", err)
		return
	}

	client, collection, ctx, err := getDBCollectionBlogs(w, r)
	if err != nil {
		log.Fatal("handlers.testDeleteBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	filter := bson.D{{"blog_id", id}}
	//opts := options.Delete().SetHint(bson.D{{"_id", 1}})

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	httpResponse.ReturnSuccessStatus(w, r, result.DeletedCount)

	log.Printf("Number of documents deleted: %d\n", result.DeletedCount)
}

/* ######################## real functions ######################## */
func (app *application) postBlog(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.postBlog")

	// decode blog json from r *http.Request
	decoder := json.NewDecoder(r.Body)
	var blog Blog
	err := decoder.Decode(&blog)
	if err != nil {
		panic(err)
	}

	blog.BlogId = snowFlake.CommonConfig.GenInt64ID()
	blog.CreateTime = primitive.NewDateTimeFromTime(time.Now())
	blog.UpdateTime = primitive.NewDateTimeFromTime(time.Now())

	client, collection, ctx, err := getDBCollectionBlogs(w, r)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)
	//
	doc, err := mongoHelper.ToDoc(blog)
	if err != nil {
		log.Fatalln(errHeader + ".testInsertBlog mongoHelper.ToDoc error: \n", err)
	}

	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatalln(errHeader + ".testInsertBlog collection.InsertOne error: \n", err)
	}

	httpResponse.ReturnSuccessStatus(w, r, result)
}

func (app *application) getBlogs(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.getBlogs")

	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		panic(err)
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		panic(err)
	}

	client, collection, ctx, err := getDBCollectionBlogs(w, r)
	if err != nil {
		log.Fatal("handlers.testUpdateBlog() error from getDBCollectionBlogs: ", err)
	}

	defer client.Disconnect(ctx)

	var startPage int64
	if start > 0 {
		startPage = int64((start - 1) * size)
	} else {
		start = 0
	}

	cur, err := collection.Find(context.TODO(), bson.D{}, options.Find().SetSkip(startPage).SetLimit(int64(size)))

	var result []Blog

	if err = cur.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	fmt.Println(result)
	httpResponse.ReturnSuccessStatus(w, r, result)
}
/* ######################## END REGION ######################## */

/* (http.MethodGet, "/testDelete/:id", app.testDeleteBlog) */
func getDBCollectionBlogs(w http.ResponseWriter, r *http.Request) (*mongo.Client, *mongo.Collection, context.Context, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.GetMongoAPI()))
	if err != nil {
		httpResponse.ReturnInternalError(w, r, err, errHeader + "testUpdateBlog mongo.NewClient error: \n")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		httpResponse.ReturnInternalError(w, r, err, errHeader + "testUpdateBlog client.Connect error: \n")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		httpResponse.ReturnInternalError(w, r, err, errHeader + "testUpdateBlog client.Ping error: \n")
	}

	collection := client.Database("my_blogs").Collection("blogs")

	return client, collection, ctx, nil
}
