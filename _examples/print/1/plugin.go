package noop

import (
	"context"
	"fmt"

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
	fmt.Println("Start Plugin 1")
	defer fmt.Println("End Plugin 1")
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "print1" }
