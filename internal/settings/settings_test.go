package settings

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagu/internal/utils"
)

var testHomeDir string

func TestMain(m *testing.M) {
	testHomeDir = utils.MustTempDir("settings_test")
	InitTest(testHomeDir)
	os.Exit(m.Run())
}

func TestReadSetting(t *testing.T) {
	load()

	// read default configs
	for _, test := range []struct {
		Name string
		Want string
	}{
		{
			Name: ConfigDataDir,
			Want: path.Join(testHomeDir, ".dagu/data"),
		},
		{
			Name: ConfigLogsDir,
			Want: path.Join(testHomeDir, ".dagu/logs"),
		},
	} {
		val, err := Get(test.Name)
		assert.NoError(t, err)
		assert.Equal(t, val, test.Want)
	}

	// read from env variables
	for _, test := range []struct {
		Name string
		Want string
	}{
		{
			Name: ConfigDataDir,
			Want: "/home/dagu/data",
		},
		{
			Name: ConfigLogsDir,
			Want: "/home/dagu/logs",
		},
	} {
		os.Setenv(test.Name, test.Want)
		load()

		val, err := Get(test.Name)
		assert.NoError(t, err)
		assert.Equal(t, val, test.Want)

		val = MustGet(test.Name)
		assert.Equal(t, val, test.Want)
	}

	_, err := Get("Invalid_Name")
	require.Equal(t, ErrConfigNotFound, err)
}
