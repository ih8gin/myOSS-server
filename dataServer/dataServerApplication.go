package main

import (
	"MyOSS/config"
	"MyOSS/dataServer/heartbeat"
	"MyOSS/dataServer/locate"
	"MyOSS/dataServer/objects"
	"MyOSS/dataServer/temp"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("dataServer starting...")
	locate.CollectObjects()
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	go temp.TempCleanWorcker()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(config.DATANODE_LISTEN_ADDRESS, nil))
}
