package metrics

import (
	"context"
	"os"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type plugin struct {
	next       corednsplugin.Handler
	reqCounter *prometheus.CounterVec
}

func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler {
	subsystem := os.Getenv("COREDNS_YAEGI_PLUGIN_METRICS_PREFIX")
	if subsystem == "" {
		subsystem = "yaegi_metrics"
	}

	return plugin{
		next: next,
		reqCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: corednsplugin.Namespace,
			Subsystem: subsystem,
			Name:      "dns_detail_requests_total",
			Help:      "Counts the number of DNS requests with more detail than core metrics.",
		}, []string{"server_name", "type"}),
	}
}

func (p plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	p.reqCounter.WithLabelValues(state.Name(), state.Type()).Inc()
	return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "metrics" }
