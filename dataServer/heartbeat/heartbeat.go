package heartbeat

import (
	"MyOSS/rabbitmq"
	"os"
	"time"
)

func StartHeartBeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		//log.Printf("heartbeat \n")
		time.Sleep(5 * time.Second)
	}
}
