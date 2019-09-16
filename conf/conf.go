package conf

import (
	"easygo/utils"
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

const (
	confDir         = ""
	defaultConfName = "easygo.toml"
	defaultLogDir   = "/easygo/logs/"
	defaultDataDir  = "/easygo/data/"
	defaultHttpPort = 25555
	defaultSyncPort = 26666
)

var (
	ErrConfNotFound = errors.New("configuration file could not be found")
)
var (
	appConfPath string
	appPath     string
	workPath    string
	Conf        = &config{}
)

type config struct {
	Log struct {
		LogDir string
	}
	Http struct {
		Timeout uint32
		Ports   struct {
			HttpPort int
			SyncPort int
		}
	}

	Data struct {
		DataDir string
	}

	Nodes []string
}

func init() {
	flag.StringVar(&appConfPath, "conf", "", "config path")
	flag.Parse()
	var err error
	if appConfPath == "" {
		var appPath string
		if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(err)
		}
		var workPath string
		if workPath, err = os.Getwd(); err != nil {
			panic(err)
		}
		appConfPath = filepath.Join(workPath, confDir, defaultConfName)
		if !utils.FileExists(appConfPath) {
			appConfPath = filepath.Join(appPath, confDir, defaultConfName)
			if !utils.FileExists(appConfPath) {
				panic(ErrConfNotFound)
			}
		}

	}
	if err = parseConfig(appConfPath); err != nil {
		panic(err)
	}

}

func parseConfig(appConfPath string) error {
	if _, err := toml.DecodeFile(appConfPath, &Conf); err != nil {
		return err
	}
	//default value
	if Conf.Log.LogDir == "" {
		Conf.Log.LogDir = defaultLogDir
	}
	if Conf.Data.DataDir == "" {
		Conf.Data.DataDir = defaultDataDir
	}
	if Conf.Http.Ports.HttpPort == 0 {
		Conf.Http.Ports.HttpPort = defaultHttpPort
	}
	if Conf.Http.Ports.SyncPort == 0 && len(Conf.Nodes) > 0 {
		Conf.Http.Ports.SyncPort = defaultSyncPort
	}
	return nil
}
