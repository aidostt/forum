package main

import (
	"flag"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"os"
)

type application struct {
	cfg           config
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	infoLog       *log.Logger
	errorLog      *log.Logger
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
	templateCache, err := NewTemplateCache()
	if err != nil {
		errorLog.Println(err)
		return
	}
	formDecoder := form.NewDecoder()
	app := &application{
		cfg:           cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		formDecoder:   formDecoder,
		templateCache: templateCache,
	}
	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
