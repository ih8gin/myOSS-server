package temp

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func TempCleanWorcker() {
	for {
		files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/temp/*")
		for i := range files {
			file_info, _ := os.Stat(files[i])
			winFileAttr := file_info.Sys().(*syscall.Win32FileAttributeData)
			if time.Unix(winFileAttr.LastWriteTime.Nanoseconds()/1e9, 0).Add(10 * time.Minute).Before(time.Now()) {
				os.Remove(files[i])
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
