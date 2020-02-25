package mutate

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/major1201/goutils"
	"github.com/major1201/k8s-mutator/internal/config"
	"github.com/major1201/k8s-mutator/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"text/template"
)

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

func mutatePods(r *http.Request, ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	getLogger(r).Debug("processing admission review")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		getLogger(r).Error("expect resource to be pod", zap.Any("podResource", podResource))
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		err = errors.Wrap(err, "deserializer decode error")
		getLogger(r).Error("deserializer decoding error", log.Error(err))
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	jsonPatch, auditAnnotations, err := patchPod(&ar, &pod)
	if err != nil {
		getLogger(r).Error("patch error", zap.Error(err))
	} else {
		if jsonPatch != nil {
			pt := v1beta1.PatchTypeJSONPatch
			reviewResponse.PatchType = &pt
			reviewResponse.Patch = jsonPatch
			reviewResponse.AuditAnnotations = auditAnnotations

			getLogger(r).Info("mutating pod", zap.String("name", pod.GenerateName))
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
	var patches []string
	for stg := range strategies {
		exist := false
		for _, configStrategy := range config.CurrentConfig.Strategies {
			if stg == configStrategy.Name {
				exist = true
				for _, patch := range configStrategy.Patches {
					var yamlPatch string
					if patch.IsTemplate {
						// template patch
						tmpl, err := template.New(pod.GenerateName).Parse(patch.Data)
						if err != nil {
							return nil, nil, errors.WithMessage(errors.WithStack(err), "parse template error")
						}
						var bbf bytes.Buffer
						err = tmpl.Execute(&bbf, td)
						if err != nil {
							return nil, nil, errors.WithMessage(errors.WithStack(err), "execute template error")
						}

						yamlPatch = bbf.String()
					} else {
						yamlPatch = patch.Data
					}

					if goutils.IsBlank(yamlPatch) {
						continue
					}

					jp, err := yaml.YAMLToJSON([]byte(yamlPatch))
					if err != nil {
						return nil, nil, errors.New(fmt.Sprintf("yaml to json error, strategy: %s", stg))
					}

					patches = append(patches, string(jp))
				}
				break
			}
		}

		if !exist {
			return nil, nil, errors.New(fmt.Sprintf("strategy not exist: %s", stg))
		}
	}

	if len(patches) == 0 {
		return nil, nil, nil
	}
	return []byte("[" + strings.Join(patches, ",") + "]"), map[string]string{"strategies": strings.Join(mapKeys(strategies), ",")}, nil
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
