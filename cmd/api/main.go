package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
	db struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	flag.Parse()

	app := &application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	srv := &http.Server {
		Addr: ":4000",
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("Starting server on port 4000")

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("server starting err: \n", err)
	}
}