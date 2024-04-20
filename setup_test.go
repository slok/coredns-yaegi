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

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	return plugin{next: next}
}

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "noop" }
`

func TestConfigParse(t *testing.T) {
	tests := map[string]struct {
		config string
		expErr bool
		expSrc []string
	}{
		"No config should error.": {
			config: ``,
			expErr: true,
		},

		"A correct configuration with multiple plugins should load the files correctly.": {
			config: `
		yeagi { ./_examples/noop/noop.go ./_examples/noop/noop.go ./_examples/noop/noop.go }
		`,
			expSrc: []string{noopSrc, noopSrc, noopSrc},
		},

		"A correct configuration with multiple plugins in different lines should load the files correctly.": {
			config: `
		yeagi {
			./_examples/noop/noop.go
			./_examples/noop/noop.go
		}
		`,
			expSrc: []string{noopSrc, noopSrc},
		},

		"A correct configuration with single plugin should load the files correctly.": {
			config: `
yeagi { ./_examples/noop/noop.go }
`,
			expSrc: []string{noopSrc},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			c := caddy.NewTestController("dns", test.config)
			gotSrc, err := corednsyaegi.ConfigParse(c)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expSrc, gotSrc)
			}
		})
	}
}
