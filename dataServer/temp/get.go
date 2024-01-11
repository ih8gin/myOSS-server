package temp

import (
	"MyOSS/config"
	"MyOSS/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, e := os.Open(config.STORAGE_ROOT + "/temp/" + uuid + ".dat")
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
