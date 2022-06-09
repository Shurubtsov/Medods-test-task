package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dshurubtsov/cmd/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// configuration flags before run app
	addr := flag.String("addr", ":4000", "Net address HTTP")
	dsn := flag.String("dsn", "mongodb://localhost:27017/", "Options for connect to MongoDB")

	// configure loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// context for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// client for mongo database
	mongoClient, err := openMongoDB(*dsn, ctx)
	if err != nil {
		infoLog.Fatal(err.Error())
	}
	defer mongoClient.Disconnect(ctx)

	// init application obj
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	// initialize custom server for our logs
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  Routes(app),
	}
	infoLog.Printf("start server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openMongoDB(dsn string, ctx context.Context) (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return client, nil
}
