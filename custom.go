package corednsyaegi

import (
	"reflect"
)

//go:generate yaegi extract --name corednsyaegi  github.com/coredns/coredns/plugin github.com/coredns/coredns/request github.com/miekg/dns

// yaegiCustomSymbols variable stores the map of custom symbols per package.
var Symbols = map[string]map[string]reflect.Value{}
