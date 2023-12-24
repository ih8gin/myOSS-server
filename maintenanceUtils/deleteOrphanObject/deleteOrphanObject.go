package main

import (
	"MyOSS/config"
	"MyOSS/es"
	"MyOSS/utils"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(config.STORAGE_ROOT + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		hashInMetadata, e := es.HasHash(hash)
		if e != nil {
			log.Println(e)
			return
		}
		if !hashInMetadata {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete", hash)
	ip, _ := utils.GetLocalIP()
	url := "http://" + ip + config.DATANODE_LISTEN_ADDRESS + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
