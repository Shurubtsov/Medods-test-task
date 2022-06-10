package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dshurubtsov/cmd/config"
	"github.com/dshurubtsov/pkg/mongodb"
	"github.com/dshurubtsov/pkg/tokens"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// configuration flags before run app
	addr := flag.String("addr", ":4000", "Net address HTTP")
	dsn := flag.String("dsn", "mongodb://localhost:27017/", "Options for connect to MongoDB")
	flag.Parse()

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

	// init jwt maker for tokens
	tokenManager, err := tokens.NewManager("secret-phrase0907-331-2356")
	if err != nil {
		errorLog.Fatal("can't create token manager")
	}

	// init application obj
	app := &config.Application{
		ErrorLog:     errorLog,
		InfoLog:      infoLog,
		UserModel:    &mongodb.UserModel{DB: mongoClient},
		TokenManager: tokenManager,
	}

	// initialize custom server for our logs
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  Routes(app),
	}

	// start server
	infoLog.Printf("start server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Open and connect to mongodb; return client for querys to db
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
