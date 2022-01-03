package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	//secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/test", app.test)

	router.HandlerFunc(http.MethodGet, "/testMongo", app.testMongo)

	router.HandlerFunc(http.MethodGet, "/testUpdate/:id", app.testUpdateBlog)

	router.HandlerFunc(http.MethodGet, "/testInsert", app.testInsertBlog)

	router.HandlerFunc(http.MethodGet, "/testDelete/:id", app.testDeleteBlog)

	router.HandlerFunc(http.MethodGet, "/blogs", app.getBlogs)

	router.HandlerFunc(http.MethodPost, "/post-blog", app.postBlog)

	//routers.HandlerFunc(http.MethodPost, "/upload", app.upload)

	return app.enableCORS(router)
}