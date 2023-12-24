package heartbeat

import (
	"MyOSS/config"
	"MyOSS/rabbitmq"
	"time"
)

func StartHeartBeat() {
	q := rabbitmq.New(config.RABBITMQ_SERVER)
	defer q.Close()
	for {
		q.Publish("apiServers", config.DATANODE_LISTEN_ADDRESS)
		//log.Printf("heartbeat \n")
		time.Sleep(5 * time.Second)
	}
}
