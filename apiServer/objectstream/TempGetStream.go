package objectstream

import (
	"fmt"
	"io"
)

type TempGetStream struct {
	reader io.Reader
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	if server == "" || uuid == "" {
		return nil, fmt.Errorf("invalid server %s uuid %s", server, uuid)
	}
	return newGetStream("http://" + server + "/temp/" + uuid)
}
