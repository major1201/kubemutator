package config

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
}

// Rule indicates each select rule
type Rule struct {
	Namespace  []string          `json:"namespace" yaml:"namespace"`
	Selector   map[string]string `json:"selector" yaml:"selector"`
	Strategies []string          `json:"strategies" yaml:"strategies"`
}
