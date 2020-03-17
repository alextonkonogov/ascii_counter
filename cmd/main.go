package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alextonkonogov/ascii_counter/pkg/config"
	"github.com/alextonkonogov/ascii_counter/pkg/ftpConnection"
	"github.com/alextonkonogov/ascii_counter/pkg/service"
)

func main() {
	asciicounter := func(rw http.ResponseWriter, r *http.Request) {
		cfg := config.NewConfig()
		err := cfg.SetConfigFromJson("./ftp.json")
		if err != nil {
			log.Fatal(err)
		}
		srv := service.NewService(
			ftpConnection.NewFtp(cfg.Host, cfg.Port, cfg.User, cfg.Password),
			rw, r,
		)
		srv.ASCIISymbolCounter(cfg.Dir)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", asciicounter).Methods("GET")
	log.Println("Service is working!")
	log.Fatal(http.ListenAndServe(":9089", r))
}
