package settings

import (
	"fmt"
	"os"
	"path"

	"github.com/yohamta/dagu/internal/utils"
)

var ErrConfigNotFound = fmt.Errorf("config not found")

var cache map[string]string = nil

const (
	ConfigDataDir   = "DAGU__DATA"
	ConfigLogsDir   = "DAGU__LOGS"
	ConfigAdminPort = "CONFIG__ADMIN_PORT"
)

func MustGet(name string) string {
	val, err := Get(name)
	if err != nil {
		panic(fmt.Errorf("failed to get %s : %w", name, err))
	}
	return val
}

func init() {
	load()
}

func Get(name string) (string, error) {
	if val, ok := cache[name]; ok {
		return val, nil
	}
	return "", ErrConfigNotFound
}

func load() {
	dir := utils.MustGetUserHomeDir()

	cache = map[string]string{}
	cache[ConfigDataDir] = config(
		ConfigDataDir,
		path.Join(dir, "/.dagu/data"))
	cache[ConfigLogsDir] = config(ConfigLogsDir,
		path.Join(dir, "/.dagu/logs"))
	cache[ConfigAdminPort] = config(ConfigAdminPort, "8000")
}

func InitTest(dir string) {
	os.Setenv("HOME", dir)
	load()
}

func config(env, def string) string {
	val := os.ExpandEnv(fmt.Sprintf("${%s}", env))
	if val == "" {
		return def
	}
	return val
}
