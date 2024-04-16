package corednsyaegi_test

import (
	"strings"
	"testing"

	"github.com/coredns/caddy"
	corednsyaegi "github.com/slok/coredns-yaegi"
	"github.com/stretchr/testify/assert"
)

func TestRecordsParse(t *testing.T) {
	tests := map[string]struct {
		config string
		expErr bool
		expSrc string
	}{
		"No config should error.": {
			config: ``,
			expErr: true,
		},

		"A correct configuration should load the file correctly.": {
			config: `
yeagi ./_examples/noop/noop.go
`,
			expSrc: `
package noop

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
	return dns.RcodeSuccess, nil
}

func (p plugin) Name() string { return "noop" }
`,
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
				assert.Equal(strings.TrimSpace(test.expSrc), strings.TrimSpace(gotSrc))
			}
		})
	}
}
