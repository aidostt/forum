package main

import (
	"flag"
	"log"
	"os"
)

type application struct {
	cfg      config
	infoLog  *log.Logger
	errorLog *log.Logger
}

type config struct {
	port int
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "server port")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		cfg:      cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
