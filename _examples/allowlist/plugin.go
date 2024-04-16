package allowlist

import (
	"context"
	"fmt"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Add here the domains that you want to be allowed.
var allowlist = map[string]struct{}{
	"slok.dev.":          {},
	"xabi.dev.":          {},
	"cloudwarlocks.com.": {},
}

type plugin struct {
	next corednsplugin.Handler
}

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	return plugin{next: next}
}

func (p plugin) Name() string { return "allowlist" }

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	_, ok := allowlist[state.Name()]
	if !ok {
		resp := new(dns.Msg)
		resp.SetRcode(r, dns.RcodeNameError)
		err := w.WriteMsg(resp)
		if err != nil {
			fmt.Println(err)
			return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
		}

		return dns.RcodeNameError, nil
	}

	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}
