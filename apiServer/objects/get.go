package objects

import (
	"MyOSS/apiServer/heartbeat"
	"MyOSS/apiServer/locate"
	"MyOSS/apiServer/rs"
	"MyOSS/es"
	"MyOSS/utils"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			utils.Logger.Warn(e.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, e := es.GetMetadata(name, version)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hash := url.PathEscape(meta.Hash)
	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes=%d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	encoder := "default"
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}

	if acceptGzip {
		encoder = "gzip"
		w.Header().Set("content-encoding", encoder)
		w2 := gzip.NewWriter(w)
		_, e = io.Copy(w2, stream)
		w2.Close()
	} else {
		_, e = io.Copy(w, stream)
	}
	if e != nil {
		utils.Logger.Warn(e.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream.Close()

	utils.Logger.Info(fmt.Sprintf("Request for downloading object-{%s}-v{%d} hash-{%s} encoded with {%s} accepted.", name, version, hash, encoder))
}

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
