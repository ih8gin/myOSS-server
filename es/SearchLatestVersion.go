package es

import (
	"MyOSS/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type hit struct {
	Source Metadata `json:"_source"`
}

type total struct {
	Value    int
	Relation string
}

type searchResult struct {
	Hits struct {
		Total total
		Hits  []hit
	}
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc", config.ES_SERVER, url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := io.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}
