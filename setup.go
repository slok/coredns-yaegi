package corednsyaegi

import (
	"fmt"
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("yaegi", setup) }

func setup(c *caddy.Controller) error {
	src, err := ConfigParse(c)
	if err != nil {
		return fmt.Errorf("could setup plugin config: %w", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		coreDNSPlugin, err := NewCoreDNSPlugin(CoreDNSPluginConfig{
			NextPlugin: next,
			PluginsSrc: src,
		})
		if err != nil {
			panic(fmt.Errorf("could not create plugin: %w", err))
		}

		return coreDNSPlugin
	})

	return nil
}

func ConfigParse(c *caddy.Controller) (pluginSrc string, err error) {
	pluginSrcPath := ""
	for c.Next() {
		if !c.NextArg() {
			// If no values then error,
			return "", c.ArgErr()
		}

		pluginSrcPath = c.Val()
	}

	if pluginSrcPath == "" {
		return "", fmt.Errorf("missing plugin file path")
	}

	return pluginSourceCodeLoader(pluginSrcPath)
}

func pluginSourceCodeLoader(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read plugin file: %w", err)
	}

	return string(b), nil
}
