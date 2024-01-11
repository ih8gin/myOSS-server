package objects

import (
	"MyOSS/config"
	"MyOSS/dataServer/locate"
	"MyOSS/utils"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
}

func getFile(name string) string {
	files, _ := filepath.Glob(config.STORAGE_ROOT + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		utils.Logger.Warn(fmt.Sprintf("object hash mismatch, remove {%s}", file))
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}

func sendFile(w io.Writer, file string) {
	f, e := os.Open(file)
	if e != nil {
		utils.Logger.Warn(e.Error())
		return
	}
	defer f.Close()
	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		utils.Logger.Warn(e.Error())
		return
	}
	io.Copy(w, gzipStream)
	gzipStream.Close()
}
