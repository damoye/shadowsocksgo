package pac

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Gen ...
func Gen(socks5Addr string) (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if b, err = base64.StdEncoding.DecodeString(string(b)); err != nil {
		return "", err
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
		return "", err
	}
	if b, err = json.MarshalIndent(lines, "", "    "); err != nil {
		return "", err
	}
	result := strings.Replace(pacJS, "__RULES__", string(b), 1)
	result = strings.Replace(result, "__SOCKS5ADDR__", socks5Addr, 2)
	return result, nil
}
