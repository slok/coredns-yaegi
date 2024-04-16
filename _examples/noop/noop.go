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
