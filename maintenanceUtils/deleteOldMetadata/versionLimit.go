package main

import (
	"MyOSS/es"
	"MyOSS/utils"
)

const MIN_VERSION_COUNT = 5

func main() {
	buckets, e := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		utils.Logger.Warn(e.Error())
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
