package temp

import (
	"MyOSS/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, e := os.Open(config.STORAGE_ROOT + "/temp/" + uuid + ".dat")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
