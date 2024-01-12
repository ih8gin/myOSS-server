package utils

import (
	"fmt"
	"io"
	"net/http"
)

var isInit bool = false
var ip string

func GetLocalIP() (ipv4 string, err error) {
	if !isInit {
		resp, err := http.Get("http://myexternalip.com/raw")
		if err != nil {
			return "", nil
		}
		defer resp.Body.Close()
		content, _ := io.ReadAll(resp.Body)
		//buf := new(bytes.Buffer)
		//buf.ReadFrom(resp.Body)
		//s := buf.String()
		Logger.Info(fmt.Sprintf("init ip: %s", content))
		ip = string(content)
		isInit = true
	}
	return ip, nil
}
