package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/alextonkonogov/ascii_counter/pkg/ftpConnection"
	"github.com/alextonkonogov/ascii_counter/pkg/service"
)

type config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

func main() {
	cfg := config{}
	path := filepath.Join("./ftp.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &cfg)

	asciicounter := func(rw http.ResponseWriter, r *http.Request) {
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
