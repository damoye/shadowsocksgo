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
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if b, err = base64.StdEncoding.DecodeString(string(b)); err != nil {
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

type server struct {
	content string
}

func (s *server) handlePAC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("HTTP method not allowed: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write([]byte(s.content)); err != nil {
		log.Print("write: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Print("GET /proxy.pac")
	}
}

// Start starts to serve PAC
func Start() {
	s := server{content: generate()}
	http.HandleFunc("/proxy.pac", s.handlePAC)
	go func() {
		panic(http.ListenAndServe(consts.HTTPAddr, nil))
	}()
}
