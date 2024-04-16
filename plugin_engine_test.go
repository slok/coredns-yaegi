package corednsyaegi_test

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"

	corednsyaegi "github.com/slok/coredns-yaegi"
)

func TestLoadPlugin(t *testing.T) {
	tests := map[string]struct {
		pluginSrc  string
		execPlugin func(t *testing.T, new corednsyaegi.NewPluginAPISignature)
		expErr     bool
	}{
		"Empty plugin should fail.": {
			pluginSrc:  "",
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {},
			expErr:     true,
		},

		"An invalid plugin syntax should fail": {
			pluginSrc:  `package test{`,
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {},
			expErr:     true,
		},

		"A plugin without the required factory function, should fail.": {
			pluginSrc:  `package test`,
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {},
			expErr:     true,
		},

		"A correct plugin should be loaded correctly.": {
			pluginSrc: `
package test

import (
	"context"
	"fmt"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type plugin struct {
	next corednsplugin.Handler
}

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	return plugin{next: next}
}

func (p plugin) Name() string { return "test" }

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}
`,
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {},
			expErr:     false,
		},

		"A correct plugin that returns error should return an error.": {
			pluginSrc: `
package test

import (
	"context"
	"fmt"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type plugin struct {
	next corednsplugin.Handler
}

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	return plugin{next: next}
}

func (p plugin) Name() string { return "test" }

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return 0, fmt.Errorf("something")
}
`,
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {
				rec := dnstest.NewRecorder(&test.ResponseWriter{})

				handler := new(nil)
				msg := dns.Msg{}
				_, err := handler.ServeDNS(context.TODO(), rec, &msg)
				assert.Error(t, err)

			},
			expErr: false,
		},

		"A correct plugin (https://github.com/coredns/demo) that returns correctly the response should execute correctly.": {
			pluginSrc: `
// Package demo implements a plugin.
package demo

import (
	"context"
	"fmt"
	"net"
	"strings"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	return demo{}
}

// Demo is a plugin in CoreDNS
type demo struct{}

// ServeDNS implements the plugin.Handler interface.
func (p demo) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	reply := "8.8.8.8"
	if strings.HasPrefix(state.IP(), "172.") || strings.HasPrefix(state.IP(), "127.") {
		reply = "1.1.1.1"
	}
	fmt.Printf("Received query %s from %s, expected to reply %s\n", qname, state.IP(), reply)

	answers := []dns.RR{}

	if state.QType() != dns.TypeA {
		return dns.RcodeNameError, nil
	}

	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET}
	rr.A = net.ParseIP(reply).To4()

	answers = append(answers, rr)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (p demo) Name() string { return "demo" }
`,
			execPlugin: func(t *testing.T, new corednsyaegi.NewPluginAPISignature) {
				assert := assert.New(t)
				handler := new(nil)

				msg := dns.Msg{}
				msg.SetQuestion("example.org.", dns.TypeA)
				rec := dnstest.NewRecorder(&test.ResponseWriter{})
				code, err := handler.ServeDNS(context.TODO(), rec, &msg)
				assert.NoError(err)
				assert.Equal(dns.RcodeSuccess, code)
				assert.Equal(rec.Msg.Answer[0].String(), `example.org.	0	IN	A	8.8.8.8`)
			},
			expErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			pluginFactory, err := corednsyaegi.LoadPlugin(test.pluginSrc)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				test.execPlugin(t, pluginFactory)
			}
		})
	}
}
