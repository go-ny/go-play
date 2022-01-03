package blog

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	BlogId int64 `bson:"blog_id"`
	Views int `bson:"views"`
	Cover string `bson:"cover"`
	Title string `bson:"title"`
	Content string `bson:"content"`
	CreateTime primitive.DateTime `bson:"create_time"`
	UpdateTime primitive.DateTime `bson:"update_time"`
}
