package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BirykovRV/miniature-broccoli/internal/lib"
)

type config struct {
	addr        string
	staticDir   string
	environment string
}

func main() {
	mux := http.NewServeMux()

	var cfg config

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address or port [:4000]")
	flag.StringVar(&cfg.staticDir, "static", "./ui/static/", "Static files web directory")
	flag.StringVar(&cfg.environment, "env", "dev", "Set environment of App")
	flag.Parse()

	fileServer := http.FileServer(lib.NeuteredFileSystem{
		Fs: http.Dir(cfg.staticDir),
	})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s with environment setting: %s", cfg.addr, cfg.environment)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
