package pac

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/damoye/ssgo/consts"
)

func generate() string {
	resp, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	s := ""
	if err == nil {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		s = string(b)
	} else {
		log.Print("get gfwlist from github failed: ", err)
		s = gfwlist
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(b))
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && line[0] != '!' && line[0] != '[' {
			lines = append(lines, line)
		}
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
	if b, err = json.MarshalIndent(lines, "", "    "); err != nil {
		panic(err)
	}
	return strings.Replace(consts.PACTemplate, "__RULES__", string(b), 1)
}

func logHTTP(r *http.Request, status int) {
	log.Printf("%s %s %s %d", r.Method, r.RequestURI, r.Proto, status)
}

type server string

func (s *server) handlePAC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logHTTP(r, http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write([]byte(*s)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logHTTP(r, http.StatusInternalServerError)
	} else {
		logHTTP(r, http.StatusOK)
	}
}

// Start starts to serve PAC
func Start() {
	s := server(generate())
	http.HandleFunc("/proxy.pac", s.handlePAC)
	go func() {
		panic(http.ListenAndServe(consts.HTTPAddr, nil))
	}()
}
