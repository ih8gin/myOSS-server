package main

import (
	"MyOSS/apiServer/heartbeat"
	"MyOSS/apiServer/locate"
	"MyOSS/apiServer/objects"
	"MyOSS/apiServer/temp"
	"MyOSS/apiServer/version"
	"MyOSS/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("apiServer starting...")
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(config.APINODE_LISTEN_ADDRESS, nil))
}
