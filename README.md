# coredns-yaegi

A [CoreDNS] plugin that allows loading other plugins on the fly using [Yaegi] engine.

## Why do I need this?

So, maybe you are wondering... we already have regular [CoreDNS plugins][coredns-plugins], why should I need this?

Opposite to the regular CoreDNS plugins, using this plugin, you don't need to recompile coreDNS and configure it to add them, [Yaegi] will be able to load them on the fly and execute as a restricted sandboxed Go runtime, however, these plugins are very limited compared to the regular plugins, Depending on the use case one or the other would be the way to go.

### Pros

- Easy to load.
- Fast to develop and setup.
- Easy to maintain.
- Portable.

### Cons

- Limited to the usage standard library and some CoreDNS packages.
- No external libraries allowed.
- Need to be on a single file.
- A bit less performant.
- Limited by [Yaegi] features and bugs

So, checking the pros and cons, you may get the idea of these, when you need a simple plugin like blocking/allow DNS, rewriting/mutating the DNS request... these plugins are very easy to set and useful. On the opposite, if you want to connect with external services like a Redis or Kubernetes, you should stick to regular plugins.

## Features and restrictions

- Made in Go.
- Heavily dependant on [Yaegi].
- Almost similar plugin syntax to the regular coreDNS plugins.
- No external dependencies allowed (except some CoreDNS and helper packages).
- Ability to measure using CoreDNS plugins Prometheus.
- Plugin must be on a single file.
- Configurable from a local file or to download from an HTTP endpoint (e.g public repo or public/private gist).

## Use cases

- DNS allow/block logic.
- Complex logic based DNS rewriting.
- Audit.

## Writing plugins

To write a plugin you only need to implement a method: `func NewPlugin(next corednsplugin.Handler) corednsplugin.Handler`.

This method should return the plugin itself. Example of a NOOP plugin:

```go
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
    return corednsplugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p plugin) Name() string { return "noop" }
```

### Allowed packages?

Besides the ability to use all the Go  standard library, you can use these external packages:

- [github.com/coredns/coredns/plugin](https://pkg.go.dev/github.com/coredns/coredns/plugin)
- [github.com/coredns/coredns/request](https://pkg.go.dev/github.com/coredns/coredns/request)
- [github.com/miekg/dns](https://pkg.go.dev/github.com/miekg/dns)
- [github.com/coredns/coredns/plugin/metrics](https://pkg.go.dev/github.com/coredns/coredns/plugin/metrics)
- [github.com/prometheus/client_golang/prometheus](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
- [github.com/prometheus/client_golang/prometheus/promauto](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promauto)

Do you miss any package? Depending on the general usefulness and safety we may add it so is available.

## Configuration

The configuration is very simple, a yaegi plugin can load a single plugin per block (if you need more, you can make your plugin have a chaining logic), and it can be loaded from a local file or a public HTTP URL. Lets see some examples:

Load from a local file:

```
. {
    yaegi /tmp/my-plugin.go
    forward . 1.1.1.1
}
```

Load from a URL:

```
. {
    yaegi https://raw.githubusercontent.com/slok/coredns-yaegi/main/_examples/allowlist/plugin.go
    forward . 1.1.1.1
}
```

Multiple blocks and different plugins:

```
slok.dev {
    log
    yaegi /plugins/slok.go
    forward . 1.1.1.1
}

google.com {
    log
    errors
    forward . 1.1.1.1
}

twitter.com {
    log
    yaegi ./twitter.go
    forward . 8.8.8.8
}

. {
    forward . 8.8.8.8
    log
    yaegi https://example.com/generic.go
    errors
    cache
}
```

## Ready to use image

You can build you own coreDNS image and set this plugin using [CoreDNS docs](https://coredns.io/manual/explugins/), however we provide a CoreDNS image that has the yaegi plugin ready to be used: 

```bash
docker pull ghcr.io/slok/coredns-yaegi
```

## Example

Let's load our [example allowlist](_examples/allowlist) plugin directly from an URL file.

```bash
$ cat ./coredns.config
. {
    yaegi https://raw.githubusercontent.com/slok/coredns-yaegi/main/_examples/allowlist/plugin.go
    forward . 1.1.1.1
}

$ docker run --rm -it --network host \
    -v ${PWD}/coredns.config:/tmp/coredns.config \
    ghcr.io/slok/coredns-yaegi:2620714c574e80582ff0072086acbc6c53072d08 \
    --conf /tmp/coredns.config  -dns.port=1053

$ dig @localhost -p 1053 a github.com +short
$ dig @localhost -p 1053 a google.com +short
$ dig @localhost -p 1053 a xabi.dev +short
172.67.193.60
104.21.20.171
```

[CoreDNS]: https://coredns.io/
[yaegi]: https://github.com/traefik/yaegi
[coredns-plugins]: https://coredns.io/explugins
