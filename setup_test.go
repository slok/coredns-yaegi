package corednsyaegi_test

import (
	"testing"

	"github.com/coredns/caddy"
	corednsyaegi "github.com/slok/coredns-yaegi"
	"github.com/stretchr/testify/assert"
)

var noopSrc = `package noop

import (
	"context"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type plugin struct {
	next corednsplugin.Handler
}

func NewPlugin(next corednsplugin.Handler, rawOpts string) corednsplugin.Handler {
	return plugin{next: next}
}

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "noop" }
`

func TestConfigParse(t *testing.T) {
	tests := map[string]struct {
		config         string
		expErr         bool
		expSrcWithOpts []corednsyaegi.PluginSrcWithOpts
	}{
		"No config should error.": {
			config: ``,
			expErr: true,
		},

		"A configuration with the plugin on the same line with options should load the plugin correctly.": {
			config: `yeagi { ./_examples/noop/noop.go "k1=v2,k3=v4" }`,
			expSrcWithOpts: []corednsyaegi.PluginSrcWithOpts{
				{Src: noopSrc, RawOptions: "k1=v2,k3=v4"},
			},
		},

		"A configuration with the plugin on the same line without options should load the plugin correctly.": {
			config: `yeagi { ./_examples/noop/noop.go }`,
			expSrcWithOpts: []corednsyaegi.PluginSrcWithOpts{
				{Src: noopSrc, RawOptions: ""},
			},
		},

		"A correct configuration with multiple plugins and different opts convinations should load the files correctly.": {
			config: `
		yeagi {
			./_examples/noop/noop.go "k1=v2,k3=v4"
			./_examples/noop/noop.go ""
			./_examples/noop/noop.go
			./_examples/noop/noop.go something
		}
				`,
			expSrcWithOpts: []corednsyaegi.PluginSrcWithOpts{
				{Src: noopSrc, RawOptions: "k1=v2,k3=v4"},
				{Src: noopSrc, RawOptions: ""},
				{Src: noopSrc, RawOptions: ""},
				{Src: noopSrc, RawOptions: "something"},
			},
		},

		"A correct configuration with single plugin without options should load the files correctly.": {
			config: `
yeagi { 
	./_examples/noop/noop.go
}
`,
			expSrcWithOpts: []corednsyaegi.PluginSrcWithOpts{
				{Src: noopSrc, RawOptions: ""},
			},
		},

		"A correct configuration with single plugin with options should load the files correctly.": {
			config: `
yeagi { 
	./_examples/noop/noop.go "k1=v2,k3=v4"
}
`,
			expSrcWithOpts: []corednsyaegi.PluginSrcWithOpts{
				{Src: noopSrc, RawOptions: "k1=v2,k3=v4"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			c := caddy.NewTestController("dns", test.config)
			gotSrcWithOpts, err := corednsyaegi.ConfigParse(c)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expSrcWithOpts, gotSrcWithOpts)
			}
		})
	}
}
