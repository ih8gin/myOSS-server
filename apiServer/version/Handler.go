package version

import (
	"MyOSS/es"
	"MyOSS/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 100
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, e := es.SearchAllVersion(name, from, size)
		if e != nil {
			utils.Logger.Warn(e.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
