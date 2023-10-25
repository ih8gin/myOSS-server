package objects

import (
	"MyOSS/dataServer/locate"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	//log.Println(os.Getenv("STORAGE_ROOT") + "/objects/" + filepath.Base(files[0]))
	//log.Println(os.Getenv("STORAGE_ROOT") + "/garbage/" + filepath.Base(files[0]))
	err := os.Rename(os.Getenv("STORAGE_ROOT")+"/objects/"+filepath.Base(files[0]), os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
	if err != nil {
		log.Println(err)
	}
	locate.Del(hash)
}
