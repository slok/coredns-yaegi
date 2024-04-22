package noop

import (
	"context"
	"encoding/json"
	"fmt"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type jsonOptions struct {
	Msg string `json:"msg"`
}

type plugin struct {
	next corednsplugin.Handler
	opts jsonOptions
}

func NewPlugin(next corednsplugin.Handler, rawOpts string) corednsplugin.Handler {
	var opts jsonOptions
	if rawOpts != "" {
		err := json.Unmarshal([]byte(rawOpts), &opts)
		if err != nil {
			panic("could not load print plugin JSON options:" + err.Error())
		}
	}

	if opts.Msg == "" {
		opts.Msg = "Hello world!"
	}

	return plugin{
		next: next,
		opts: opts,
	}
}

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	fmt.Printf("Start %q\n", p.opts.Msg)
	defer fmt.Printf("End %q\n", p.opts.Msg)
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "print" }
