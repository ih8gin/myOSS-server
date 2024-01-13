package objects

import (
	"MyOSS/apiServer/heartbeat"
	"MyOSS/apiServer/locate"
	"MyOSS/apiServer/rs"
	"MyOSS/es"
	"MyOSS/utils"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]

	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		utils.Logger.Warn("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.Logger.Info(fmt.Sprintf("Received resumable uploading request for object with hash {%s}", hash))

	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if locate.Exist(url.PathEscape(hash)) {
		e = es.AddVersion(name, hash, size)
		if e != nil {
			utils.Logger.Warn(e.Error())
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		utils.Logger.Warn("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
	utils.Logger.Info(fmt.Sprintf("Prepared to receive resumable transmission successfully, hash {%s}, target dataNode {%s}", hash, ds))
}
