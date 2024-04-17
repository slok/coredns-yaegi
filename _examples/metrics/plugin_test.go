package metrics

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	plugintest "github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	corednsyaegi "github.com/slok/coredns-yaegi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	pluginSrc, _ := os.ReadFile("./plugin.go")

	tests := map[string]struct {
		request    func() *dns.Msg
		expMetrics string
		expCode    int
		expErr     bool
	}{
		"Making a requuest should measure as a metric and return the correct rcode.": {
			request: func() *dns.Msg {
				msg := dns.Msg{}
				msg.SetQuestion("slok.dev.", dns.TypeA)
				return &msg
			},
			expCode: dns.RcodeSuccess,
			expMetrics: `
# HELP coredns_yaegi_metrics_dns_detail_requests_total Counts the number of DNS requests with more detail than core metrics.
# TYPE coredns_yaegi_metrics_dns_detail_requests_total counter
coredns_yaegi_metrics_dns_detail_requests_total{server_name="slok.dev.",type="A"} 1
`,
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

				// Get metrics.
				promHandler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
				resp := httptest.NewRecorder()
				promHandler.ServeHTTP(resp, &http.Request{Method: http.MethodGet})

				metrics, err := io.ReadAll(resp.Body)
				require.NoError(err)

				gotMetrics := strings.TrimSpace(string(metrics))
				expectedMetrics := strings.Split(strings.TrimSpace(test.expMetrics), "\n")
				for _, m := range expectedMetrics {
					assert.Contains(gotMetrics, m)
				}
			}
		})
	}
}
