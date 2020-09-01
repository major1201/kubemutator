package mutate

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/major1201/goutils"
	"github.com/major1201/kubemutator/internal/config"
	"github.com/major1201/kubemutator/pkg/log"
	"github.com/major1201/kubemutator/pkg/tmpl"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// JSONPatch object
type JSONPatch struct {
	Op    JSONPatchOp `json:"op" yaml:"op"`
	Path  string      `json:"path" yaml:"path"`
	From  string      `json:"from,omitempty" yaml:"from,omitempty"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

// JSONPatchOp json patch operator type
type JSONPatchOp string

const (
	// JSONPatchOpAdd json patch operator add
	JSONPatchOpAdd JSONPatchOp = "add"
	// JSONPatchOpRemove json patch operator remove
	JSONPatchOpRemove JSONPatchOp = "remove"
	// JSONPatchOpReplace json patch operator replace
	JSONPatchOpReplace JSONPatchOp = "replace"
	// JSONPatchOpMove json patch operator move
	JSONPatchOpMove JSONPatchOp = "move"
	// JSONPatchOpCopy json patch operator copy
	JSONPatchOpCopy JSONPatchOp = "copy"
	// JSONPatchOpTest json patch operator test
	JSONPatchOpTest JSONPatchOp = "test"
)

// Valid returns if the json patch operator is valid or not
func (op JSONPatchOp) Valid() bool {
	return goutils.Contains(op, JSONPatchOpAdd, JSONPatchOpRemove, JSONPatchOpReplace, JSONPatchOpMove, JSONPatchOpCopy, JSONPatchOpTest)
}

type templateData struct {
	Pod *corev1.Pod
}

// ContainsString indicates if the target string in the list of string
func ContainsString(obj string, v ...string) bool {
	for _, o := range v {
		if obj == o {
			return true
		}
	}
	return false
}

func mutatePods(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	logger.Debug("processing admission review")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		logger.Error("expect resource to be pod", zap.Any("podResource", podResource))
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		err = errors.Wrap(err, "deserializer decode error")
		logger.Error("deserializer decoding error", log.Error(err))
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	jsonPatch, auditAnnotations, err := patchPod(&ar, &pod)
	if err != nil {
		logger.Error("patch error", zap.Error(err))
	} else {
		if jsonPatch != nil {
			pt := v1beta1.PatchTypeJSONPatch
			reviewResponse.PatchType = &pt
			reviewResponse.Patch = jsonPatch
			reviewResponse.AuditAnnotations = auditAnnotations

			logger.Info("mutating pod", zap.String("name", pod.GenerateName))
		}
	}
	return &reviewResponse
}

func patchPod(ar *v1beta1.AdmissionReview, pod *corev1.Pod) (jsonPatch []byte, auditAnnotations map[string]string, err error) {
	strategies := make(map[string]bool)
	// check rules
	for _, rule := range config.CurrentConfig.Rules {
		if matchRule(ar, pod, rule) {
			// group rule strategies
			for _, stg := range rule.Strategies {
				if !strategies[stg] {
					strategies[stg] = true
				}
			}
		}
	}
	// check annotations
	if requestAnnotation, ok := pod.Annotations[config.CurrentConfig.AnnotationKey]; ok {
		requestStrategies := strings.Split(requestAnnotation, ",")
		for _, stg := range requestStrategies {
			if !strategies[stg] {
				strategies[stg] = true
			}
		}
	}

	// prepare template data
	td := templateData{
		Pod: pod,
	}

	// generate json patch
	var patches []JSONPatch
	for stg := range strategies {
		exist := false
		for _, configStrategy := range config.CurrentConfig.Strategies {
			if stg == configStrategy.Name {
				exist = true
				for _, patch := range configStrategy.Patches {
					var yamlPatch string
					if patch.IsTemplate {
						// template patch
						yamlPatch, err = tmpl.ExecuteTextTemplate(patch.Data, td)
						if err != nil {
							return nil, nil, errors.WithMessage(errors.WithStack(err), "execute template error")
						}
					} else {
						yamlPatch = patch.Data
					}

					if goutils.IsBlank(yamlPatch) {
						continue
					}

					if patch.IsArray {
						var jps []JSONPatch
						if err := yaml.Unmarshal([]byte(yamlPatch), &jps); err != nil {
							return nil, nil, errors.Errorf("yaml unmarshal patch array error, strategy: %s", stg)
						}
						for _, jp := range jps {
							if !jp.Op.Valid() {
								return nil, nil, errors.Errorf("unknown op: %s, strategy: %s", jp.Op, stg)
							}
						}
						patches = append(patches, jps...)
					} else {
						jp := JSONPatch{}
						if err := yaml.Unmarshal([]byte(yamlPatch), &jp); err != nil {
							return nil, nil, errors.Errorf("yaml unmarshal patch error, strategy: %s", stg)
						}
						if !jp.Op.Valid() {
							return nil, nil, errors.Errorf("unknown op: %s, strategy: %s", jp.Op, stg)
						}
						patches = append(patches, jp)
					}
				}
				break
			}
		}

		if !exist {
			return nil, nil, errors.Errorf("strategy not exist: %s", stg)
		}
	}

	if len(patches) == 0 {
		return nil, nil, nil
	}
	jpsb, err := json.Marshal(patches)
	if err != nil {
		return nil, nil, errors.Wrap(err, "json marshal error")
	}
	return jpsb, map[string]string{"strategies": strings.Join(mapKeys(strategies), ",")}, nil
}

func matchRule(ar *v1beta1.AdmissionReview, pod *corev1.Pod, rule *config.Rule) bool {
	// check namespace
	if rule.Namespace != nil {
		if !ContainsString(ar.Request.Namespace, rule.Namespace...) {
			return false
		}
	}
	return rule.Matches(pod.Labels)
}

func mapKeys(m map[string]bool) (arr []string) {
	for key := range m {
		arr = append(arr, key)
	}
	return
}
