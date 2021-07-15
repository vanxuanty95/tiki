package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"tiki/internal/pkg/logger"
	"time"
)

type (
	Config struct {
		State      string
		RestfulAPI struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		} `yaml:"restful_api"`
		DB struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"db"`
		JWTKey   string        `yaml:"jwt_key"`
		SSExpire time.Duration `yaml:"ss_expire"`
		Distance int           `yaml:"distance"`
	}
)

func Load(state *string) (*Config, error) {
	cfgPath := fmt.Sprintf("%v/config/config.%v.yml", RootDir(), *state)
	f, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Errorf("Fail to open configurations file: %v", err)
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		logger.Errorf("Fail to decode configurations file: %v", err)
		return nil, err
	}
	cfg.State = *state
	return &cfg, nil
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
