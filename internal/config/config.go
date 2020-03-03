package config

import (
	"github.com/ghodss/yaml"
	"github.com/major1201/goutils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if err != nil {
		return errors.Wrapf(err, "open config file error: %s", Path)
	}
	defer yamlFile.Close()

	// read all
	yamlByte, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return errors.Wrap(err, "read config file error")
	}

	config := &MutatorConfig{}

	if err := yaml.Unmarshal(yamlByte, config); err != nil {
		return errors.Wrap(err, "unmarshal config file error")
	}

	setDefaultValues(config)
	if err := setRuleSelectors(config.Rules); err != nil {
		return err
	}

	CurrentConfig = config

	zap.L().Named("config").Info("config file loaded", zap.String("path", Path), zap.Any("currentConfig", CurrentConfig))

	return nil
}

func setDefaultValues(config *MutatorConfig) {
	if goutils.IsBlank(config.AnnotationKey) {
		config.AnnotationKey = "kubemutator.example.com/requests"
	}
}

func setRuleSelectors(rules []*Rule) error {
	for _, rule := range rules {
		selector, err := metav1.LabelSelectorAsSelector(rule.Selector)
		if err != nil {
			return errors.Wrap(err, "LabelSelectorAsSelector error")
		}
		rule.selector = selector
	}
	return nil
}
