package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

type config struct {
	addr        string
	staticDir   string
	environment string
}

func main() {
	var cfg config

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address or port [:4000]")
	flag.StringVar(&cfg.staticDir, "static", "./ui/static/", "Static files web directory")
	flag.StringVar(&cfg.environment, "env", "dev", "Set environment of App")
	flag.Parse()

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg),
	}

	infoLog.Printf("Starting server on %s with environment setting: %s", cfg.addr, cfg.environment)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
