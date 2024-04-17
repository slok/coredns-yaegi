package allowlist_test

import (
	"context"
	"os"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	plugintest "github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	corednsyaegi "github.com/slok/coredns-yaegi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	pluginSrc, _ := os.ReadFile("./plugin.go")

	tests := map[string]struct {
		request func() *dns.Msg
		expCode int
		expErr  bool
	}{
		"A not allowed domain should not respond correctly.": {
			request: func() *dns.Msg {
				msg := dns.Msg{}
				msg.SetQuestion("example.org.", dns.TypeA)
				return &msg
			},
			expCode: dns.RcodeNameError,
		},

		"An allowed domain should respond correctly.": {
			request: func() *dns.Msg {
				msg := dns.Msg{}
				msg.SetQuestion("slok.dev.", dns.TypeA)
				return &msg
			},
			expCode: dns.RcodeSuccess,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			// We use coredns-yaegi to run the plugin on top in the test to make it more a kind of "integration test".
			plugin, err := corednsyaegi.NewCoreDNSPlugin(corednsyaegi.CoreDNSPluginConfig{
				NextPlugin: plugintest.NextHandler(dns.RcodeSuccess, nil),
				PluginsSrc: string(pluginSrc),
			})
			require.NoError(err)

			rec := dnstest.NewRecorder(&plugintest.ResponseWriter{})
			gotCode, err := plugin.ServeDNS(context.TODO(), rec, test.request())

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expCode, gotCode)
			}
		})
	}
}
