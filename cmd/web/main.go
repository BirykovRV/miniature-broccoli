package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BirykovRV/miniature-broccoli/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

type config struct {
	addr        string
	staticDir   string
	environment string
	dsn         string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address or port [:4000]")
	flag.StringVar(&cfg.staticDir, "static", "./ui/static/", "Static files web directory")
	flag.StringVar(&cfg.environment, "env", "dev", "Set environment of App")
	flag.StringVar(&cfg.dsn, "dsn", "web:Trd19afo@(127.0.0.1:3307)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(cfg),
	}

	app.infoLog.Printf("Starting server on %s with environment setting: %s", cfg.addr, cfg.environment)

	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
