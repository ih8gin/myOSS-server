package temp

import (
	"MyOSS/config"
	"MyOSS/dataServer/locate"
	"MyOSS/utils"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempInfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := config.STORAGE_ROOT + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempInfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempInfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f.Close()
	commitTempObject(datFile, tempInfo)
}

func commitTempObject(datFile string, tempInfo *tempInfo) {
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	var e error
	_, e = f.Seek(0, io.SeekStart)
	if e != nil {
		log.Println(e)
	}
	w, _ := os.Create(config.STORAGE_ROOT + "/objects/" + tempInfo.Name + "." + d)
	w2 := gzip.NewWriter(w)
	_, e = io.Copy(w2, f)
	if e != nil {
		log.Println(e)
	}
	f.Close()
	w2.Close()
	e = os.Remove(datFile)
	if e != nil {
		log.Println(e)
	}
	//os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Name+"."+d)
	//err := os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Name+"."+d)
	//if err != nil {
	//	log.Println(err)
	//}
	locate.Add(tempInfo.hash(), tempInfo.id())
}
