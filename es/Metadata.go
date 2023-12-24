package es

import (
	"MyOSS/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Metadata struct {
	Name    string
	Version int   `json:"version,string"`
	Size    int64 `json:"size,string"`
	Hash    string
}

func getMetaData(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_source/%s_%d", config.ES_SERVER, name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := io.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetaData(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":"%d","size":"%d","hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", config.ES_SERVER, name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	// new es(8.+) api require to declare content-type in headers
	request.Header.Set("content-type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := io.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", config.ES_SERVER, name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}
