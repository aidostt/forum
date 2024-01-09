package main

import (
	"flag"
	"log"
)

type application struct {
	cfg config
}

type config struct {
	port int
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "server port")
	flag.Parse()
	app := &application{
		cfg: cfg,
	}
	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
