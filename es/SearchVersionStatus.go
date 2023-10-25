package es

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`{
		"size": 0,
		"aggs": {
			"group_by_name": {
				"terms": {
					"field": "name",
					"min_doc_count": %d
				},
				"aggs": {
					"min_version": {
						"min": {
							"field": "version"
						}
					}
				}
			}
		}
	}`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	// new es(8.+) api require to declare content-type in headers
	request.Header.Set("content-type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("receive status code {%d} from es_server!", r.StatusCode))
	}
	b, _ := io.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}
