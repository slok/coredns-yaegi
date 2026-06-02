# Changelog

## [Unreleased]

### Added

- Initial Yaegi plugin engine for CoreDNS.
- Allow loading multiple plugins.
- Allow loading plugins from local files.
- Allow loading plugins from HTTP public URL.
- Allow using package `github.com/coredns/coredns/plugin`.
- Allow using package `github.com/coredns/coredns/request`.
- Allow using package `github.com/miekg/dns`.
- Allow using package `github.com/coredns/coredns/plugin/metrics`.
- Allow using package `github.com/prometheus/client_golang/prometheus`.
- Allow using package `github.com/prometheus/client_golang/prometheus/promauto`.
- Allow using package `github.com/coredns/coredns/plugin/pkg/rcode`.
- Allow using string configuration on each plugin load.
- Use CoreDNS v1.11.1.
- Use Yaegi v0.16.1.

### Changed

- Upgrade CoreDNS to v1.14.3.
- Upgrade Go to 1.26.
- Upgrade `github.com/coredns/caddy` to v1.1.4.
- Upgrade `github.com/miekg/dns` to v1.1.72.
- Upgrade `github.com/prometheus/client_golang` to v1.23.2.
- Upgrade `github.com/stretchr/testify` to v1.11.1.
- Regenerate Yaegi custom symbols for CoreDNS 1.14.3 APIs.

[unreleased]: https://github.com/slok/coredns-yaegi/compare/v0.1.0...HEAD
