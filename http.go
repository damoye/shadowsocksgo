package main

import (
	"log"
	"net/http"

	"github.com/damoye/ssgo/pac"
)

type pacServer string

func (s pacServer) getPac(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("not allowed HTTP method: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write([]byte(s)); err != nil {
		log.Print("write: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print("GET /proxy.pac")
}

func startHTTP(conf *config) {
	pacJS, err := pac.Gen(conf.LocalPort)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/proxy.pac", pacServer(pacJS).getPac)
	go func() {
		panic(http.ListenAndServe(conf.HTTPAddr, nil))
	}()
}
