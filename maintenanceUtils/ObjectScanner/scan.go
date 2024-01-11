package main

import (
	"MyOSS/apiServer/objects"
	"MyOSS/config"
	"MyOSS/es"
	"MyOSS/utils"
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(config.STORAGE_ROOT + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
	}
}

func verify(hash string) {
	utils.Logger.Warn(fmt.Sprintf("verify {%s}", hash))
	size, e := es.SearchHashSize(hash)
	if e != nil {
		utils.Logger.Warn(e.Error())
		return
	}
	stream, e := objects.GetStream(hash, size)
	if e != nil {
		utils.Logger.Warn(e.Error())
		return
	}
	d := utils.CalculateHash(stream)
	if d != hash {
		utils.Logger.Warn(fmt.Sprintf("object hash mimatch, culculated=%s, requested=%s", d, hash))
		return
	}
	stream.Close()
}
