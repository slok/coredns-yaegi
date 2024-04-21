package corednsyaegi

import (
	"context"
	"fmt"

	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"

	"github.com/slok/coredns-yaegi/internal/yaegi"
)

type CoreDNSPlugin struct {
	yaegiPlugin corednsplugin.Handler
}

type CoreDNSPluginConfig struct {
	NextPlugin corednsplugin.Handler
	PluginsSrc string
}

func (c *CoreDNSPluginConfig) defaults() error {
	if c.PluginsSrc == "" {
		return fmt.Errorf("plugin source is required")
	}

	return nil
}

func NewCoreDNSPlugin(config CoreDNSPluginConfig) (*CoreDNSPlugin, error) {
	err := config.defaults()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Load plugin factory.
	pluginfactory, err := yaegi.LoadYaegiPlugin(config.PluginsSrc)
	if err != nil {
		return nil, fmt.Errorf("could not load plugin source: %w", err)
	}

	return &CoreDNSPlugin{
		yaegiPlugin: pluginfactory(config.NextPlugin),
	}, nil
}

func (c CoreDNSPlugin) Ready() bool  { return true }
func (c CoreDNSPlugin) Name() string { return c.yaegiPlugin.Name() }

func (c CoreDNSPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return c.yaegiPlugin.ServeDNS(ctx, w, r)
}
