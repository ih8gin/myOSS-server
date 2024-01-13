package main

import (
	"MyOSS/apiServer/heartbeat"
	"MyOSS/apiServer/locate"
	"MyOSS/apiServer/objects"
	"MyOSS/apiServer/temp"
	"MyOSS/apiServer/version"
	"MyOSS/config"
	"MyOSS/utils"
	"net/http"
)

func main() {
	utils.InitLogger()
	utils.Logger.Info("apiServer starting...")

	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	utils.Logger.Fatal(http.ListenAndServe(config.APINODE_LISTEN_ADDRESS, nil).Error())
}
