// Code generated by 'yaegi extract github.com/prometheus/client_golang/prometheus/promauto'. DO NOT EDIT.

package custom

import (
	"github.com/prometheus/client_golang/prometheus/promauto"
	"reflect"
)

func init() {
	Symbols["github.com/prometheus/client_golang/prometheus/promauto/promauto"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewCounter":      reflect.ValueOf(promauto.NewCounter),
		"NewCounterFunc":  reflect.ValueOf(promauto.NewCounterFunc),
		"NewCounterVec":   reflect.ValueOf(promauto.NewCounterVec),
		"NewGauge":        reflect.ValueOf(promauto.NewGauge),
		"NewGaugeFunc":    reflect.ValueOf(promauto.NewGaugeFunc),
		"NewGaugeVec":     reflect.ValueOf(promauto.NewGaugeVec),
		"NewHistogram":    reflect.ValueOf(promauto.NewHistogram),
		"NewHistogramVec": reflect.ValueOf(promauto.NewHistogramVec),
		"NewSummary":      reflect.ValueOf(promauto.NewSummary),
		"NewSummaryVec":   reflect.ValueOf(promauto.NewSummaryVec),
		"NewUntypedFunc":  reflect.ValueOf(promauto.NewUntypedFunc),
		"With":            reflect.ValueOf(promauto.With),

		// type definitions
		"Factory": reflect.ValueOf((*promauto.Factory)(nil)),
	}
}