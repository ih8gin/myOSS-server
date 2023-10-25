package main

import (
	"MyOSS/dataServer/heartbeat"
	"MyOSS/dataServer/locate"
	"MyOSS/dataServer/objects"
	"MyOSS/dataServer/temp"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("dataServer starting...")
	locate.CollectObjects()
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	go temp.TempCleanWorcker()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
