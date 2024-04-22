package corednsyaegi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("yaegi", setup) }

type PluginSrcWithOpts struct {
	Src        string
	RawOptions string
}

func setup(c *caddy.Controller) error {
	srcs, err := ConfigParse(c)
	if err != nil {
		return fmt.Errorf("could setup plugin config: %w", err)
	}

	// Reverse plugins to execute in declared order.
	slices.Reverse(srcs)

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		for _, src := range srcs {
			next, err = NewCoreDNSPlugin(CoreDNSPluginConfig{
				NextPlugin: next,
				PluginsSrc: src.Src,
				RawOptions: src.RawOptions,
			})
			if err != nil {
				panic(fmt.Errorf("could not create plugin: %w", err))
			}
		}

		return next
	})

	return nil
}

func ConfigParse(c *caddy.Controller) ([]PluginSrcWithOpts, error) {
	pluginSrcPaths := []PluginSrcWithOpts{}

	// Get block.
	c.Next()
	for c.NextBlock() {
		pluginSrc := PluginSrcWithOpts{Src: c.Val()}
		opts := c.RemainingArgs()
		// Sanitize end block (remove end curly brace) if we are on the same line.
		if len(opts) > 0 && opts[len(opts)-1] == "}" {
			opts = opts[:len(opts)-1]
		}

		switch {
		case len(opts) == 1:
			pluginSrc.RawOptions = opts[0]
		case len(opts) > 1:
			return nil, fmt.Errorf("0 or 1 argument raw string options is required, %d received: %#v", len(opts), opts)
		}

		pluginSrcPaths = append(pluginSrcPaths, pluginSrc)
	}

	if len(pluginSrcPaths) == 0 {
		return nil, fmt.Errorf("missing plugin file paths")
	}

	pluginsSrc := []PluginSrcWithOpts{}
	for _, p := range pluginSrcPaths {
		pluginSrc, err := pluginSourceCodeLoader(p.Src)
		if err != nil {
			return nil, fmt.Errorf("could not load plugin src on %q: %w", p, err)
		}
		pluginsSrc = append(pluginsSrc, PluginSrcWithOpts{
			Src:        pluginSrc,
			RawOptions: p.RawOptions,
		})
	}

	return pluginsSrc, nil
}

func pluginSourceCodeLoader(pathOrURL string) (string, error) {
	u, err := url.ParseRequestURI(pathOrURL)
	if err != nil || !strings.HasPrefix(u.Scheme, "http") {
		// Load from local file.
		b, err := os.ReadFile(pathOrURL)
		if err != nil {
			return "", fmt.Errorf("could not read plugin file: %w", err)
		}
		return string(b), nil
	}

	// Download plugin.
	resp, err := http.Get(pathOrURL)
	if err != nil {
		return "", fmt.Errorf("could not download plugin: %w", err)
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
