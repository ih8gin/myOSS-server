package locate

import (
	"MyOSS/config"
	"MyOSS/rabbitmq"
	"MyOSS/types"
	"MyOSS/utils"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	q := rabbitmq.New(config.RABBITMQ_SERVER)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		id := Locate(hash)
		if id != -1 {
			ip, _ := utils.GetLocalIP()
			q.Send(msg.ReplyTo, types.LocateMessage{ip + config.DATANODE_LISTEN_ADDRESS, id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(config.STORAGE_ROOT + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
	}
}
