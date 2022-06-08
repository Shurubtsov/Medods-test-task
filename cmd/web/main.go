package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dshurubtsov/cmd/config"
)

func main() {
	// configuration flags before run app
	addr := flag.String("addr", ":4000", "Net address HTTP")

	// configure loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
