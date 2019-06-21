package config

import (
	"github.com/ghodss/yaml"
	"github.com/major1201/goutils"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

// Path indicates the config file path
var Path string

// CurrentConfig indicates the current running config
var CurrentConfig *MutatorConfig

// SetPath sets the config file path
func SetPath(path string) {
	Path = path
}

// LoadConfig loads the config from the file
func LoadConfig() error {
	yamlFile, err := os.Open(Path)
	defer yamlFile.Close()
	if err != nil {
		zap.L().Named("config").Fatal("open config file error", zap.String("path", Path), zap.Error(err))
	}

	// read all
	yamlByte, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		zap.L().Named("config").Fatal("read config file error", zap.Error(err))
	}

	config := &MutatorConfig{}

	if err = yaml.Unmarshal(yamlByte, config); err != nil {
		return err
	}

	setDefaultValues(config)

	CurrentConfig = config

	zap.L().Named("config").Info("config file loaded", zap.String("path", Path), zap.Any("currentConfig", CurrentConfig))

	return nil
}

func setDefaultValues(config *MutatorConfig) {
	if goutils.IsBlank(config.AnnotationKey) {
		config.AnnotationKey = "k8s-mutator.example.com/requests"
	}
}
