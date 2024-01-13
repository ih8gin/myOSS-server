package heartbeat

import (
	"MyOSS/config"
	"MyOSS/rabbitmq"
	"MyOSS/utils"
	"time"
)

func StartHeartBeat() {
	q := rabbitmq.New(config.RABBITMQ_SERVER)
	defer q.Close()
	for {
		ip, _ := utils.GetLocalIP()
		q.Publish("apiServers", ip+config.DATANODE_LISTEN_ADDRESS)
		//log.Printf("heartbeat \n")
		time.Sleep(5 * time.Second)
	}
}
