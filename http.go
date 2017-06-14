package main

import (
	"net/http"
	"os"

	"github.com/damoye/ssgo/pac"
	"github.com/gorilla/handlers"
)

type pacServer string

func (s pacServer) getPac(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write([]byte(s)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func startHTTP(conf *config) {
	pacJS, err := pac.Gen(conf.LocalAddr)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/proxy.pac", pacServer(pacJS).getPac)
	go func() {
		panic(http.ListenAndServe(
			conf.HTTPAddr,
			handlers.LoggingHandler(os.Stdout, http.DefaultServeMux),
		))
	}()
}
