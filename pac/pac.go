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

var pacContent = gen()

func gen() string {
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

func getPac(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("not allowed HTTP method: ", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write([]byte(pacContent)); err != nil {
		log.Print("write: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print("GET /proxy.pac")
}

// Start ...
func Start() {
	http.HandleFunc("/proxy.pac", getPac)
	go func() {
		panic(http.ListenAndServe(consts.HTTPAddr, nil))
	}()
}
