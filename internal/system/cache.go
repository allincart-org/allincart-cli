package system

import (
	"os"
	"path"
)

func GetAllincartCliCacheDir() string {
	cacheDir, _ := os.UserCacheDir()

	return path.Join(cacheDir, "allincart-cli")
}
