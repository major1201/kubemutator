package tmpl

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
)

func mulQuantity(c float64, qstr interface{}) string {
	q, err := resource.ParseQuantity(fmt.Sprintf("%v", qstr))
	if err != nil {
		return ""
	}

	q.SetMilli(int64(float64(q.MilliValue()) * c))
	return q.String()
}

func divQuantity(c float64, qstr interface{}) string {
	q, err := resource.ParseQuantity(fmt.Sprintf("%v", qstr))
	if err != nil {
		return ""
	}

	q.SetMilli(int64(float64(q.MilliValue()) / c))
	return q.String()
}
