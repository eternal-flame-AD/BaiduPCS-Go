package pcscache

import (
	"github.com/eternal-flame-AD/BaiduPCS-Go/baidupcs"
	"time"
)

var (
	// DirCache 网盘目录缓存
	DirCache = &dirCache{
		fdl:      map[string]*baidupcs.FileDirectoryList{},
		lifeTime: 1 * time.Hour,
	}
)
