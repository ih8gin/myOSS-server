package temp

import (
	"MyOSS/config"
	"MyOSS/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := config.STORAGE_ROOT + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		utils.Logger.Warn(fmt.Sprintf("actual size {%d}, exceeds {%d}", actual, tempinfo.Size))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	f, e := os.Open(config.STORAGE_ROOT + "/temp/" + uuid)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}
