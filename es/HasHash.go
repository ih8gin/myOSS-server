package es

import (
	"MyOSS/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", config.ES_SERVER, hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := io.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total.Value != 0, nil
}
