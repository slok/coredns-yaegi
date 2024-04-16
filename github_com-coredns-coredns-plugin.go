// Code generated by 'yaegi extract github.com/coredns/coredns/plugin'. DO NOT EDIT.

package corednsyaegi

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["github.com/coredns/coredns/plugin/plugin"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"A":                            reflect.ValueOf(plugin.A),
		"AAAA":                         reflect.ValueOf(plugin.AAAA),
		"BackendError":                 reflect.ValueOf(plugin.BackendError),
		"CNAME":                        reflect.ValueOf(plugin.CNAME),
		"ClientWrite":                  reflect.ValueOf(plugin.ClientWrite),
		"Done":                         reflect.ValueOf(plugin.Done),
		"ErrOnce":                      reflect.ValueOf(&plugin.ErrOnce).Elem(),
		"Error":                        reflect.ValueOf(plugin.Error),
		"MX":                           reflect.ValueOf(plugin.MX),
		"NS":                           reflect.ValueOf(plugin.NS),
		"Namespace":                    reflect.ValueOf(constant.MakeFromLiteral("\"coredns\"", token.STRING, 0)),
		"NextOrFailure":                reflect.ValueOf(plugin.NextOrFailure),
		"OriginsFromArgsOrServerBlock": reflect.ValueOf(plugin.OriginsFromArgsOrServerBlock),
		"PTR":                          reflect.ValueOf(plugin.PTR),
		"Register":                     reflect.ValueOf(plugin.Register),
		"SOA":                          reflect.ValueOf(plugin.SOA),
		"SRV":                          reflect.ValueOf(plugin.SRV),
		"SlimTimeBuckets":              reflect.ValueOf(&plugin.SlimTimeBuckets).Elem(),
		"SplitHostPort":                reflect.ValueOf(plugin.SplitHostPort),
		"TXT":                          reflect.ValueOf(plugin.TXT),
		"TimeBuckets":                  reflect.ValueOf(&plugin.TimeBuckets).Elem(),

		// type definitions
		"Handler":        reflect.ValueOf((*plugin.Handler)(nil)),
		"HandlerFunc":    reflect.ValueOf((*plugin.HandlerFunc)(nil)),
		"Host":           reflect.ValueOf((*plugin.Host)(nil)),
		"Name":           reflect.ValueOf((*plugin.Name)(nil)),
		"Options":        reflect.ValueOf((*plugin.Options)(nil)),
		"Plugin":         reflect.ValueOf((*plugin.Plugin)(nil)),
		"ServiceBackend": reflect.ValueOf((*plugin.ServiceBackend)(nil)),
		"Zones":          reflect.ValueOf((*plugin.Zones)(nil)),

		// interface wrapper definitions
		"_Handler":        reflect.ValueOf((*_github_com_coredns_coredns_plugin_Handler)(nil)),
		"_ServiceBackend": reflect.ValueOf((*_github_com_coredns_coredns_plugin_ServiceBackend)(nil)),
	}
}

// _github_com_coredns_coredns_plugin_Handler is an interface wrapper for Handler type
type _github_com_coredns_coredns_plugin_Handler struct {
	IValue    interface{}
	WName     func() string
	WServeDNS func(a0 context.Context, a1 dns.ResponseWriter, a2 *dns.Msg) (int, error)
}

func (W _github_com_coredns_coredns_plugin_Handler) Name() string {
	return W.WName()
}
func (W _github_com_coredns_coredns_plugin_Handler) ServeDNS(a0 context.Context, a1 dns.ResponseWriter, a2 *dns.Msg) (int, error) {
	return W.WServeDNS(a0, a1, a2)
}

// _github_com_coredns_coredns_plugin_ServiceBackend is an interface wrapper for ServiceBackend type
type _github_com_coredns_coredns_plugin_ServiceBackend struct {
	IValue       interface{}
	WIsNameError func(err error) bool
	WLookup      func(ctx context.Context, state request.Request, name string, typ uint16) (*dns.Msg, error)
	WMinTTL      func(state request.Request) uint32
	WRecords     func(ctx context.Context, state request.Request, exact bool) ([]msg.Service, error)
	WReverse     func(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error)
	WSerial      func(state request.Request) uint32
	WServices    func(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error)
}

func (W _github_com_coredns_coredns_plugin_ServiceBackend) IsNameError(err error) bool {
	return W.WIsNameError(err)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) Lookup(ctx context.Context, state request.Request, name string, typ uint16) (*dns.Msg, error) {
	return W.WLookup(ctx, state, name, typ)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) MinTTL(state request.Request) uint32 {
	return W.WMinTTL(state)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) Records(ctx context.Context, state request.Request, exact bool) ([]msg.Service, error) {
	return W.WRecords(ctx, state, exact)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) Reverse(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error) {
	return W.WReverse(ctx, state, exact, opt)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) Serial(state request.Request) uint32 {
	return W.WSerial(state)
}
func (W _github_com_coredns_coredns_plugin_ServiceBackend) Services(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error) {
	return W.WServices(ctx, state, exact, opt)
}
