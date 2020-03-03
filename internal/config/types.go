package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// MutatorConfig indicates the k8s-mutator main config file structure
type MutatorConfig struct {
	AnnotationKey string      `json:"annotationKey" yaml:"annotationKey"`
	Strategies    []*Strategy `json:"strategies" yaml:"strategies"`
	Rules         []*Rule     `json:"rules" yaml:"rules"`
}

// Strategy indicates each strategy
type Strategy struct {
	Name    string  `json:"name" yaml:"name"`
	Patches []Patch `json:"patches" yaml:"patches"`
}

// Patch indicates each patch
type Patch struct {
	Data       string `json:"data" yaml:"data"`
	IsTemplate bool   `json:"isTemplate" yaml:"isTemplate"`
	IsArray    bool   `json:"isArray" yaml:"isArray"`
}

// Rule indicates each select rule
type Rule struct {
	Namespace  []string              `json:"namespace" yaml:"namespace"`
	Selector   *metav1.LabelSelector `json:"selector" yaml:"selector"`
	Strategies []string              `json:"strategies" yaml:"strategies"`

	selector labels.Selector
}

// Matches returns if the selector matches the labels
func (r *Rule) Matches(l map[string]string) bool {
	if r.selector == nil {
		return false
	}
	return r.selector.Matches(labels.Set(l))
}
