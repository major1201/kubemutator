package mutate

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/major1201/k8s-mutator/internal/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"text/template"
)

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
	zap.L().Named("mutator").Debug("mutating pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		zap.L().Named("mutator").Error("expect resource to be pod", zap.Any("podResource", podResource))
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		zap.L().Named("mutator").Error("deserializer decoding error", zap.Error(err))
		return toAdmissionResponse(err)
	}

	// debug
	zap.L().Named("mutator").Debug("Pod object", zap.Any("pod", pod))

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	jsonPatch, err := patchPod(&ar, &pod)
	if err != nil {
		zap.L().Named("mutator").Error("patch error", zap.Error(err))
	} else {
		if jsonPatch != nil {
			pt := v1beta1.PatchTypeJSONPatch
			reviewResponse.PatchType = &pt
			reviewResponse.Patch = jsonPatch

			zap.L().Named("mutator").Info("mutating pod", zap.String("name", pod.Name), zap.ByteString("jsonPatch", jsonPatch))
		}
	}
	return &reviewResponse
}

func patchPod(ar *v1beta1.AdmissionReview, pod *corev1.Pod) (jsonPatch []byte, err error) {
	var strategies []string
	for _, rule := range config.CurrentConfig.Rules {
		if matchRule(ar, pod, rule) {
			// group rule strategies
			for _, stg := range rule.Strategies {
				if !ContainsString(stg, strategies...) {
					strategies = append(strategies, stg)
				}
			}
		}
	}

	var patches []string
	for _, stg := range strategies {
		exist := false
		for _, configStrategy := range config.CurrentConfig.Strategies {
			if stg == configStrategy.Name {
				exist = true
				for _, patch := range configStrategy.Patches {
					jp, err := yaml.YAMLToJSON([]byte(patch.Data))
					if err != nil {
						return nil, errors.New(fmt.Sprintf("yaml to json error, strategy: %s", stg))
					}

					patchString := string(jp)
					if patch.IsTemplate {
						// template patch
						tmpl, err := template.New(pod.GenerateName).Parse(patchString)
						if err != nil {
							return nil, errors.WithMessage(errors.WithStack(err), "parse template error")
						}
						var bbf bytes.Buffer
						err = tmpl.Execute(&bbf, pod)
						if err != nil {
							return nil, errors.WithMessage(errors.WithStack(err), "execute template error")
						}

						patches = append(patches, bbf.String())
					} else {
						patches = append(patches, patchString)
					}
				}
				break
			}
		}

		if !exist {
			return nil, errors.New(fmt.Sprintf("strategy not exist: %s", stg))
		}
	}

	if len(patches) == 0 {
		return nil, nil
	}
	return []byte("[" + strings.Join(patches, ",") + "]"), nil
}

func matchRule(ar *v1beta1.AdmissionReview, pod *corev1.Pod, rule *config.Rule) bool {
	// check namespace
	if rule.Namespace != nil {
		if !ContainsString(ar.Request.Namespace, rule.Namespace...) {
			return false
		}
	}

	// check labels
	podLabels := pod.Labels
	for k, v := range rule.Selector {
		val, ok := podLabels[k]
		if !ok {
			return false
		}
		if v != val {
			return false
		}
	}

	return true
}
