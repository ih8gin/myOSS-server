package main

import (
	"MyOSS/config"
	"MyOSS/dataServer/heartbeat"
	"MyOSS/dataServer/locate"
	"MyOSS/dataServer/objects"
	"MyOSS/dataServer/temp"
	"MyOSS/utils"
	"net/http"
)

func main() {
	utils.InitLogger()
	utils.Logger.Info("dataServer starting...")

	locate.CollectObjects()
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	go temp.TempCleanWorcker()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	utils.Logger.Fatal(http.ListenAndServe(config.DATANODE_LISTEN_ADDRESS, nil).Error())
}
