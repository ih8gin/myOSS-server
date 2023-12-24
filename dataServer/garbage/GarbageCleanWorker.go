package garbage

import (
	"MyOSS/config"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// TODO 定期删除/garbage/路径下的长期文件，恢复存在引用的文件
func GarbageCleanWorcker() {
	for {
		files, _ := filepath.Glob(config.STORAGE_ROOT + "/temp/*")
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
